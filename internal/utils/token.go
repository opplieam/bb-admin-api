package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"aidanwoods.dev/go-paseto"
)

// GenerateToken return paseto token
func GenerateToken(expire time.Duration, userID int32) (string, error) {
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetIssuer("bb-admin")
	token.SetExpiration(time.Now().Add(expire))
	token.SetString("user_id", strconv.Itoa(int(userID)))

	key, err := paseto.V4SymmetricKeyFromHex(os.Getenv("TOKEN_ENCODED"))
	if err != nil {
		return "", fmt.Errorf("failed to generate symmetric key: %w", err)
	}
	encrypted := token.V4Encrypt(key, nil)

	return encrypted, nil
}

// VerifyToken Validate token and return nil if it successes
func VerifyToken(token string) (*paseto.Token, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.IssuedBy("bb-admin"))
	parser.AddRule(paseto.NotExpired())

	key, err := paseto.V4SymmetricKeyFromHex(os.Getenv("TOKEN_ENCODED"))
	if err != nil {
		return nil, fmt.Errorf("failed to generate symmetric key: %w", err)
	}
	parsedToken, err := parser.ParseV4Local(key, token, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	return parsedToken, nil
}
