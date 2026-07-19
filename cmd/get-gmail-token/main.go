// cmd/get-gmail-token — one-time OAuth2 flow to obtain a Gmail refresh token.
//
// Run this once to get GMAIL_REFRESH_TOKEN. After that Emily Prime handles
// token refresh automatically at runtime — no browser needed again.
//
// Usage:
//
//	GMAIL_CLIENT_ID=... GMAIL_CLIENT_SECRET=... go run ./cmd/get-gmail-token
//
// It starts a local HTTP server on 127.0.0.1, prints a URL to open in a
// browser, and captures the authorization code automatically from the
// redirect — no manual copy/paste of a code. Uses a loopback redirect_uri
// (http://127.0.0.1:<port>/callback), not Google's out-of-band
// ("urn:ietf:wg:oauth:2.0:oob") flow: Google stopped issuing OOB support to
// any OAuth client created after Feb 28, 2022, so that flow fails outright
// against a client made today. Loopback redirects on "Desktop app"-type
// clients remain fully supported and don't need the exact port
// pre-registered in Google Cloud Console.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
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

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not open a local port: %v\n", err)
		os.Exit(1)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d/callback", port)

	params := url.Values{
		"client_id":     {clientID},
		"redirect_uri":  {redirectURI},
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
	fmt.Println("2. Grant access. The browser will redirect to 127.0.0.1 — that's this")
	fmt.Println("   program listening locally, not a real website; the redirect will hang")
	fmt.Println("   or show a blank page in the browser tab, which is expected. Come back")
	fmt.Println("   to this terminal.")
	fmt.Println()
	fmt.Println("Waiting for the redirect (5 minute timeout)...")

	codeCh := make(chan string, 1)
	errCh := make(chan string, 1)
	srv := &http.Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if errParam := r.URL.Query().Get("error"); errParam != "" {
			fmt.Fprintln(w, "Authorization failed — you can close this tab and check the terminal.")
			errCh <- errParam
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "missing code", http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, "Authorization received — you can close this tab and go back to the terminal.")
		codeCh <- code
	})
	srv.Handler = mux
	go srv.Serve(listener) //nolint:errcheck

	var code string
	select {
	case code = <-codeCh:
	case errParam := <-errCh:
		fmt.Fprintf(os.Stderr, "\nerror: Google returned an authorization error: %s\n", errParam)
		os.Exit(1)
	case <-time.After(5 * time.Minute):
		fmt.Fprintln(os.Stderr, "\nerror: timed out waiting for authorization — no redirect received in 5 minutes")
		os.Exit(1)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx) //nolint:errcheck

	// Exchange the auth code for tokens.
	resp, err := http.PostForm(tokenURL, url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"code":          {code},
		"redirect_uri":  {redirectURI},
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
