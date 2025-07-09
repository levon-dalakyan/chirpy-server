package tests

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/levon-dalakyan/chirpy-server/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestMakeAndValidateJWT(t *testing.T) {
	// Setup
	secret := "super-secret-key"
	userID := uuid.New()
	expiresIn := 1 * time.Hour

	// Generate token
	token, err := auth.MakeJWT(userID, secret, expiresIn)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate token
	returnedID, err := auth.ValidateJWT(token, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, returnedID)
}

func TestValidateJWT_InvalidSecret(t *testing.T) {
	secret := "correct-secret"
	badSecret := "wrong-secret"
	userID := uuid.New()
	expiresIn := 1 * time.Hour

	token, err := auth.MakeJWT(userID, secret, expiresIn)
	assert.NoError(t, err)

	_, err = auth.ValidateJWT(token, badSecret)
	assert.Error(t, err)
}

func TestValidateJWT_Expired(t *testing.T) {
	secret := "super-secret-key"
	userID := uuid.New()

	// Token that expired 1 second ago
	token, err := auth.MakeJWT(userID, secret, -1*time.Second)
	assert.NoError(t, err)

	_, err = auth.ValidateJWT(token, secret)
	assert.Error(t, err)
}
