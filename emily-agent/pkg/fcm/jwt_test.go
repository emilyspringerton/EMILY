package fcm

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"strings"
	"testing"
)

func testRSAKeyPEM(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}
	der := x509.MarshalPKCS1PrivateKey(key)
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: der}
	return string(pem.EncodeToMemory(block))
}

func TestBuildServiceAccountJWTStructure(t *testing.T) {
	pemKey := testRSAKeyPEM(t)
	token, err := buildServiceAccountJWT("svc@project.iam.gserviceaccount.com", pemKey,
		"https://oauth2.googleapis.com/token", "https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		t.Fatalf("buildServiceAccountJWT: %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("JWT must have 3 parts, got %d", len(parts))
	}

	// Header
	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatalf("decode header: %v", err)
	}
	var header map[string]string
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		t.Fatalf("unmarshal header: %v", err)
	}
	if header["alg"] != "RS256" {
		t.Errorf("alg = %q, want RS256", header["alg"])
	}
	if header["typ"] != "JWT" {
		t.Errorf("typ = %q, want JWT", header["typ"])
	}

	// Claims
	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("decode claims: %v", err)
	}
	var claims map[string]any
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		t.Fatalf("unmarshal claims: %v", err)
	}
	if claims["iss"] != "svc@project.iam.gserviceaccount.com" {
		t.Errorf("iss = %q", claims["iss"])
	}
	if claims["aud"] != "https://oauth2.googleapis.com/token" {
		t.Errorf("aud = %q", claims["aud"])
	}
	iat, ok1 := claims["iat"].(float64)
	exp, ok2 := claims["exp"].(float64)
	if !ok1 || !ok2 {
		t.Fatalf("iat/exp not present or not numeric")
	}
	if exp-iat != 3600 {
		t.Errorf("exp-iat = %v, want 3600", exp-iat)
	}
}

func TestBuildServiceAccountJWTInvalidPEM(t *testing.T) {
	_, err := buildServiceAccountJWT("svc@example.com", "not a PEM",
		"https://oauth2.googleapis.com/token", "")
	if err == nil {
		t.Fatal("expected error for invalid PEM, got nil")
	}
}

func TestBuildServiceAccountJWTUnsupportedKeyType(t *testing.T) {
	// Provide a PEM block with wrong type
	block := &pem.Block{Type: "CERTIFICATE", Bytes: []byte("dummy")}
	pem := string(pem.EncodeToMemory(block))
	_, err := buildServiceAccountJWT("svc@example.com", pem, "", "")
	if err == nil || !strings.Contains(err.Error(), "unsupported PEM block type") {
		t.Errorf("expected unsupported PEM error, got %v", err)
	}
}
