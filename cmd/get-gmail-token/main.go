// cmd/get-gmail-token — one-time OAuth2 flow to obtain a Gmail refresh token.
//
// Run this once to get GMAIL_REFRESH_TOKEN. After that Emily Prime handles
// token refresh automatically at runtime — no browser needed again.
//
// Usage:
//   GMAIL_CLIENT_ID=... GMAIL_CLIENT_SECRET=... go run ./cmd/get-gmail-token
//
// It will print a URL. Open it in a browser, sign in as the CEO account,
// grant access, copy the auth code from the redirect URL, paste it here.
// The refresh token is printed — save it as GMAIL_REFRESH_TOKEN.

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	authURL  = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenURL = "https://oauth2.googleapis.com/token"

	// Scopes required by Emily Prime's gmail.go
	scopeRead = "https://www.googleapis.com/auth/gmail.readonly"
	scopeSend = "https://www.googleapis.com/auth/gmail.send"
)

func main() {
	clientID := os.Getenv("GMAIL_CLIENT_ID")
	clientSecret := os.Getenv("GMAIL_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		fmt.Fprintln(os.Stderr, "error: GMAIL_CLIENT_ID and GMAIL_CLIENT_SECRET must be set")
		os.Exit(1)
	}

	// Step 1 — print the authorization URL for the user to visit.
	params := url.Values{
		"client_id":     {clientID},
		"redirect_uri":  {"urn:ietf:wg:oauth:2.0:oob"},
		"response_type": {"code"},
		"scope":         {scopeRead + " " + scopeSend},
		"access_type":   {"offline"},
		"prompt":        {"consent"}, // force consent screen so refresh_token is always issued
	}
	authLink := authURL + "?" + params.Encode()

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  Emily Prime — Gmail OAuth2 Setup")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("1. Open this URL in your browser and sign in as the CEO account:")
	fmt.Println()
	fmt.Println(" ", authLink)
	fmt.Println()
	fmt.Println("2. Grant access, then copy the authorisation code shown.")
	fmt.Println()
	fmt.Print("Paste the code here: ")

	reader := bufio.NewReader(os.Stdin)
	code, _ := reader.ReadString('\n')
	code = strings.TrimSpace(code)
	if code == "" {
		fmt.Fprintln(os.Stderr, "error: no code entered")
		os.Exit(1)
	}

	// Step 2 — exchange the auth code for tokens.
	resp, err := http.PostForm(tokenURL, url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"code":          {code},
		"redirect_uri":  {"urn:ietf:wg:oauth:2.0:oob"},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error exchanging code: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "token exchange failed (status %d):\n%s\n", resp.StatusCode, body)
		os.Exit(1)
	}

	var tok struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}
	if err := json.Unmarshal(body, &tok); err != nil {
		fmt.Fprintf(os.Stderr, "parse token response: %v\n%s\n", err, body)
		os.Exit(1)
	}
	if tok.RefreshToken == "" {
		fmt.Fprintln(os.Stderr, "error: no refresh_token in response — try re-running with a fresh browser session")
		fmt.Fprintf(os.Stderr, "full response: %s\n", body)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  Success! Save these as environment variables:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Printf("export GMAIL_CLIENT_ID=%s\n", clientID)
	fmt.Printf("export GMAIL_CLIENT_SECRET=%s\n", clientSecret)
	fmt.Printf("export GMAIL_REFRESH_TOKEN=%s\n", tok.RefreshToken)
	fmt.Printf("export GMAIL_CEO_ADDRESS=%s\n", "emilyspringerton@gmail.com")
	fmt.Println()
	fmt.Println("Then start Emily Prime — gmail_read_inbox and gmail_send_alert will be available.")
	fmt.Println()
}
