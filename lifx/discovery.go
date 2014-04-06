package lifx

import "net"
import "log"
import "time"
import "bytes"

func DiscoverGateways(findAll bool) ([]*Gateway, error) {
	addr, _ := net.ResolveUDPAddr("udp4", ":56700")
	conn, _ := net.ListenUDP("udp4", addr)
	broadcast, _ := net.ResolveUDPAddr("udp4", "255.255.255.255:56700")
	buf := make([]byte, 2048)
	outbuf := new(bytes.Buffer)
	gateways := make(map[string]*Gateway)

	message, _ := CreateMessage(MSG_GET_PAN_GATEWAY)
	EncodeMessage(message, outbuf)

	conn.SetDeadline(time.Now().Add(10 * time.Second))
	_, err := conn.WriteToUDP(outbuf.Bytes(), broadcast)

	if err != nil {
		return nil, err
	}

	for {
		_, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Print("ERROR: Error reading message: ", err)
			break
		}

		message, err := DecodeMessage(bytes.NewBuffer(buf))
		if err != nil {
			log.Print("ERROR: Error decoding message: ", err)
			continue
		}

		if message.Header.PacketType == MSG_STATE_PAN_GATEWAY {
			payload := message.Payload.(*StatePanGateway)
			address := remote.String()
			gateway := NewGateway(address, message.Header.Site, payload.Service)
			gateways[gateway.Id] = gateway

			if !findAll {
				break
			}
		}
	}

	var gatewayList []*Gateway
	for k := range gateways {
		gatewayList = append(gatewayList, gateways[k])
	}

	return gatewayList, err
}
