package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

func verify(authHeader string) {
	var oauthToken oauth2.Token
	json.Unmarshal([]byte(authHeader), &oauthToken)
	fmt.Printf("token: %v\n", oauthToken)
	fmt.Println("oauth2 token is valid:", oauthToken.Valid())

	ctx := context.Background()
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: customTransport}
	http.DefaultClient = client
	provider, err := oidc.NewProvider(ctx, "https://keycloak.dev.leetcyber.com/auth/realms/basic")
	if err != nil {
		log.Fatal(err)
	}
	userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(&oauthToken))
	if err != nil {
		log.Fatal(err)
		return
	}
	//resp := struct {
	//	OAuth2Token *oauth2.Token
	//	UserInfo    *oidc.UserInfo
	//}{&oauthToken, userInfo}
	//data, err := json.MarshalIndent(resp, "", "    ")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	fmt.Printf("UserInfo: %v\n", userInfo)

	// keycloak stuff
	//	ctx := context.Background()
	//	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	//	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	//	client := &http.Client{Transport: customTransport}
	//	http.DefaultClient = client
	//	provider, err := oidc.NewProvider(ctx, "https://keycloak.dev.leetcyber.com/auth/realms/basic")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	oidcConfig := &oidc.Config{
	//		ClientID: "client-secret",
	//	}
	//	verifier := provider.Verifier(oidcConfig)
	//	rawIDToken, ok := oauthToken.Extra("id_token").(string)
	//	if !ok {
	//		log.Fatal("No id_token field in oauth2 token.")
	//		return
	//	}
	//	idToken, err := verifier.Verify(ctx, rawIDToken)
	//	if err != nil {
	//		log.Fatal(err)
	//		return
	//	}
	//	fmt.Println("idToken:", idToken)
	//	fmt.Println("OIDC token provided to the server was validated successfully!")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)

			if name == "Authorization" {
				verify(h)
			}
		}
	}
}

func main() {

	http.HandleFunc("/", headers)

	http.ListenAndServe(":8090", nil)
}
