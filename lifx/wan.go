package lifx

const (
	MSG_CONNECT_PLAIN uint16 = 201
	MSG_CONNECT_KEY   uint16 = 202
	MSG_STATE_CONNECT uint16 = 203
	MSG_SUB           uint16 = 204
	MSG_UNSUB         uint16 = 205
	MSG_STATE_SUB     uint16 = 206
)

type ConnectPlain struct {
	User [32]byte
	Pass [32]byte
}

type ConnectKey struct {
	AuthKey [32]byte
}

type StateConnect struct {
	AuthKey [32]byte
}

type Sub struct {
	Target [8]byte
	Site   [6]byte
	Device bool // 0 - Targets a device. 1 - Targets a tag
}

type Unsub struct {
	Target [8]byte
	Site   [6]byte
	Device bool // 0 - Targets a device. 1 - Targets a tag
}

type StateSub struct {
	Target [8]byte
	Site   [6]byte
	Device bool // 0 - Targets a device. 1 - Targets a tag
}
