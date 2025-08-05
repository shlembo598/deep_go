package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	data []error
}

func (e *MultiError) Error() string {
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("%v errors occured:\n", len(e.data)))

	for _, v := range e.data {
		if e != nil {
			builder.WriteString(fmt.Sprintf("\t* %v", v.Error()))
		}
	}

	builder.WriteString("\n")

	return builder.String()

}

func Append(err error, errs ...error) *MultiError {
	newMultiError := &MultiError{}

	if err != nil {
		if multiErr, ok := err.(*MultiError); ok {
			newMultiError.data = append(newMultiError.data, multiErr.data...)
		} else {
			newMultiError.data = append(newMultiError.data, err)
		}

	}

	if len(errs) > 0 {
		newMultiError.data = append(newMultiError.data, errs...)
	}

	return newMultiError
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
