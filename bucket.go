package bucket

import "errors"

var ErrFull = errors.New("Bucket is full")

type Bucket interface {
	C() chan interface{}
	Put(interface{}) error
}
