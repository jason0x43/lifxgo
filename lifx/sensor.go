package lifx

const (
	MSG_GET_AMBIENT_LIGHT uint16    = 401
	MSG_STATE_AMBIENT_LIGHT uint16  = 402
	MSG_GET_DIMMER_VOLTAGE uint16   = 403
	MSG_STATE_DIMMER_VOLTAGE uint16 = 404
)

type GetAmbientLight struct {
}

type StateAmbientLight struct {
	Lux float32
}

type GetDimmerVoltage struct {
}

type StateDimmerVoltage struct {
	Voltage uint32
}
