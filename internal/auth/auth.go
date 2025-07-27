package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// BasicAuth holds the username and password for basic authentication
type BasicAuth struct {
	Username string
	Password string
}

// Authenticator manages authentication
type Authenticator struct {
	users map[string]string // username -> password
}

// NewAuthenticator creates a new authenticator with predefined users
func NewAuthenticator() *Authenticator {
	// In a real application, these would come from a database or config
	users := map[string]string{
		"admin":  "password123",
		"client": "client456",
		"test":   "test789",
	}
	return &Authenticator{users: users}
}

// AddUser adds a new user to the authenticator
func (a *Authenticator) AddUser(username, password string) {
	a.users[username] = password
}

// ValidateCredentials checks if the username and password are valid
func (a *Authenticator) ValidateCredentials(username, password string) bool {
	storedPassword, exists := a.users[username]
	return exists && storedPassword == password
}

// UnaryInterceptor returns a gRPC unary server interceptor for basic authentication
func (a *Authenticator) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip authentication for health checks or specific methods if needed
		if strings.HasSuffix(info.FullMethod, "/Health") {
			return handler(ctx, req)
		}

		err := a.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// StreamInterceptor returns a gRPC stream server interceptor for basic authentication
func (a *Authenticator) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := a.authenticate(stream.Context())
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

// authenticate extracts and validates credentials from the gRPC metadata
func (a *Authenticator) authenticate(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return status.Error(codes.Unauthenticated, "missing authorization header")
	}

	authHeader := authHeaders[0]
	if !strings.HasPrefix(authHeader, "Basic ") {
		return status.Error(codes.Unauthenticated, "invalid authorization header format")
	}

	// Extract base64 encoded credentials
	encodedCreds := strings.TrimPrefix(authHeader, "Basic ")
	decodedCreds, err := base64.StdEncoding.DecodeString(encodedCreds)
	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid base64 encoding")
	}

	// Parse username:password
	credentials := string(decodedCreds)
	parts := strings.SplitN(credentials, ":", 2)
	if len(parts) != 2 {
		return status.Error(codes.Unauthenticated, "invalid credentials format")
	}

	username, password := parts[0], parts[1]

	// Validate credentials
	if !a.ValidateCredentials(username, password) {
		return status.Error(codes.Unauthenticated, "invalid username or password")
	}

	return nil
}

// EncodeBasicAuth encodes username and password for basic auth header
func EncodeBasicAuth(username, password string) string {
	credentials := fmt.Sprintf("%s:%s", username, password)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return fmt.Sprintf("Basic %s", encoded)
}
