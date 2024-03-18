package yandex

import "auto_scaling/storage"

type Scaler struct {
	
}

func New() *Scaler {return &Scaler{}}

func (s *Scaler) ApplyCall(call *storage.Call) error {
	return nil
}
