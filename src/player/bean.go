package player

import "time"

type WordPlayInfo struct {
	IdleTime   time.Duration
	EngMp3Path string
	ChMp3Path  string
	Word       string
}
