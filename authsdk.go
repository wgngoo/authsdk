package authsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type AuthClient struct {
	BaseURL string
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	Error string `json:"error,omitempty"`
}

func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{BaseURL: baseURL}
}

func (c *AuthClient) Authenticate(username, password string) (string, error) {
	authReq := AuthRequest{Username: username, Password: password}
	reqBody, err := json.Marshal(authReq)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(c.BaseURL+"/auth", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return "", err
	}

	if authResp.Error != "" {
		return "", errors.New(authResp.Error)
	}

	return authResp.Token, nil
}
