package authsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

func (c *AuthClient) HttpNewRequest(path, token string) (string, error) {
	// 定义请求的 URL
	//path := "/getList"

	// 创建一个新的 GET 请求
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// 设置请求头
	req.Header.Set("Authorization", token)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}
	// 打印响应
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
	return string(body), nil
}
