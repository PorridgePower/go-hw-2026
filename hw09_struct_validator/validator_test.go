package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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
		Versions   []string `validate:"len:5"`
		Timestamps []int64  `validate:"min:1767225600"`
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
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "IncompatibleTagAndTypeError",
			in: struct {
				Field int `validate:"in:abc"`
			}{Field: 42},
			expectedErr: fmt.Errorf("%w: 'abc' is not a number", errTagParser),
		},
		{
			name: "MixedRules",
			in: struct {
				Field int `validate:"min:20|max:10"`
			}{Field: 25},
			expectedErr: ValidationErrors{
				{Field: "Field", Err: fmt.Errorf("value does not satisfy validation rule: 25 > 10")},
			},
		},
		{
			in: struct {
				Field int `validate:"min:20|max:10"`
			}{Field: 25},
			expectedErr: ValidationErrors{
				{Field: "Field", Err: fmt.Errorf("value does not satisfy validation rule: 25 > 10")},
			},
		},
		{
			name: "IncompatibleTagAndTypeSkipped",
			in: struct {
				Field int `validate:"regexp:.*"`
			}{Field: 42},
			expectedErr: nil,
		},
		{
			name: "TestOk",
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
			name: "TestNok",
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
				{
					Field: "Email",
					Err:   fmt.Errorf("value does not satisfy validation rule: value not match regexp '^\\w+@\\w+\\.\\w+$'"),
				},
				{Field: "Role", Err: fmt.Errorf("value does not satisfy validation rule: value not in [admin stuff]")},
				{Field: "Phones", Err: fmt.Errorf("value does not satisfy validation rule: length != 11")},
			},
		},
		{
			// collect both err from Versions+Timestamps
			name: "CollectErrors",
			in: App{
				Versions:   []string{"1.1"},
				Timestamps: []int64{1735689600},
			},
			expectedErr: ValidationErrors{
				{Field: "Versions", Err: fmt.Errorf("value does not satisfy validation rule: length != 5")},
				{Field: "Timestamps", Err: fmt.Errorf("value does not satisfy validation rule: 1735689600 < 1767225600")},
			},
		},
		{
			// stops at the first Versions array item
			name: "StopsAtFirstArrayItem",
			in: App{
				Versions:   []string{"1.1", "2.2"},
				Timestamps: []int64{1767225600},
			},
			expectedErr: ValidationErrors{
				{Field: "Versions", Err: fmt.Errorf("value does not satisfy validation rule: length != 5")},
			},
		},
		{
			name:        "SkipWithoutTagField",
			in:          Response{Code: 200, Body: ""},
			expectedErr: nil,
		},
		{
			name: "TestIn",
			in:   Response{Code: 999},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: fmt.Errorf("value does not satisfy validation rule: value not in [200 404 500]")},
			},
		},
		{
			name:        "UnsupportedType",
			in:          Token{Header: []byte("header_data"), Payload: []byte("payload_data")},
			expectedErr: nil,
		},
	}

	//nolint:errorlint
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d:%s", i, tt.name), func(t *testing.T) {
			tt := tt
			t.Parallel()
			result := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, result)
				return
			}

			var ve ValidationErrors
			if !errors.As(result, &ve) {
				require.ErrorIs(t, result, errTagParser)
				require.Equal(t, tt.expectedErr, result)
				return
			}

			require.Equal(t, len(ve), len(tt.expectedErr.(ValidationErrors)))
			for i, e := range ve {
				require.Equal(t, tt.expectedErr.(ValidationErrors)[i].Field, e.Field)
				require.Equal(t, tt.expectedErr.(ValidationErrors)[i].Err.Error(), e.Err.Error())
			}
			_ = tt
		})
	}
}

func TestValidateNonStruct(t *testing.T) {
	result := Validate("data")
	require.Equal(t, fmt.Errorf("v is not a struct"), result)
}

func TestParseTag(t *testing.T) {
	tests := []struct {
		name        string
		tag         string
		expectedErr string
	}{
		{
			name:        "EmptyRule",
			tag:         `validate:"len:"`,
			expectedErr: "invalid validate tag: empty rule 'len:'",
		},
		{
			name:        "IncompatibleRule1",
			tag:         `validate:"len:AAAA"`,
			expectedErr: "invalid validate tag: 'AAAA' is not a number",
		},
		{
			name:        "IncompatibleRule2",
			tag:         `validate:"min:abc"`,
			expectedErr: "invalid validate tag: 'abc' is not a number",
		},
		{
			name:        "InvalidRegexp",
			tag:         `validate:"regexp:[invalid"`,
			expectedErr: "invalid validate tag: invalid regexp '[invalid'",
		},
		{
			name:        "UnknownTag",
			tag:         `validate:"re:.*"`,
			expectedErr: "invalid validate tag: unknown tag value 're'",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("case %s", tt.name), func(t *testing.T) {
			t.Parallel()

			_, err := parseTag(reflect.StructTag(tt.tag))
			require.Error(t, err)
			require.ErrorIs(t, err, errTagParser)
			require.Equal(t, tt.expectedErr, err.Error())
		})
	}
}
