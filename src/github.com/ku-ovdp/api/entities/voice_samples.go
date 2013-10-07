package entities

import (
	"time"
)

type VoiceSample struct {
	Id        int
	SessionId int
	Created   time.Time
	Length    time.Duration
	MimeType  string
	Size      int
	Bitrate   int
	AudioURL  string `json:"-"`
}
