package ssgo

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	AuthorizationHeaderKey = "Authorization"

	DefaultJWTTimeToLive = time.Hour
)

const (
	PermissionWrite Permission = "write"
	PermissionRead  Permission = "read"
)

type (
	// Permission for entity
	Permission string

	// Claims to store permissions per entity
	Claims struct {
		jwt.RegisteredClaims

		Permissions map[string][]Permission `json:"permissions"`
	}
)

// NewRegisteredClaims creates default jwt claims with predefined constraints
func NewRegisteredClaims(userID uuid.UUID, issuedAt time.Time, ttl time.Duration) jwt.RegisteredClaims {
	issuedAt = issuedAt.UTC()

	return jwt.RegisteredClaims{
		ID:        uuid.NewString(),
		Issuer:    "ss-go",
		Subject:   userID.String(),
		Audience:  jwt.ClaimStrings{"ss-go"},
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(issuedAt.Add(ttl)),
	}
}

// NewClaims constructor
func NewClaims(userID uuid.UUID, permissions map[string][]Permission) Claims {
	return Claims{
		RegisteredClaims: NewRegisteredClaims(userID, time.Now(), DefaultJWTTimeToLive),
		Permissions:      permissions,
	}
}

// NewSignedJWT returns raw jwt token signed
func NewSignedJWT(secret []byte, userID uuid.UUID, permissions map[string][]Permission) (string, error) {
	claims := NewClaims(userID, permissions)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ClaimsFromJWT parses raw jwt token and returns Claims
func ClaimsFromJWT(rawToken string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(rawToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("jwt is invalid: %s", rawToken)
}
