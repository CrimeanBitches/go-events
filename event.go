package protocol

import (
	"errors"
	"reflect"
)

// EventArgs is todo
type EventArgs interface{}

// Event is todo
type Event struct {
	listeners []interface{}
	chanType  reflect.Type
	argsType  reflect.Type
}

// NewEvent does todo
func NewEvent(t interface{}) *Event {
	return &Event{
		listeners: make([]interface{}, 0),
		chanType:  reflect.PtrTo(reflect.ChanOf(reflect.BothDir, reflect.TypeOf(t))),
		argsType:  reflect.TypeOf(t),
	}
}

// NewEventReceiver does todo
func NewEventReceiver(t interface{}) interface{} {
	return reflect.MakeChan(reflect.ChanOf(reflect.BothDir, reflect.TypeOf(t)), 1).Interface()
}

// NewEventReceiverSize does todo
func NewEventReceiverSize(t interface{}, buffer int) interface{} {
	return reflect.MakeChan(reflect.ChanOf(reflect.BothDir, reflect.TypeOf(t)), buffer).Interface()
}

// Add does todo
func (e *Event) Add(c interface{}) (err error) {
	if !reflect.TypeOf(c).AssignableTo(e.chanType) {
		err = errors.New("invalid chan type")
		return
	}
	e.listeners = append(e.listeners, c)
	return
}

// Remove does todo
func (e *Event) Remove(c interface{}) (ok bool, err error) {
	if !reflect.TypeOf(c).AssignableTo(e.chanType) {
		err = errors.New("invalid chan type")
		return
	}
	for i, ec := range e.listeners {
		if ec == c {
			e.listeners = append(e.listeners[:i], e.listeners[i+1:]...)
			ok = true
			return
		}
	}
	return
}

// Call does todo
func (e *Event) Call(v interface{}) (err error) {
	if !reflect.TypeOf(v).AssignableTo(e.argsType) {
		err = errors.New("invalid argument type")
		return
	}
	for _, c := range e.listeners {
		reflect.ValueOf(c).Elem().Send(reflect.ValueOf(v))
	}
	return
}

// Clear does todo
func (e *Event) Clear() {
	// TODO: is it right?
	for _, c := range e.listeners {
		reflect.ValueOf(c).Close()
	}
	e.listeners = make([]interface{}, 0)
}
