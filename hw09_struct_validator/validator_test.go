package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	// "github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "111111111111111111111111111111111111",
				Name:   "Some name",
				Age:    20,
				Email:  "mail@mail.com",
				Role:   "admin",
				Phones: []string{"89995556677", "11122233344"},
				meta:   json.RawMessage(`{"data": "skipped"}`),
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "1",
				Name:   "Some name",
				Age:    100,
				Email:  "not-an-email",
				Role:   "user",
				Phones: []string{"12345"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("value does not satisfy validation rule: length != 36")},
				{Field: "Age", Err: fmt.Errorf("value does not satisfy validation rule: 100 > 50")},
				{Field: "Email", Err: fmt.Errorf("value does not satisfy validation rule: invalid regexp '^\\w+@\\w+\\.\\w+$'")},
				{Field: "Role", Err: fmt.Errorf("value does not satisfy validation rule: value not in [admin stuff]")},
				{Field: "Phones", Err: fmt.Errorf("value does not satisfy validation rule: length != 11")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			result := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, result)
				return
			}

			// require.Equal(t, len(tt.expectedErr), len(result))

			// for i := range result {
			// 	require.Equal(t, tt.expectedErr[i].Field, result[i].Field)
			// 	require.Equal(t, tt.expectedErr[i].Err.Error(), result[i].Err.Error())
			// }
			// Place your code here.

			//=====================================
			// if tt.expectedErr != nil {
			// 	if result == nil {
			// 		t.Errorf("expected error, got nil")
			// 		return
			// 	}
			var ve ValidationErrors
			if !errors.As(result, &ve) {
				t.Error(result.Error())
				return
			}
			// check specific fields
			for i, ve := range ve {
				if ve.Field == tt.expectedErr.(ValidationErrors)[i].Field {
					return
				}
			}
			t.Errorf("expected error with field %s", tt.expectedErr)
			// } else {
			// 	if result != nil {
			// 		t.Errorf("expected nil error, got %v", result)
			// 	}
			// }
			//=====================================
			_ = tt
		})
	}
}
