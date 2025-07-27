package auth

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestNewAuthenticator(t *testing.T) {
	auth := NewAuthenticator()

	if auth == nil {
		t.Fatal("NewAuthenticator() returned nil")
	}

	if len(auth.users) == 0 {
		t.Fatal("NewAuthenticator() should create default users")
	}
}

func TestValidateCredentials(t *testing.T) {
	auth := NewAuthenticator()

	tests := []struct {
		name     string
		username string
		password string
		expected bool
	}{
		{"valid admin", "admin", "password123", true},
		{"valid client", "client", "client456", true},
		{"valid test", "test", "test789", true},
		{"invalid username", "invalid", "password123", false},
		{"invalid password", "admin", "wrongpassword", false},
		{"empty username", "", "password123", false},
		{"empty password", "admin", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := auth.ValidateCredentials(tt.username, tt.password)
			if result != tt.expected {
				t.Errorf("ValidateCredentials(%q, %q) = %v, want %v",
					tt.username, tt.password, result, tt.expected)
			}
		})
	}
}

func TestAddUser(t *testing.T) {
	auth := NewAuthenticator()

	auth.AddUser("newuser", "newpassword")

	if !auth.ValidateCredentials("newuser", "newpassword") {
		t.Error("AddUser() should add user successfully")
	}
}

func TestEncodeBasicAuth(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		expected string
	}{
		{"admin credentials", "admin", "password123", "Basic YWRtaW46cGFzc3dvcmQxMjM="},
		{"client credentials", "client", "client456", "Basic Y2xpZW50OmNsaWVudDQ1Ng=="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EncodeBasicAuth(tt.username, tt.password)
			if result != tt.expected {
				t.Errorf("EncodeBasicAuth(%q, %q) = %q, want %q",
					tt.username, tt.password, result, tt.expected)
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {
	auth := NewAuthenticator()

	tests := []struct {
		name        string
		authHeader  string
		expectError bool
	}{
		{
			name:        "valid credentials",
			authHeader:  "Basic YWRtaW46cGFzc3dvcmQxMjM=", // admin:password123
			expectError: false,
		},
		{
			name:        "invalid credentials",
			authHeader:  "Basic aW52YWxpZDppbnZhbGlk", // invalid:invalid
			expectError: true,
		},
		{
			name:        "malformed header",
			authHeader:  "Bearer token",
			expectError: true,
		},
		{
			name:        "invalid base64",
			authHeader:  "Basic invalid-base64",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create context with metadata
			md := metadata.New(map[string]string{
				"authorization": tt.authHeader,
			})
			ctx := metadata.NewIncomingContext(context.Background(), md)

			err := auth.authenticate(ctx)

			if tt.expectError && err == nil {
				t.Error("authenticate() should return error but didn't")
			}
			if !tt.expectError && err != nil {
				t.Errorf("authenticate() should not return error but got: %v", err)
			}
		})
	}
}

func TestAuthenticateNoMetadata(t *testing.T) {
	auth := NewAuthenticator()

	// Test with context without metadata
	ctx := context.Background()
	err := auth.authenticate(ctx)

	if err == nil {
		t.Error("authenticate() should return error for missing metadata")
	}

	if status.Code(err).String() != "Unauthenticated" {
		t.Errorf("authenticate() should return Unauthenticated error, got: %v", err)
	}
}

func TestAuthenticateNoAuthHeader(t *testing.T) {
	auth := NewAuthenticator()

	// Test with context with metadata but no authorization header
	md := metadata.New(map[string]string{
		"other-header": "value",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	err := auth.authenticate(ctx)

	if err == nil {
		t.Error("authenticate() should return error for missing authorization header")
	}

	if status.Code(err).String() != "Unauthenticated" {
		t.Errorf("authenticate() should return Unauthenticated error, got: %v", err)
	}
}
