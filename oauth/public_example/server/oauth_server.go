package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc"
)

func verify(rawIDToken string) {
	ctx := context.Background()
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: customTransport}
	http.DefaultClient = client
	provider, err := oidc.NewProvider(ctx, "https://keycloak.dev.leetcyber.com/auth/realms/basic")
	clientID := "client-secret"
	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		panic(err)
	}
	var claims struct {
		Email         string `json:"email"`
		Name          string `json:"name"`
		Id            string `json:"sid"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		panic(err)
	}
	fmt.Printf("idClaims: %v\n", claims)
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)

			if name == "Id" {
				verify(h)
			}
		}
	}
}

func main() {

	http.HandleFunc("/", headers)

	http.ListenAndServe(":8090", nil)
}
