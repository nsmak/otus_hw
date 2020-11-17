package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Slices struct {
		Strings []string `validate:"len:5|in:first,12345|min:1"`
		Ints    []int    `validate:"min:1|max:50|in:25,30,50|len:40"`
	}

	PrivateWithTags struct {
		field1 string `validate:"in:private,life"`
		field2 int    `validate:"min:18|max:50"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in             interface{}
		expectedFields []string
		expectedErrs   []error
	}{
		{
			in: App{
				Version: "12345",
			},
		},
		{
			in: Slices{
				Strings: []string{"12345", "first"},
				Ints:    []int{25, 30, 50},
			},
		},
		{
			in: Response{
				Code: 404,
				Body: "Not found",
			},
		},
		{
			in: PrivateWithTags{
				field1: "LoLKek",
				field2: 100500,
			},
		},
		{
			in: Token{},
		},

		{
			in: App{
				Version: "3.3.5a",
			},
			expectedFields: []string{
				"Version",
			},
			expectedErrs: []error{
				ErrInvalidLength,
			},
		},
		{
			in: Slices{
				Strings: []string{"1234567", "first"},
				Ints:    []int{100, 40},
			},
			expectedFields: []string{
				"Strings",
				"Strings",
				"Ints",
				"Ints",
				"Ints",
			},
			expectedErrs: []error{
				ErrInvalidLength,
				ErrNotIncludedInSet,
				ErrMaxMoreMax,
				ErrNotIncludedInSet,
				ErrNotIncludedInSet,
			},
		},
		{
			in: User{
				Age:   14,
				Email: "test@test.com",
			},
			expectedFields: []string{
				"ID",
				"Age",
				"Role",
			},
			expectedErrs: []error{
				ErrInvalidLength,
				ErrLessThanMin,
				ErrNotIncludedInSet,
			},
		},
		{
			in: nil,
			expectedFields: []string{
				"-",
			},
			expectedErrs: []error{
				ErrInputIsNotStruct,
			},
		},
		{
			in: "1",
			expectedFields: []string{
				"-",
			},
			expectedErrs: []error{
				ErrInputIsNotStruct,
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			errs := unwrapErrors(err)

			require.Len(t, errs, len(tt.expectedErrs))
			if len(errs) > 0 {
				for i, err := range errs {
					require.Equal(t, err.Field, tt.expectedFields[i])
					require.True(t, errors.Is(err.Err, tt.expectedErrs[i]))
				}
			}
		})
	}
}

func unwrapErrors(err error) ValidationErrors {
	var vErrs ValidationErrors
	if !errors.As(err, &vErrs) {
		return nil
	}
	return vErrs
}
