// pkg/fcm/jwt.go — RS256 service account JWT builder for Google OAuth2.
//
// Builds the signed JWT assertion used to exchange a service account key
// for a short-lived OAuth2 access token (RFC 7523 JWT bearer grant).

package fcm

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"
)

func buildServiceAccountJWT(clientEmail, privateKeyPEM, audience, scope string) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", errors.New("invalid PEM block in private key")
	}

	var privateKey *rsa.PrivateKey
	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("parse PKCS1 private key: %w", err)
		}
		privateKey = key
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("parse PKCS8 private key: %w", err)
		}
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return "", errors.New("PKCS8 key is not RSA")
		}
	default:
		return "", fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}

	now := time.Now()
	header := map[string]string{"alg": "RS256", "typ": "JWT"}
	claims := map[string]any{
		"iss":   clientEmail,
		"sub":   clientEmail,
		"aud":   audience,
		"scope": scope,
		"iat":   now.Unix(),
		"exp":   now.Add(time.Hour).Unix(),
	}

	headerJSON, _ := json.Marshal(header)
	claimsJSON, _ := json.Marshal(claims)

	headerEnc := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsEnc := base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := headerEnc + "." + claimsEnc

	hash := sha256.New()
	hash.Write([]byte(signingInput))
	digest := hash.Sum(nil)

	sig, err := rsa.SignPKCS1v15(rand.Reader, privateKey, 0, digest)
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}
	sigEnc := base64.RawURLEncoding.EncodeToString(sig)

	return strings.Join([]string{headerEnc, claimsEnc, sigEnc}, "."), nil
}
