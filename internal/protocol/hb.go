package protocol

// Protocol Format
// +------+------------+--------------------+
// | Type | Heartbeat ID | Timestamp (Unix)  |
// +------+------------+--------------------+
// |  1   |   4 bytes   |      8 bytes       |
// +------+------------+--------------------+

const (
	TypeHeartbeatRequest  = 1
	TypeHeartbeatResponse = 2
)
