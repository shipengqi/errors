//go:build go1.20
// +build go1.20

package errors

import (
	"reflect"
	"testing"
)

func TestJoinReturnsNil(t *testing.T) {
	if err := Join(); err != nil {
		t.Errorf("errors.Join() = %v, want nil", err)
	}
	if err := Join(nil); err != nil {
		t.Errorf("errors.Join(nil) = %v, want nil", err)
	}
	if err := Join(nil, nil); err != nil {
		t.Errorf("errors.Join(nil, nil) = %v, want nil", err)
	}
}

func TestJoin(t *testing.T) {
	err := New("test")
	err1 := New("err1")
	err2 := New("err2")
	err3 := WithStack(err)
	err4 := WithMessage(err, "msg")
	err5 := WithCode(err, 99)
	for _, test := range []struct {
		errs []error
		want []error
	}{{
		errs: []error{err1},
		want: []error{err1},
	}, {
		errs: []error{err1, err2},
		want: []error{err1, err2},
	}, {
		errs: []error{err1, nil, err2},
		want: []error{err1, err2},
	}, {
		errs: []error{err1, err2, err3, err4, err5},
		want: []error{err1, err2, err3, err4, err5},
	}} {
		got := Join(test.errs...).(interface{ Unwrap() []error }).Unwrap()
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Join(%v) = %v; want %v", test.errs, got, test.want)
		}
		if len(got) != cap(got) {
			t.Errorf("Join(%v) returns errors with len=%v, cap=%v; want len==cap", test.errs, len(got), cap(got))
		}
	}
}

func TestJoinErrorMethod(t *testing.T) {
	err1 := New("err1")
	err2 := New("err2")
	for _, test := range []struct {
		errs []error
		want string
	}{{
		errs: []error{err1},
		want: "err1",
	}, {
		errs: []error{err1, err2},
		want: "err1\nerr2",
	}, {
		errs: []error{err1, nil, err2},
		want: "err1\nerr2",
	}} {
		got := Join(test.errs...).Error()
		if got != test.want {
			t.Errorf("Join(%v).Error() = %q; want %q", test.errs, got, test.want)
		}
	}
}
