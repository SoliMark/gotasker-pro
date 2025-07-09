package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	secretKey := "test_secret_key"
	jwtMaker := NewJWTMaker(secretKey)

	userID := uint(123)
	duration := time.Minute

	token, err := jwtMaker.GenerateToken(userID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := jwtMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotNil(t, claims)
	require.Equal(t, userID, claims.UserID)

	require.WithinDuration(t,
		time.Now().Add(duration),
		claims.ExpiresAt.Time,
		time.Second*2,
	)
}

func TestExpiredToken(t *testing.T) {

	secretKey := "test_secret_key"
	jwtMaker := NewJWTMaker(secretKey)

	userID := uint(123)
	duration := -time.Minute

	token, err := jwtMaker.GenerateToken(userID, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := jwtMaker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, claims)
}
