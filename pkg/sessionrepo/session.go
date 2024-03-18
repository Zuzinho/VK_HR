package sessionrepo

import (
	"time"
)

type Session struct {
	Sub string
	Exp time.Time
}

type SessionPacker interface {
	Pack(sub string) (*string, error)
	Unpack(inToken string) (*Session, error)
}
