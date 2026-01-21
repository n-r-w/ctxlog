package serrors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type someStruct struct{}

func (s *someStruct) someStructErrorf(format string, a ...any) error {
	return Errorf(format, a...)
}

func TestErrorf(t *testing.T) {
	err := Errorf("someError: %w", fmt.Errorf("wrapped some error"))
	expectedErrorMessage := "{serrors.TestErrorf.17} someError: wrapped some error"
	require.Equal(t, expectedErrorMessage, err.Error())

	var testStruct someStruct
	err = testStruct.someStructErrorf("someError: %w", fmt.Errorf("wrapped some error"))
	expectedErrorMessage = "{serrors.(*someStruct).someStructErrorf.13} someError: wrapped some error"

	require.Equal(t, expectedErrorMessage, err.Error())
}

func TestNew(t *testing.T) {
	err := New("someError")
	expectedErrorMessage := "{serrors.TestNew.29} someError"
	require.Equal(t, expectedErrorMessage, err.Error())

	var testStruct someStruct
	err = testStruct.someStructErrorf("someError")
	expectedErrorMessage = "{serrors.(*someStruct).someStructErrorf.13} someError"

	require.Equal(t, expectedErrorMessage, err.Error())
}
