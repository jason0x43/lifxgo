package lifx

import "net"
import "log"
import "time"
import "bytes"
import "encoding/hex"

const (
	TCP_MAX_ATTEMPTS int = 3
	SERVICE_UDP byte = 1
	SERVICE_TCP byte = 2
)

type Gateway struct {
	Site [6]byte
	Id string
	Lights map[string]*Light
	Protocol byte
	Address string
	transport net.Conn
	readBuf []byte
	writeBuf *bytes.Buffer
}

func NewGateway(address string, site [6]byte, protocol byte) *Gateway {
	id := hex.EncodeToString(site[:])
	return &Gateway{
		Site: site,
		Id: id,
		Protocol: protocol,
		Lights: make(map[string]*Light, 0),
		Address: address,
		readBuf: make([]byte, 2048),
		writeBuf: new(bytes.Buffer),
	}
}

func (gateway *Gateway) dial(netw string) error {
	transport, err := net.DialTimeout(netw, gateway.Address, 10 * time.Second)
	if err == nil {
		gateway.transport = transport
	}
	return err
}

func (gateway *Gateway) Dial() error {
	var err error

	if gateway.Protocol == SERVICE_TCP {
		log.Print("Establishing TCP connection...")
		err = gateway.dial("tcp4")
	}

	if gateway.transport == nil {
		gateway.Protocol = SERVICE_UDP
		log.Print("Establishing UDP connection...")
		err = gateway.dial("udp4")
	}

	return err
}

func (gateway *Gateway) send(message *Message) error {
	EncodeMessage(message, gateway.writeBuf)
	_, err := gateway.transport.Write(gateway.writeBuf.Bytes())
	return err
}

func (gateway *Gateway) createMessage(messageType uint16) (*Message, error) {
	message, err := CreateMessage(messageType)
	message.Header.Site = gateway.Site
	return message, err
}

func (gateway *Gateway) RefreshLights() ([]*Light, error) {
	msg, _ := gateway.createMessage(MSG_LIGHT_GET)
	err := gateway.send(msg)
	if err != nil {
		log.Print("ERROR: Error sending state message: ", err)
		return nil, err
	}

	gateway.transport.SetReadDeadline(time.Now().Add(2 * time.Second))
	lights := make(map[string]*Light)

	for {
		_, err := gateway.transport.Read(gateway.readBuf)
		if err != nil {
			if !err.(net.Error).Timeout() {
				log.Print("ERROR: Error reading: ", err)
			}
			break
		}

		message, err := DecodeMessage(bytes.NewBuffer(gateway.readBuf))
		if err != nil {
			log.Print("ERROR: Error decoding message: ", err)
			continue
		}

		if message.Header.PacketType == MSG_STATE_LIGHT {
			light := NewLight(gateway, message)
			lights[light.Id] = light
		}
	}

	var zeroTime time.Time
	gateway.transport.SetReadDeadline(zeroTime)
	gateway.Lights = lights

	var lightList []*Light
	for k := range lights {
		lightList = append(lightList, lights[k])
	}
	return lightList, nil
}
