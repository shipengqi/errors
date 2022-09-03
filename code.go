package errors

import (
	"fmt"
	"net/http"
	"sync"
)

type Coder interface {
	// HTTPStatus that should be used for the associated error code.
	HTTPStatus() int

	// String error message.
	String() string

	// Reference returns the detail documents for user.
	Reference() string

	icoder
}

type icoder interface {
	// Code returns the code of the coder
	Code() int
}

type defaultCoder struct {
	code   int
	status int
	msg    string
	ref    string
}

func (d defaultCoder) Code() int         { return d.code }
func (d defaultCoder) String() string    { return d.msg }
func (d defaultCoder) Reference() string { return d.ref }
func (d defaultCoder) HTTPStatus() int {
	if d.status == 0 {
		return http.StatusInternalServerError
	}
	return d.status
}

type causer interface {
	Cause() error
}

var (
	unknownCode = defaultCoder{code: 1, status: http.StatusInternalServerError,
		msg: "Internal server error"}
	// _codes registered codes.
	_codes = make(map[int]Coder)
	mux    = &sync.Mutex{}
)

// Register registers an Coder.
func Register(code Coder) {
	if code.Code() == unknownCode.Code() {
		panic(fmt.Sprintf("code `%d` is reserved by `github.com/shipengqi/errors` as Unknown Code", code.Code()))
	}
	if _, ok := _codes[code.Code()]; ok {
		panic(fmt.Sprintf("code `%d` already registered", code.Code()))
	}
	mux.Lock()
	defer mux.Unlock()

	_codes[code.Code()] = code
}

// ParseCoder parse any error into icoder interface.
// nil error will return nil direct.
// None withStack error will be parsed as Unknown Code.
func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(icoder); ok {
		if coder, ok := _codes[v.Code()]; ok {
			return coder
		}
	}

	return unknownCode
}

// IsCode reports whether any error in err's contains the given code.
func IsCode(err error, code int) bool {
	if v, ok := err.(icoder); ok {
		if v.Code() == code {
			return true
		}
	}
	if v, ok := err.(causer); ok {
		err = v.Cause()
		return IsCode(err, code)
	}

	return false
}

func unregister(code Coder) {
	if _, ok := _codes[code.Code()]; ok {
		mux.Lock()
		defer mux.Unlock()

		delete(_codes, code.Code())
	}
}

func init() {
	_codes[unknownCode.Code()] = unknownCode
}
