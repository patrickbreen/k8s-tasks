package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

func getBearerToken() []byte {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://keycloak.dev.leetcyber.com/auth/realms/basic")
	if err != nil {
		log.Fatal(err)
	}
	clientID := "client-secret"
	clientSecret := "client-secret"
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://oauth.dev.leetcyber.com/callbacks/redirect",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	oauth2Token, err := config.PasswordCredentialsToken(ctx, "patrick", "star")
	if err != nil {
		panic(err)
	}
	// maybe this auto-refreshes the token?
	config.Client(ctx, oauth2Token)

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		panic("couldn't get the rawIDToken")
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		panic(err)
	}

	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{oauth2Token, new(json.RawMessage)}

	if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(resp, "", "    ")
	fmt.Println("data:", string(data))
	if err != nil {
		panic(err)
	}
	//
	return []byte(rawIDToken)
}

func main() {
	// if this was for real, I'd trust the self signed cert ;)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	http.DefaultClient = client

	tokenData := getBearerToken()

	//TODO if token about to expire renew it, if renewal token about to expire, get new token

	// send request to server
	req, err := http.NewRequest("GET", "http://oauth.dev.leetcyber.com:8090/", bytes.NewBuffer(nil))
	req.Header.Set("Id", string(tokenData))
	req.Header.Add("Accept", "application/json")
	if err != nil {
		panic(err)
	}
	response, err := client.Do(req)
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if http.StatusOK != response.StatusCode {
		panic(err)
	}
	fmt.Println("response:", string(buf.Bytes()))
}
