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
		in          interface{}
		expectedErr error
	}{
		{
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
		},
		{
			in: Slices{
				Strings: []string{"12345", "first"},
				Ints:    []int{25, 30, 50},
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 404,
				Body: "Not found",
			},
			expectedErr: nil,
		},
		{
			in: PrivateWithTags{
				field1: "LoLKek",
				field2: 100500,
			},
			expectedErr: nil,
		},
		{
			in:          Token{},
			expectedErr: nil,
		},

		{
			in: App{
				Version: "3.3.5a",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   errors.New("length is invalid need: 5, now: 6"),
				},
			},
		},
		{
			in: Slices{
				Strings: []string{"1234567", "first"},
				Ints:    []int{100, 40},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Strings",
					Err:   errors.New("length is invalid need: 5, now: 7"),
				},
				ValidationError{
					Field: "Strings",
					Err:   errors.New("value (1234567) not included in validation set ([first 12345])"),
				},
				ValidationError{
					Field: "Ints",
					Err:   errors.New("value (100) is more thаn condition (50)"),
				},
				ValidationError{
					Field: "Ints",
					Err:   errors.New("value (100) not included in validation set ([25 30 50])"),
				},
				ValidationError{
					Field: "Ints",
					Err:   errors.New("value (40) not included in validation set ([25 30 50])"),
				},
			},
		},
		{
			in: User{
				Age:   14,
				Email: "test@test.com",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   errors.New("length is invalid need: 36, now: 0"),
				},
				ValidationError{
					Field: "Age",
					Err:   errors.New("value (14) is less thаn condition (18)"),
				},
				ValidationError{
					Field: "Role",
					Err:   errors.New("value () not included in validation set ([admin stuff])"),
				},
			},
		},
		{
			in:          nil,
			expectedErr: errors.New("input value is not a struct"),
		},
		{
			in:          "1",
			expectedErr: errors.New("input value is not a struct"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)

			require.Equal(t, tt.expectedErr, err)
		})
	}
}
