package api

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/sethvargo/go-password/password"
)

func GenerateSecret() (string, error) {
	codeVerifier, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		return "", errors.New("couldnt create a new secret")
	}

	entropy := sha256.Sum256([]byte(codeVerifier))
	return base64.StdEncoding.EncodeToString(entropy[:]), nil
}
