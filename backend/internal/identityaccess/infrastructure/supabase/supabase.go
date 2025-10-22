package supabase

import (
	"context"
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
	observability "github.com/context-space/cloud-observability"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// SupabaseConfig holds configuration for Supabase
type SupabaseConfig struct {
	ProjectRef  string
	ServiceRole string
	JWTSecret   string
}

// Claims represents the JWT claims structure from Supabase Auth
type Claims struct {
	jwt.RegisteredClaims
	Role         string                 `json:"role"`
	Email        string                 `json:"email,omitempty"`
	Phone        string                 `json:"phone,omitempty"`
	AppMetadata  map[string]interface{} `json:"app_metadata,omitempty"`
	UserMetadata map[string]interface{} `json:"user_metadata,omitempty"`
	Aud          string                 `json:"aud,omitempty"`
}

// UserInfo represents a user info response from Supabase
type UserInfo struct {
	ID           string                 `json:"id"`
	Email        string                 `json:"email"`
	Phone        string                 `json:"phone"`
	IsAnonymous  bool                   `json:"is_anonymous"`
	InfoMetadata map[string]interface{} `json:"info_metadata"`
}

// SupabaseAuthService provides Supabase authentication capabilities
type SupabaseAuthService struct {
	client *Client
	config *SupabaseConfig
	obs    *observability.ObservabilityProvider
}

// NewSupabaseAuthService creates a new Supabase auth service
func NewSupabaseAuthService(config *SupabaseConfig, observabilityProvider *observability.ObservabilityProvider) (*SupabaseAuthService, error) {
	if config == nil {
		return nil, errors.New("supabase config cannot be nil")
	}

	client := NewClient(config.ProjectRef, config.ServiceRole)

	return &SupabaseAuthService{
		client: client,
		config: config,
		obs:    observabilityProvider,
	}, nil
}

// ValidateToken validates a JWT token from Supabase Auth
func (s *SupabaseAuthService) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	// Parse the token with claims validation
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing algorithm is HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Use the configured JWT secret
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	// Verify the issuer is from your Supabase project
	expectedIssuer := fmt.Sprintf("https://%s.supabase.co/auth/v1", s.config.ProjectRef)
	if claims.Issuer != expectedIssuer {
		return nil, errors.New("invalid token issuer")
	}

	// Success - token is valid
	return claims, nil
}

// GetUserInfo extracts user information from an access token
func (s *SupabaseAuthService) GetUserInfo(ctx context.Context, supUserID string) (*UserInfo, error) {
	req := AdminGetUserRequest{UserID: uuid.MustParse(supUserID)}

	resp, err := s.client.AdminGetUser(req)
	if err != nil {
		return nil, err
	}

	data, err := sonic.Marshal(resp)
	if err != nil {
		return nil, err
	}

	infoMetadata := make(map[string]interface{})
	err = sonic.Unmarshal(data, &infoMetadata)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:           supUserID,
		Email:        resp.Email,
		Phone:        resp.Phone,
		IsAnonymous:  resp.IsAnonymous,
		InfoMetadata: infoMetadata,
	}, nil
}
