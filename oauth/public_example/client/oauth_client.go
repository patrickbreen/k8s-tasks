package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

func getBearerToken(client *http.Client) []byte {
	request, err := http.NewRequest("POST",
		"https://keycloak.dev.leetcyber.com/auth/realms/basic/protocol/openid-connect/token",
		bytes.NewBufferString(`client_id=client-secret&username=patrick&password=star&grant_type=password`))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		panic(err)
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if http.StatusOK != response.StatusCode {
		panic(err)
	}
	var oauthToken oauth2.Token
	json.Unmarshal(buf.Bytes(), &oauthToken)
	fmt.Printf("token: %v\n", oauthToken)
	return buf.Bytes()
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

	tokenData := getBearerToken(client)

	//TODO if token about to expire renew it, if renewal token about to expire, get new token

	// send request to server
	req, err := http.NewRequest("GET", "http://oauth.dev.leetcyber.com:8090/", bytes.NewBuffer(nil))
	req.Header.Set("Authorization", string(tokenData))
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
