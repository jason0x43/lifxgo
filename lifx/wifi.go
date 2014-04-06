package lifx

const (
	MSG_WIFI_GET uint16           = 301
	MSG_WIFI_SET uint16           = 302
	MSG_STATE_WIFI uint16         = 303
	MSG_GET_ACCESS_POINT uint16   = 304
	MSG_SET_ACCESS_POINT uint16   = 305
	MSG_STATE_ACCESS_POINT uint16 = 306
)

type WifiGet struct {
	Iface uint8
}

type WifiSet struct {
	Iface uint8
	Active bool
}

type StateWifi struct {
	Iface uint8
	Status uint8
	Ipv4 uint32
	Ipv6 [16]byte
}

type GetAccessPoint struct {
}

type SetAccessPoint struct {
	Iface uint8
	Ssid [32]byte
	Pass [64]byte
	Security uint8
}

type StateAccessPoint struct {
	Iface uint8
	Ssid [32]byte
	Security uint8
	Strength int16
	Channel uint16
}
