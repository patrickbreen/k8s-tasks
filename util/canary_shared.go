package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"leet/models"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var tokenSource oauth2.TokenSource = nil
var clientID = "client-secret"
var clientSecret = "client-secret"

// hard code all oauth/oidc stuff (for now) and return a valid TokenSource
func getTokenSource() oauth2.TokenSource {
	ctx := context.Background()
	if err != nil {
		log.Fatalln(err)
	}
	provider, err := oidc.NewProvider(ctx, "https://keycloak.dev.leetcyber.com/auth/realms/basic")
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://oauth.dev.leetcyber.com/callbacks/redirect",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	oauth2Token, err := config.PasswordCredentialsToken(ctx, "patrick", "star")
	if err != nil {
		log.Fatalln(err)
	}
	return config.TokenSource(ctx, oauth2Token)
}

// return a valid IDToken as a string
func getIDToken() string {
	// init
	if tokenSource == nil {
		tokenSource = getTokenSource()
	}
	// tokenSource.Token() will refresh the token as-needed
	oauth2Token, err := tokenSource.Token()
	if err != nil {
		log.Fatalln(err)
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Fatalln("couldn't get the rawIDToken")
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://keycloak.dev.leetcyber.com/auth/realms/basic")
	verifier := provider.Verifier(oidcConfig)
	_, err = verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Fatalln(err)
	}

	return rawIDToken
}

func RunCanary(serverDomain string) {
	c := http.DefaultClient

	// create
	request, err := http.NewRequest("POST",
		serverDomain+"/api/v1/tasks/",
		bytes.NewBufferString(`{"Title": "test", "Completed": false}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Id", getIDToken())
	if err != nil {
		panic(err)
	}
	response, err := c.Do(request)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if http.StatusOK != response.StatusCode {
		panic(err)
	}
	var task models.Task
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		panic(err)
	}
	if "test" != task.Title {
		panic(err)
	}

	// verify get, TODO this should be a lookup by ID
	request, err = http.NewRequest("GET",
		serverDomain+"/api/v1/tasks/",
		bytes.NewBufferString(``))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Id", getIDToken())
	if err != nil {
		panic(err)
	}
	response, err = c.Do(request)
	if err != nil {
		panic(err)
	}
	buf = new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if http.StatusOK != response.StatusCode {
		panic(err)
	}
	var tasks []models.Task
	err = json.Unmarshal(buf.Bytes(), &tasks)
	if err != nil {
		panic(err)
	}
	foundTask := false
	for _, returnedTask := range tasks {
		if returnedTask.ID == task.ID {
			foundTask = true
		}
	}
	if true != foundTask {
		panic(err)
	}

	// update
	// id := task.ID
	request, err = http.NewRequest("PUT",
		serverDomain+"/api/v1/tasks/?id="+fmt.Sprint(task.ID),
		bytes.NewBufferString(`{"title": "changedit", "completed": false}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Id", getIDToken())
	if err != nil {
		panic(err)
	}
	response, err = c.Do(request)
	if err != nil {
		panic(err)
	}
	buf = new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if http.StatusOK != response.StatusCode {
		panic(err)
	}
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		panic(err)
	}
	if "changedit" != task.Title {
		panic(err)
	}

	// delete
	request, err = http.NewRequest("DELETE",
		serverDomain+"/api/v1/tasks/?id="+fmt.Sprint(task.ID),
		bytes.NewBufferString(``))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Id", getIDToken())
	if err != nil {
		panic(err)
	}
	response, err = c.Do(request)
	if err != nil {
		panic(err)
	}
	buf = new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if http.StatusOK != response.StatusCode {
		panic(err)
	}

	// verify no get, TODO this should be a lookup by ID
	request, err = http.NewRequest("GET",
		serverDomain+"/api/v1/tasks/",
		bytes.NewBufferString(``))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Id", getIDToken())
	if err != nil {
		panic(err)
	}
	response, err = c.Do(request)
	if err != nil {
		panic(err)
	}
	buf = new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if http.StatusOK != response.StatusCode {
		panic(err)
	}
	err = json.Unmarshal(buf.Bytes(), &tasks)
	if err != nil {
		panic(err)
	}
	foundTask = false
	for _, returnedTask := range tasks {
		if returnedTask.ID == task.ID {
			foundTask = true
		}
	}
	if false != foundTask {
		panic(err)
	}
}
