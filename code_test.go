package errors

import (
	"errors"
	"testing"
)

func TestRegister(t *testing.T) {
	mockSuccessCode := defaultCoder{
		code:   0,
		status: 200,
		msg:    "SUCCESS",
	}
	Register(mockSuccessCode)
	defer unregister(mockSuccessCode)

	if "SUCCESS" != mockSuccessCode.String() {
		t.Errorf("code string: want: %s, got: %s", "SUCCESS", mockSuccessCode.String())
	}
	if 200 != mockSuccessCode.HTTPStatus() {
		t.Errorf("code http status: want: %d, got: %d", 200, mockSuccessCode.HTTPStatus())
	}
	if "" != mockSuccessCode.Reference() {
		t.Errorf("code reference: want: %s, got: %s", "", mockSuccessCode.Reference())
	}

	t.Run("HTTP status 0", func(t *testing.T) {
		mockSuccessCode2 := defaultCoder{
			code: 3,
			msg:  "SUCCESS",
		}
		Register(mockSuccessCode2)
		defer unregister(mockSuccessCode2)
		if 500 != mockSuccessCode2.HTTPStatus() {
			t.Errorf("code http status: want: %d, got: %d", 500, mockSuccessCode2.HTTPStatus())
		}
	})
}

func TestRegisterPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			if err != "code `1` is reserved by `github.com/shipengqi/errors` as Unknown Code" {
				t.Errorf("code string: want: %s, got: %s", "code `1` is reserved by `github.com/shipengqi/errors` as Unknown Code", err)
			}
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

func TestRegisterPanic2(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			if err != "code `3` already registered" {
				t.Errorf("code string: want: %s, got: %s", "code `3` already registered", err)
			}
		} else {
			t.Fatal("no panic")
		}
	}()
	mockErrCode := defaultCoder{
		code:   3,
		status: 200,
		msg:    "error",
	}
	Register(mockErrCode)
	defer unregister(mockErrCode)
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
		if got != r.expected {
			t.Errorf("IsCode: want: %v, got: %v", r.expected, got)
		}
	}
}

func TestParseCoder(t *testing.T) {
	errUnknown := WithCode(errors.New(unknown), 1)
	err := ParseCoder(errUnknown)
	if err != unknownCode {
		t.Errorf("ParseCoder: want: unknown, got: %s", err)
	}

	err = ParseCoder(nil)
	if err != nil {
		t.Errorf("ParseCoder: want: nil, got: %s", err)
	}

	errUnknown2 := WithCode(errors.New(unknown), 2)
	err = ParseCoder(errUnknown2)
	if err != unknownCode {
		t.Errorf("ParseCoder: want: unknown, got: %s", err)
	}

	mockSuccessCode := defaultCoder{
		code:   10010,
		status: 200,
		msg:    "SUCCESS",
	}
	Register(mockSuccessCode)
	defer unregister(mockSuccessCode)

	embedErr1 := WithMessage(WithCode(errors.New("embedded"), 10010), "")
	err = ParseCoder(embedErr1)
	if err.Code() != 10010 {
		t.Errorf("ParseCoder: want: 2, got: %s", err)
	}
}

func TestParseCoderRecursively(t *testing.T) {
	codes := []Coder{
		defaultCoder{
			code:   10010,
			status: 200,
			msg:    "SUCCESS",
		},
		defaultCoder{
			code:   10011,
			status: 200,
			msg:    "SUCCESS",
		},
	}
	for _, v := range codes {
		Register(v)
	}

	defer func() {
		for _, v := range codes {
			unregister(v)
		}
	}()

	embedErr1 := WithMessage(WithCode(errors.New("embedded"), 10010), "embed1")
	err := ParseCoder(embedErr1)
	if err.Code() != 10010 {
		t.Errorf("ParseCoder: want: 2, got: %s", err)
	}

	embedErr2 := WithMessage(WithMessage(WithCode(errors.New("embedded"), 10010), "embed1"), "embed2")
	err = ParseCoder(embedErr2)
	if err.Code() != 10010 {
		t.Errorf("ParseCoder: want: 2, got: %s", err)
	}

	embedErr3 := WithMessage(WithCode(WithCode(errors.New("embedded"), 10010), 10011), "embed2")
	err = ParseCoder(embedErr3)
	if err.Code() != 10011 {
		t.Errorf("ParseCoder: want: 2, got: %s", err)
	}
}
