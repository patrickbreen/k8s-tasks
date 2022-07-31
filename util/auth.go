package util

import (
	"context"
	"net/http"

	"github.com/coreos/go-oidc"
)

type UserClaims struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	Id            string `json:"sub"`
	EmailVerified bool   `json:"email_verified"`
}

func getRawIDToken(req *http.Request) string {
	for name, headers := range req.Header {
		for _, h := range headers {
			if name == "Id" {
				return h
			}
		}
	}
	return ""
}

func ValidateAndGetClaims(req *http.Request) (UserClaims, error) {
	rawIDToken := getRawIDToken(req)
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, keycloakUrl())
	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return UserClaims{}, err
	}
	claims := UserClaims{}
	if err := idToken.Claims(&claims); err != nil {
		return UserClaims{}, err
	}
	return claims, nil
}
