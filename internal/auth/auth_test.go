package auth

import (
	"github.com/google/uuid"
	"net/http"
	"testing"
	"time"
)

func TestMakeJWTAndValidateJWT(t *testing.T) {
	// create test object
	secret := "test-object-secret"
	userID := uuid.New()
	expiresIn := time.Hour

	// test function
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// validate
	returnredID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// asset
	if returnredID != userID {
		t.Errorf("expected %v, got %v", userID, returnredID)
	}
}

func TestJWTEdgeCase(t *testing.T) {
	// create test object with expired duration
	secret := "test-object-secret"
	userID := uuid.New()
	expiresIn := -1 * time.Hour // expired an hour ago already

	// test Make function
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// validate
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("expected an error for an expired token, but got nil")
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		headerValue string
		expectToken string
		expectError bool
	}{
		{
			name:        "valid token",
			headerValue: "Bearer abc123",
			expectToken: "abc123",
			expectError: false,
		},
		{
			name:        "missing header",
			headerValue: "",
			expectError: true,
		},
		{
			name:        "wrong prefix",
			headerValue: "Token abc123",
			expectError: true,
		},
		{
			name:        "empty token",
			headerValue: "Bearer ",
			expectError: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			headers := http.Header{}

			if testCase.headerValue != "" {
				headers.Set("Authorization", testCase.headerValue)
			}

			token, err := GetBearerToken(headers)
			if testCase.expectError {
				if err == nil {
					t.Error("expect error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if token != testCase.expectToken {
				t.Errorf("expected token %s, got %s", testCase.expectToken, token)
			}
		})
	}
}
