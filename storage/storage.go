package storage

import (
	"auto_scaling/lib/e"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"time"
)

var ErrEmpty = errors.New("no calls exist")
var ErrNoDir = errors.New("dir does not exist")
var ErrOutOfLimit = errors.New("amount is out of set limit")
var ErrUnknownType = errors.New("unknown type of call")

type Storage interface {
	Save(*Call) error
	PickLastCalls(string) ([]*Call, error)
}

type Call struct {
	Type      int
	Amount    int
	UserName  string
	CreatedAt time.Time
}

func (p *Call) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.CreatedAt.String()); err != nil {
		return "", e.WrapErr("can't create hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.WrapErr("can't create hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
