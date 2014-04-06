package lifx

import "fmt"
import "bytes"
import "errors"
import "encoding/binary"

type header struct {
	Size uint16
	Protocol uint16
	Reserved1 [4]byte
	BulbAddress [6]byte
	Reserved2 [2]byte
	Site [6]byte
	Reserved3 [2]byte
	Timestamp uint64
	PacketType uint16
	Reserved4 [2]byte
}

func (item *header) init() *header {
	item.Protocol = 13312
	item.Size = 36
	return item
}

type payload interface {}

type Message struct {
	Header *header
	Payload payload
}

func (message *Message) String() string {
	return fmt.Sprintf("<<%T:%+v %T:%+v>>", message.Header, message.Header, message.Payload, message.Payload)
}

func DecodeMessage(buf *bytes.Buffer) (*Message, error) {
	header := new(header).init()
	binary.Read(buf, binary.LittleEndian, header)
	payload, err := decodePayload(header.PacketType, buf)
	if err != nil {
		return nil, err
	}
	header.Size += (uint16)(binary.Size(payload))
	return &Message{header, payload}, nil
}

func EncodeMessage(message *Message, buf *bytes.Buffer) error {
	err := binary.Write(buf, binary.LittleEndian, message.Header)
	if err != nil {
		return err
	}
	err = binary.Write(buf, binary.LittleEndian, message.Payload)
	return err
}

func CreateMessage(typeId uint16) (*Message, error) {
	header := new(header).init()
	header.PacketType = typeId
	payload, err := selectPayload(typeId)
	if err != nil {
		return nil, err
	}
	header.Size += (uint16)(binary.Size(payload))
	return &Message{header, payload}, nil
}

func selectPayload(typeId uint16) (payload, error) {
	switch (typeId) {
		// device
		case MSG_SET_SITE:               return new(SetSite), nil
		case MSG_GET_PAN_GATEWAY:        return new(GetPanGateway), nil
		case MSG_STATE_PAN_GATEWAY:      return new(StatePanGateway), nil
		case MSG_GET_TIME:               return new(GetTime), nil
		case MSG_SET_TIME:               return new(SetTime), nil
		case MSG_STATE_TIME:             return new(StateTime), nil
		case MSG_GET_RESET_SWITCH:       return new(GetResetSwitch), nil
		case MSG_STATE_RESET_SWITCH:     return new(StateResetSwitch), nil
		case MSG_GET_MESH_INFO:          return new(GetMeshInfo), nil
		case MSG_STATE_MESH_INFO:        return new(StateMeshInfo), nil
		case MSG_GET_MESH_FIRMWARE:      return new(GetMeshFirmware), nil
		case MSG_STATE_MESH_FIRMWARE:    return new(StateMeshFirmware), nil
		case MSG_GET_WIFI_INFO:          return new(GetWifiInfo), nil
		case MSG_STATE_WIFI_INFO:        return new(StateWifiInfo), nil
		case MSG_GET_WIFI_FIRMWARE:      return new(GetWifiFirmware), nil
		case MSG_STATE_WIFI_FIRMWARE:    return new(StateWifiFirmware), nil
		case MSG_GET_POWER:              return new(GetPower), nil
		case MSG_SET_POWER:              return new(SetPower), nil
		case MSG_STATE_POWER:            return new(StatePower), nil
		case MSG_GET_LABEL:              return new(GetLabel), nil
		case MSG_SET_LABEL:              return new(SetLabel), nil
		case MSG_STATE_LABEL:            return new(StateLabel), nil
		case MSG_GET_TAGS:               return new(GetTags), nil
		case MSG_SET_TAGS:               return new(SetTags), nil
		case MSG_STATE_TAGS:             return new(StateTags), nil
		case MSG_GET_TAG_LABELS:         return new(GetTagLabels), nil
		case MSG_SET_TAG_LABELS:         return new(SetTagLabels), nil
		case MSG_STATE_TAG_LABELS:       return new(StateTagLabels), nil
		case MSG_GET_VERSION:            return new(GetVersion), nil
		case MSG_STATE_VERSION:          return new(StateVersion), nil
		case MSG_GET_INFO:               return new(GetInfo), nil
		case MSG_STATE_INFO:             return new(StateInfo), nil
		case MSG_GET_MCU_RAIL_VOLTAGE:   return new(GetMcuRailVoltage), nil
		case MSG_STATE_MCU_RAIL_VOLTAGE: return new(StateMcuRailVoltage), nil
		case MSG_REBOOT:                 return new(Reboot), nil
		// light, nil
		case MSG_LIGHT_GET:              return new(LightGet), nil
		case MSG_LIGHT_SET:              return new(LightSet), nil
		case MSG_SET_WAVEFORM:           return new(SetWaveform), nil
		case MSG_SET_DIM_ABSOLUTE:       return new(SetDimAbsolute), nil
		case MSG_SET_DIM_RELATIVE:       return new(SetDimRelative), nil
		case MSG_SET_RGBW:               return new(SetRgbw), nil
		case MSG_STATE_LIGHT:            return new(StateLight), nil
		case MSG_GET_RAIL_VOLTAGE:       return new(GetRailVoltage), nil
		case MSG_STATE_RAIL_VOLTAGE:     return new(StateRailVoltage), nil
		case MSG_GET_TEMPERATURE:        return new(GetTemperature), nil
		case MSG_STATE_TEMPERATURE:      return new(StateTemperature), nil
		// wan, nil
		case MSG_CONNECT_PLAIN:          return new(ConnectPlain), nil
		case MSG_CONNECT_KEY:            return new(ConnectKey), nil
		case MSG_STATE_CONNECT:          return new(StateConnect), nil
		case MSG_SUB:                    return new(Sub), nil
		case MSG_UNSUB:                  return new(Unsub), nil
		case MSG_STATE_SUB:              return new(StateSub), nil
		// wifi, nil
		case MSG_WIFI_GET:               return new(WifiGet), nil
		case MSG_WIFI_SET:               return new(WifiSet), nil
		case MSG_STATE_WIFI:             return new(StateWifi), nil
		case MSG_GET_ACCESS_POINT:       return new(GetAccessPoint), nil
		case MSG_SET_ACCESS_POINT:       return new(SetAccessPoint), nil
		case MSG_STATE_ACCESS_POINT:     return new(StateAccessPoint), nil
		// sensor, nil
		case MSG_GET_AMBIENT_LIGHT:      return new(GetAmbientLight), nil
		case MSG_STATE_AMBIENT_LIGHT:    return new(StateAmbientLight), nil
		case MSG_GET_DIMMER_VOLTAGE:     return new(GetDimmerVoltage), nil
		case MSG_STATE_DIMMER_VOLTAGE:   return new(StateDimmerVoltage), nil
	}
	return nil, errors.New("Invalid type")
}

func decodePayload(typeId uint16, buf *bytes.Buffer) (payload, error) {
	payload, err := selectPayload(typeId)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, payload)
	return payload, err
}
