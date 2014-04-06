package lifx

const (
	MSG_SET_SITE uint16               = 1
	MSG_GET_PAN_GATEWAY uint16        = 2
	MSG_STATE_PAN_GATEWAY uint16      = 3
	MSG_GET_TIME uint16               = 4
	MSG_SET_TIME uint16               = 5
	MSG_STATE_TIME uint16             = 6
	MSG_GET_RESET_SWITCH uint16       = 7
	MSG_STATE_RESET_SWITCH uint16     = 8
	MSG_GET_MESH_INFO uint16          = 12
	MSG_STATE_MESH_INFO uint16        = 13
	MSG_GET_MESH_FIRMWARE uint16      = 14
	MSG_STATE_MESH_FIRMWARE uint16    = 15
	MSG_GET_WIFI_INFO uint16          = 16
	MSG_STATE_WIFI_INFO uint16        = 17
	MSG_GET_WIFI_FIRMWARE uint16      = 18
	MSG_STATE_WIFI_FIRMWARE uint16    = 19
	MSG_GET_POWER uint16              = 20
	MSG_SET_POWER uint16              = 21
	MSG_STATE_POWER uint16            = 22
	MSG_GET_LABEL uint16              = 23
	MSG_SET_LABEL uint16              = 24
	MSG_STATE_LABEL uint16            = 25
	MSG_GET_TAGS uint16               = 26
	MSG_SET_TAGS uint16               = 27
	MSG_STATE_TAGS uint16             = 28
	MSG_GET_TAG_LABELS uint16         = 29
	MSG_SET_TAG_LABELS uint16         = 30
	MSG_STATE_TAG_LABELS uint16       = 31
	MSG_GET_VERSION uint16            = 32
	MSG_STATE_VERSION uint16          = 33
	MSG_GET_INFO uint16               = 34
	MSG_STATE_INFO uint16             = 35
	MSG_GET_MCU_RAIL_VOLTAGE uint16   = 36
	MSG_STATE_MCU_RAIL_VOLTAGE uint16 = 37
	MSG_REBOOT uint16                 = 38
)

type SetSite struct {
	Site [6]byte
}

type GetPanGateway struct {
}

type StatePanGateway struct {
	Service byte
	Port uint32
}

type GetTime struct {
}

type SetTime struct {
	Time uint64 // nanoseconds since epoch
}

type StateTime struct {
	Time uint64 // nanoseconds since epoch
}

type GetResetSwitch struct {
}

type StateResetSwitch struct {
	Position byte
}

type GetMeshInfo struct {
}

type StateMeshInfo struct {
	Signal float32 // milliwatts
	Tx uint32 // bytes
	Rx uint32 // bytes
	McuTemperature uint16 // Deci-celsius. 25.45 celsius is 2545
}

type GetMeshFirmware struct {
}

type StateMeshFirmware struct {
	Build uint64 // Firmware build nanoseconds since epoch
	Install uint64 // Firmware install nanoseconds since epoch
	Version uint32 // Firmware human readable version
}

type GetWifiInfo struct {
}

type StateWifiInfo struct {
	Signal float32 // Milliwatts
	Tx uint32 // Bytes
	Rx uint32 // Bytes
	McuTemperature int16 // Deci-celsius. 25.45 celsius is 2545
}

type GetWifiFirmware struct {
}

type StateWifiFirmware struct {
	Build uint64 // Firmware build nanoseconds since epoch
	Install uint64 // Firmware install nanoseconds since epoch
	Version uint32 // Firmware human readable version
}

type GetPower struct {
}

type SetPower struct {
	Level uint16 // 0 Standby. > 0 On
}

type StatePower struct {
	Level uint16 // 0 Standby. > 0 On
}

type GetLabel struct {
}

type SetLabel struct {
	Label [32]byte
}

type StateLabel struct {
	Label [32]byte
}

type GetTags struct {
}

type SetTags struct {
	Tags uint64
}

type StateTags struct {
	Tags uint64
}

type GetTagLabels struct {
	Tags uint64
}

type SetTagLabels struct {
	Tags uint64
	Label [32]byte
}

type StateTagLabels struct {
	Tags uint64
	Label [32]byte
}

type GetVersion struct {
}

type StateVersion struct {
	Vendor uint32
	Product uint32
	Version uint32
}

type GetInfo struct {
}

type StateInfo struct {
	Time uint64 // Nanoseconds since epoch
	Uptime uint64 // Nanoseconds since boot
	Downtype uint64 // Nanoseconds off last power cycle
}

type GetMcuRailVoltage struct {
}

type StateMcuRailVoltage struct {
	Voltage uint32
}

type Reboot struct {
}
