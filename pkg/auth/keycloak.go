package auth

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"os"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type AuthenticatorService interface {
	ExtractToken(context.Context) (string, error)
	ValidateTokenMiddleware(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error)
	VerifyToken(context.Context, string) (*oidc.IDToken, error)
}

type key int

const (
	// ContextKeyUser is used to store the authenticated user's claims in context.
	ContextKeyUser key = iota
)

// keycloakAuthenticator handles OpenID Connect token validation.
type keycloakAuthenticator struct {
	Verifier *oidc.IDTokenVerifier
	client   *oauth2.Config
}

func (a *keycloakAuthenticator) VerifyToken(ctx context.Context, s string) (*oidc.IDToken, error) {
	return a.Verifier.Verify(ctx, s)
}

// ValidateTokenMiddleware validates the JWT token in the authorization header.
func (a *keycloakAuthenticator) ValidateTokenMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Extract and validate the token from metadata (authorization header).
	token, err := a.ExtractToken(ctx)
	if err != nil {
		return nil, err
	}

	// Parse and verify the token.
	idToken, err := a.VerifyToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// Get the claims from the token.
	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %v", err)
	}

	// Pass the claims into the context for further use in the handler.
	ctx = context.WithValue(ctx, ContextKeyUser, claims)
	return handler(ctx, req)
}

// ExtractToken extracts the bearer token from the gRPC metadata (authorization header).
func (a *keycloakAuthenticator) ExtractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata")
	}

	// Look for the authorization header.
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return "", errors.New("missing authorization header")
	}

	// The authorization header should be in the form "Bearer <token>".
	parts := strings.SplitN(authHeader[0], " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization header")
	}

	return parts[1], nil
}

// NewAuthenticator creates a new OIDC authenticator using the given issuer URL and client configuration.
func NewAuthenticator() (AuthenticatorService, error) {
	issuerURL := os.Getenv("KC.ISSUER_URL")
	clientID := os.Getenv("KC.CLIENT_ID")

	provider, err := oidc.NewProvider(context.Background(), issuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	// Configure an IDTokenVerifier to validate tokens.
	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	// Create an OAuth2 client config.
	client := &oauth2.Config{
		ClientID: clientID,
	}

	return &keycloakAuthenticator{
		Verifier: verifier,
		client:   client,
	}, nil
}
