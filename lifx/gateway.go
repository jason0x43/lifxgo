package lifx

import "net"
import "log"
import "time"
import "bytes"
import "encoding/hex"

const (
	TCP_MAX_ATTEMPTS int  = 3
	SERVICE_UDP      byte = 1
	SERVICE_TCP      byte = 2
)

type Gateway struct {
	Site      [6]byte
	Id        string
	Lights    map[string]*Light
	Protocol  byte
	Address   string
	transport net.Conn
	readBuf   []byte
	writeBuf  *bytes.Buffer
}

func NewGateway(address string, site [6]byte, protocol byte) *Gateway {
	id := hex.EncodeToString(site[:])
	return &Gateway{
		Site:     site,
		Id:       id,
		Protocol: protocol,
		Lights:   make(map[string]*Light, 0),
		Address:  address,
		readBuf:  make([]byte, 2048),
		writeBuf: new(bytes.Buffer),
	}
}

func (gateway *Gateway) Dial() error {
	var err error
	var transport net.Conn

	if gateway.Protocol == SERVICE_TCP {
		log.Print("Establishing TCP connection...")
		transport, err = net.DialTimeout("tcp4", gateway.Address, 10*time.Second)
	}

	if transport == nil {
		gateway.Protocol = SERVICE_UDP
		log.Print("Establishing UDP connection...")
		transport, err = net.DialTimeout("udp4", gateway.Address, 10*time.Second)
	}

	if err == nil {
		gateway.transport = transport
	}
	return err
}

func (gateway *Gateway) send(message *Message) error {
	EncodeMessage(message, gateway.writeBuf)
	_, err := gateway.transport.Write(gateway.writeBuf.Bytes())
	return err
}

func (gateway *Gateway) NewMessage(messageType uint16) (*Message, error) {
	message, err := NewMessage(messageType)
	message.Site = gateway.Site
	return message, err
}

func (gateway *Gateway) RefreshLights() ([]*Light, error) {
	msg, _ := gateway.NewMessage(MSG_LIGHT_GET)
	err := gateway.send(msg)
	if err != nil {
		log.Print("ERROR: Error sending state message: ", err)
		return nil, err
	}

	gateway.transport.SetReadDeadline(time.Now().Add(3 * time.Second))
	lights := make(map[string]*Light)

	for {
		_, err := gateway.transport.Read(gateway.readBuf)
		if err != nil {
			if !err.(net.Error).Timeout() {
				log.Print("ERROR: Error reading: ", err)
			}
			break
		}
		log.Print("Received data");

		message, err := DecodeMessage(bytes.NewBuffer(gateway.readBuf))
		if err != nil {
			log.Print("ERROR: Error decoding message: ", err)
			continue
		}
		log.Println("Decoded message", message);

		if message.PacketType == MSG_STATE_LIGHT {
			light := NewLight(gateway, message)
			lights[light.Id] = light
		}
	}

	gateway.Lights = lights

	var lightList []*Light
	for k := range lights {
		lightList = append(lightList, lights[k])
	}
	return lightList, nil
}
