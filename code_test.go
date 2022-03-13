package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	mockSuccessCode := defaultCoder{
		code:   0,
		status: 200,
		msg:    "SUCCESS",
	}
	Register(mockSuccessCode)
	defer unregister(mockSuccessCode)

	assert.Equal(t, "SUCCESS", mockSuccessCode.String())
	assert.Equal(t, 200, mockSuccessCode.HTTPStatus())
	assert.Equal(t, "", mockSuccessCode.Reference())
}

func TestRegisterPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			assert.Equal(t, err, "code `1` is reserved by `github.com/shipengqi/errors` as Unknown Code")
		} else {
			t.Fatal("no panic")
		}
	}()
	mockErrCode := defaultCoder{
		code:   1,
		status: 200,
		msg:    "error",
	}
	Register(mockErrCode)
}

func TestIsCode(t *testing.T) {
	ok := WithCode(errors.New("ok"), 0)
	errUnknown := WithCode(errors.New(unknown), 1)
	type run struct {
		expected bool
		code     int
		err      error
	}
	runs := []run{
		{true, 0, ok},
		{false, 0, errUnknown},
		{true, 1, errUnknown},
		{true, 1, WithCode(errUnknown, 2)},
		{true, 1, WithCode(New("test1"), 1)},
		{true, 1, WithCode(WithMessage(New("test2"), "msg2"), 1)},
		{true, 1, WithMessage(errUnknown, "msg3")},
		{true, 1, WithMessage(WithCode(errUnknown, 2), "msg4")},
		{true, 1, Wrap(errUnknown, "msg5")},
		{true, 1, Wrap(WithCode(WithCode(errUnknown, 2), 3), "msg6")},
		{true, 2, Wrap(WithCode(WithCode(errUnknown, 2), 3), "msg7")},
		{true, 3, Wrap(WithCode(WithCode(errUnknown, 2), 3), "msg8")},
	}
	for _, r := range runs {
		got := IsCode(r.err, r.code)
		assert.Equal(t, got, r.expected, fmt.Sprintf("IsCode(%s, %d)", r.err.Error(), r.code))
	}
}

func TestParseCoder(t *testing.T) {
	errUnknown := WithCode(errors.New(unknown), 1)
	err := ParseCoder(errUnknown)
	assert.Equal(t, unknownCode, err)

	err = ParseCoder(nil)
	assert.Nil(t, err)

	errUnknown2 := WithCode(errors.New(unknown), 2)
	err = ParseCoder(errUnknown2)
	assert.Equal(t, unknownCode, err)
}

func TestCodef(t *testing.T) {
	err := Codef(3, "test codef")
	assert.Equal(t, "code: 3: test codef", err.Error())
}
