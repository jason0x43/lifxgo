package lifx

import "net"
import "log"
import "time"
import "bytes"

var zeroTime time.Time

func DiscoverGateways(findAll bool) ([]*Gateway, error) {
	addr, _ := net.ResolveUDPAddr("udp4", ":56700")
	conn, _ := net.ListenUDP("udp4", addr)
	broadcast, _ := net.ResolveUDPAddr("udp4", "255.255.255.255:56700")
	buf := make([]byte, 2048)
	outbuf := new(bytes.Buffer)
	gateways := make(map[string]*Gateway)

	log.Print("Creating GetPanGateway message");
	message, _ := NewMessage(MSG_GET_PAN_GATEWAY)
	EncodeMessage(message, outbuf)

	log.Print("Sending discovery message");
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	_, err := conn.WriteToUDP(outbuf.Bytes(), broadcast)

	if err != nil {
		return nil, err
	}

	for {
		_, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			if !err.(net.Error).Timeout() {
				log.Print("ERROR: Error reading message: ", err)
			}
			break
		}
		log.Print("Received data");

		message, err := DecodeMessage(bytes.NewBuffer(buf))
		if err != nil {
			log.Print("ERROR: Error decoding message: ", err)
			continue
		}
		log.Println("Decoded message", message);

		if message.PacketType == MSG_STATE_PAN_GATEWAY {
			payload := message.Payload.(*StatePanGateway)
			address := remote.String()
			gateway := NewGateway(address, message.Site, payload.Service)
			gateways[gateway.Id] = gateway

			if !findAll {
				break
			}
		}
	}

	conn.SetDeadline(zeroTime)

	var gatewayList []*Gateway
	for k := range gateways {
		gatewayList = append(gatewayList, gateways[k])
	}

	return gatewayList, err
}
