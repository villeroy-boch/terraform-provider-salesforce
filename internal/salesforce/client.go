package salesforce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	HostURL    string
	ApiVersion string
	HTTPClient *http.Client
	Auth       AuthStruct
}

type RespBody struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type AuthStruct struct {
	AuthHost     string `json:"authHost"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	GrantType    string `json:"grantType"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	BearerToken  string `json:"bearerToken"`
}

func NewClient(apiHost, apiVersion, authHost, clientID, clientSecret, grantType, username, password *string) (*Client, error) {

	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    *apiHost,
		ApiVersion: *apiVersion,
	}

	// If username or password not provided, return empty client
	if authHost == nil || clientID == nil || clientSecret == nil || grantType == nil || username == nil || password == nil {
		return &c, nil
	}

	c.Auth = AuthStruct{
		AuthHost:     *authHost,
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		GrantType:    *grantType,
		Username:     *username,
		Password:     *password,
	}
	var err error
	c.Auth.BearerToken, err = getBearerToken(c.Auth)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func getBearerToken(auth AuthStruct) (string, error) {

	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)

	fw, err := writer.CreateFormField("client_id")
	if err != nil {
		return "", err
	}
	_, err = io.Copy(fw, strings.NewReader(auth.ClientID))
	if err != nil {
		return "", err
	}
	fw, err = writer.CreateFormField("client_secret")
	if err != nil {
		return "", err
	}
	_, err = io.Copy(fw, strings.NewReader(auth.ClientSecret))
	if err != nil {
		return "", err
	}
	fw, err = writer.CreateFormField("grant_type")
	if err != nil {
		return "", err
	}
	_, err = io.Copy(fw, strings.NewReader(auth.GrantType))
	if err != nil {
		return "", err
	}
	fw, err = writer.CreateFormField("username")
	if err != nil {
		return "", err
	}
	_, err = io.Copy(fw, strings.NewReader(auth.Username))
	if err != nil {
		return "", err
	}
	fw, err = writer.CreateFormField("password")
	if err != nil {
		return "", err
	}
	_, err = io.Copy(fw, strings.NewReader(auth.Password))
	if err != nil {
		return "", err
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		auth.AuthHost,
		bytes.NewReader(reqBody.Bytes()),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	d := http.Client{Timeout: time.Duration(5) * time.Second}
	resp, err := d.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Fehler 3: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("status: %d, body: %s", resp.StatusCode, body)
	}

	respBody := &RespBody{}
	err = json.Unmarshal(body, respBody)
	if err != nil {
		fmt.Printf("Fehler 3: %s", err)
	}
	return respBody.AccessToken, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	// Note: this will have problems if there are redirects
	// see https://stackoverflow.com/a/31309385
	req.Header.Set("Authorization", "Bearer "+c.Auth.BearerToken)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err

}
