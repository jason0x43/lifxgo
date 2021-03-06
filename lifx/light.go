package lifx

import "fmt"
import "log"
import "time"
import "bytes"
import "encoding/hex"

type Light struct {
	Id      string
	Site    [6]byte
	Color   hsbk
	Label   string
	seen    time.Time
	gateway *Gateway
}

func (light *Light) String() string {
	return fmt.Sprintf("{Light %s '%s'}", light.Id, light.Label)
}

func NewLight(gateway *Gateway, message *Message) *Light {
	light := Light{gateway: gateway}
	light.Id = hex.EncodeToString(message.BulbAddress[:])
	light.HandleMessage(message)

	return &light
}

func (light *Light) NewMessage(messageType uint16) (*Message, error) {
	message, err := light.gateway.NewMessage(MSG_SET_POWER)
	if err != nil {
		return nil, err
	}
	message.BulbAddress = light.Site
	return message, nil
}

func (light *Light) SetPower(percent float32) error {
	level := uint16(percent * 0xFFFF)
	log.Println("Setting light level to", level);

	message, err := light.NewMessage(MSG_SET_POWER)
	if err != nil {
		return err
	}

	message.Payload.(*SetPower).Level = level
	err = light.gateway.send(message)
	if err != nil {
		return err
	}

	_, err = light.gateway.transport.Read(light.gateway.readBuf)
	if err != nil {
		return err
	}

	message, err = DecodeMessage(bytes.NewBuffer(light.gateway.readBuf))
	log.Printf("Got message: %#v", message.Payload)

	return err
}

func (light *Light) HandleMessage(message *Message) {
	light.seen = time.Now()

	switch message.Payload.(type) {
	case *StateLight:
		payload := message.Payload.(*StateLight)
		light.Color = payload.Color
		light.Label = string(payload.Label[:])
	}
}

const (
	WAVEFORM_SAW       uint8 = 0
	WAVEFORM_SINE      uint8 = 1
	WAVEFORM_HALF_SINE uint8 = 2
	WAVEFORM_TRIANGLE  uint8 = 3
	WAVEFORM_PULSE     uint8 = 4
)

const (
	MSG_LIGHT_GET          uint16 = 101
	MSG_LIGHT_SET          uint16 = 102
	MSG_SET_WAVEFORM       uint16 = 103
	MSG_SET_DIM_ABSOLUTE   uint16 = 104
	MSG_SET_DIM_RELATIVE   uint16 = 105
	MSG_SET_RGBW           uint16 = 106
	MSG_STATE_LIGHT        uint16 = 107
	MSG_GET_RAIL_VOLTAGE   uint16 = 108
	MSG_STATE_RAIL_VOLTAGE uint16 = 109
	MSG_GET_TEMPERATURE    uint16 = 110
	MSG_STATE_TEMPERATURE  uint16 = 111
)

type hsbk struct {
	Hue        uint16 // 0..65_535 scaled to 0° - 360°
	Saturation uint16 // 0..65_535 scaled to 0% - 100%
	Brightness uint16 // 0..65_535 scaled to 0% - 100%
	Kelvin     uint16 // Explicit 2_400..10_000
}

type LightGet struct {
}

type LightSet struct {
	Stream   uint8 // 0 is no stream
	Color    hsbk
	Duration uint32 // Milliseconds
}

type SetWaveform struct {
	Stream    uint8 // 0 is no stream
	Transient bool
	Color     hsbk
	Period    uint32 // Milliseconds per cycle
	Cycles    float32
	DutyCycle int16
	Waveform  uint8
}

type SetDimAbsolute struct {
	Brightness int16  // 0 is no change
	Duration   uint32 // Milliseconds
}

type SetDimRelative struct {
	Brightness int32  // 0 is no change
	Duration   uint32 // Milliseconds
}

type rgbw struct {
	Red   uint16
	Green uint16
	Blue  uint16
	White uint16
}

type SetRgbw struct {
	Color rgbw
}

type StateLight struct {
	Color hsbk
	Dim   int16
	Power uint16
	Label [32]byte
	Tags  uint64
}
func (sl *StateLight) String() string {
	return fmt.Sprintf("{%T Label='%s', Power=%d}", sl, sl.Label[:], sl.Power)
}

type GetRailVoltage struct {
}

type StateRailVoltage struct {
	Voltage uint32
}

type GetTemperature struct {
}

type StateTemperature struct {
	Temperature int16 // Deci-celsius. 25.45 celsius is 2545
}
