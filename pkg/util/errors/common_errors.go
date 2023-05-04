package errors

import (
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/pkg/errors"
)

var (
	_messages atomic.Value         // NOTE: stored map[string]map[int]string
	_codes    = map[int]struct{}{} // register codes.
)

type Code int

func (c Code) Error() string {
	return strconv.FormatInt(int64(c), 10)
}

func (c Code) Code() int {
	return int(c)
}

func (c Code) Message() string {
	if cm, ok := _messages.Load().(map[int]string); ok {
		if msg, ok := cm[c.Code()]; ok {
			return msg
		}
	}
	return c.Error()
}

func (c Code) Details() []interface{} {
	return nil
}

func (c Code) Equal(err error) bool {
	return EqualError(c, err)
}

type Codes interface {
	// sometimes Error return Code in string form
	// NOTE: don't use Error in monitor report even it also work for now
	Error() string
	// Code get error code.
	Code() int
	// Message get code message.
	Message() string
	//Detail get error detail,it may be nil.
	Details() []interface{}
	// Equal for compatible.
	// Deprecated: please use ecode.EqualError.
	Equal(error) bool
}

func String(e string) Code {
	if e == "" {
		return ErrCodeNone
	}
	// try error string
	i, err := strconv.Atoi(e)
	if err != nil {
		return ParameterTransferErr
	}
	return Code(i)
}

// Cause cause from error to ecode.
func Cause(e error) Codes {
	if e == nil {
		return ErrCodeNone
	}
	ec, ok := errors.Cause(e).(Codes)
	if ok {
		return ec
	}
	return String(e.Error())
}

// Equal equal a and b by code int.
func Equal(a, b Codes) bool {
	if a == nil {
		a = ErrCodeNone
	}
	if b == nil {
		b = ErrCodeNone
	}
	return a.Code() == b.Code()
}

// EqualError equal error
func EqualError(code Codes, err error) bool {
	return Cause(err).Code() == code.Code()
}

func New(e int, msg string) Code {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[e] = struct{}{}
	messages, ok := _messages.Load().(map[int]string)
	if ok {
		messages[e] = msg
		_messages.Store(messages)
	} else {
		m := make(map[int]string)
		m[e] = msg
		_messages.Store(m)
	}
	return Int(e)
}

func Int(e int) Code {
	return Code(e)
}
