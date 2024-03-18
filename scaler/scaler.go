package scaler

import "auto_scaling/storage"

type Scaler interface {
	ApplyCall(*storage.Call) error
}