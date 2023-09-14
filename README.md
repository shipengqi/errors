# errors

Package errors provides simple error handling primitives.
Migrated from [golib](https://github.com/shipengqi/golib).

Based on [github.com/pkg/errors](https://github.com/pkg/errors), and fully compatible with `github.com/pkg/errors`.

[![test](https://github.com/shipengqi/errors/actions/workflows/test.yml/badge.svg)](https://github.com/shipengqi/errors/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/shipengqi/errors/branch/main/graph/badge.svg?token=SMU4SI304O)](https://codecov.io/gh/shipengqi/errors)
[![Go Report Card](https://goreportcard.com/badge/github.com/shipengqi/errors)](https://goreportcard.com/report/github.com/shipengqi/errors)
[![release](https://img.shields.io/github/release/shipengqi/errors.svg)](https://github.com/shipengqi/errors/releases)
[![license](https://img.shields.io/github/license/shipengqi/errors)](https://github.com/shipengqi/errors/blob/main/LICENSE)

## Getting Started

### Coder

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/shipengqi/errors"
)

type fakeCoder struct {
	code   int
	status int
	msg    string
	ref    string
}

func (d fakeCoder) Code() int         { return d.code }
func (d fakeCoder) String() string    { return d.msg }
func (d fakeCoder) Reference() string { return d.ref }
func (d fakeCoder) HTTPStatus() int {
	if d.status == 0 {
		return http.StatusInternalServerError
	}
	return d.status
}

type parseCoder struct {
	code int
}

func (d parseCoder) Code() int { return d.code }

func main() {
	// annotates err with a code.
	codeErr := errors.WithCode(fmt.Errorf("demo error"), 20010)

	// reports whether any error in err's contains the given code.
	fmt.Println(errors.IsCode(codeErr, 20010)) // true
	
	// returns a code error with the format specifier.
	_ = errors.Codef(20011, "codef %s", "demo")
	// returns an error annotating err with a code and a stack trace at the point WrapC is called, and the format specifier.
	_ = errors.WrapC(fmt.Errorf("wrap error"), 20012, "wrap %s", "demo")
	

	demoCoder := fakeCoder{
		code:   20013,
		status: http.StatusBadRequest,
		msg:    "bad request",
		ref:    "https://docs.example.com/codes",
	}
	
	// registers a Coder to the global cache.
	errors.Register(demoCoder)
	
	// parse any error into icoder interface, find the corresponding Coder from global cache.
	errors.ParseCoder(parseCoder{code: 20013})
}

```

### Aggregate

```go
package main

import (
	"fmt"
	
	"github.com/shipengqi/errors"
)

func main() {
	// Aggregate represents an object that contains multiple errors, but does not 
	// necessarily have singular semantic meaning
	var errs []error
	errs = append(errs, 
		errors.New("error 1"), 
		errors.New("error 2"), 
		errors.New("error 3"),
	)

	agge := errors.NewAggregate(errs)
	fmt.Println(agge.Error()) // [error 1, error 2, error 3]
}
```

## Documentation

You can find the docs at [go docs](https://pkg.go.dev/github.com/shipengqi/errors).

## 🔋 JetBrains OS licenses

`errors` had been being developed with **GoLand** under the **free JetBrains Open Source license(s)** granted by JetBrains s.r.o., hence I would like to express my thanks here.

<a href="https://www.jetbrains.com/?from=errors" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg" alt="JetBrains Logo (Main) logo." width="250" align="middle"></a>
