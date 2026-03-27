package common

import "testing"

func TestDTOSchemaValidation(t *testing.T) {
	type SignUpRequest struct {
		Username      string `json:"username" binding:"required"`
		Password      string `json:"password" binding:"required"`
		ApplicationID string `json:"application_id" binding:"required"`
		Nickname      string `json:"nickname"`
	}

	type NoJSONTagRequest struct {
		Name string `binding:"required"`
	}

	type NumericRequest struct {
		Age int `json:"age" binding:"required"`
	}

	type BoolRequest struct {
		Enabled bool `json:"enabled" binding:"required"`
	}

	tests := []struct {
		name    string
		input   any
		wantErr string
	}{
		{
			name: "valid struct",
			input: SignUpRequest{
				Username:      "alice",
				Password:      "secret",
				ApplicationID: "app-1",
			},
			wantErr: "",
		},
		{
			name: "missing required string field",
			input: SignUpRequest{
				Username:      "alice",
				Password:      "",
				ApplicationID: "app-1",
			},
			wantErr: `field 'password' is required`,
		},
		{
			name: "pointer to valid struct",
			input: &SignUpRequest{
				Username:      "alice",
				Password:      "secret",
				ApplicationID: "app-1",
			},
			wantErr: "",
		},
		{
			name:    "nil interface input",
			input:   nil,
			wantErr: "expected struct or pointer to struct",
		},
		{
			name: "nil pointer input",
			input: func() any {
				var req *SignUpRequest
				return req
			}(),
			wantErr: "nil pointer provided",
		},
		{
			name:    "non-struct input",
			input:   123,
			wantErr: "expected struct or pointer to struct",
		},
		{
			name:    "missing required field without json tag uses field name",
			input:   NoJSONTagRequest{},
			wantErr: `field 'Name' is required`,
		},
		{
			name: "non-required field is ignored",
			input: SignUpRequest{
				Username:      "alice",
				Password:      "secret",
				ApplicationID: "app-1",
				Nickname:      "",
			},
			wantErr: "",
		},
		{
			name:    "required int zero value fails",
			input:   NumericRequest{},
			wantErr: `field 'age' is required`,
		},
		{
			name:    "required bool zero value fails",
			input:   BoolRequest{},
			wantErr: `field 'enabled' is required`,
		},
		{
			name: "required int non-zero passes",
			input: NumericRequest{
				Age: 18,
			},
			wantErr: "",
		},
		{
			name: "required bool true passes",
			input: BoolRequest{
				Enabled: true,
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				err := DTOSchemaValidation(tt.input)

				if tt.wantErr == "" {
					if err != nil {
						t.Fatalf("expected no error, got %v", err)
					}
					return
				}

				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}

				if err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
				}
			},
		)
	}
}
