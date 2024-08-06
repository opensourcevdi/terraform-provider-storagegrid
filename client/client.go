package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	ApiUrl    string
	AccountId string
	Username  string
	Password  string
	Token     string
}

func (c *Client) authorize() (token string, err error) {
	return token, c.Api("POST", "/authorize", map[string]any{
		"accountId": c.AccountId,
		"username":  c.Username,
		"password":  c.Password,
	}, &token)
}

type Response struct {
	Status  string
	Data    any
	Message struct {
		Text string
		Key  string
	}
}

func (c *Client) Api(method, endpoint string, data any, res any) (err error) {
	if c.Token == "" && endpoint != "/authorize" {
		token, err := c.authorize()
		if err != nil {
			return err
		}
		c.Token = token
	}
	var body io.Reader
	if data != nil {
		dataJson, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(dataJson)
	}
	url, err := url.JoinPath(c.ApiUrl, "/api/v3/", endpoint)
	if err != nil {
		return
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	if c.Token != "" {
		req.Header.Add("authorization", "Bearer "+c.Token)
	}
	respHttp, err := (&http.Client{}).Do(req)
	if err != nil {
		return
	}
	respBody, err := io.ReadAll(respHttp.Body)
	if err != nil {
		return
	}
	if (respHttp.StatusCode == 200 || respHttp.StatusCode == 204) && len(respBody) == 0 {
		return nil
	}
	resp := Response{
		Data: &res,
	}
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return
	}
	if resp.Status != "success" {
		return fmt.Errorf("%v: %v", resp.Message.Key, resp.Message.Text)
	}
	return
}

type User struct {
	Id         string
	AccountId  string
	UniqueName string
	UserURN    string
}

func (c *Client) CreateUser(username string) (res User, err error) {
	return res, c.Api("POST", "/org/users", map[string]any{
		"uniqueName": "user/" + username,
		"fullName":   username,
		"disable":    true,
	}, &res)
}

func (c *Client) GetUsers() (res []User, err error) {
	return res, c.Api("GET", "/org/users", nil, &res)
}

func (c *Client) GetUser(userId string) (res User, err error) {
	return res, c.Api("GET", "/org/users/"+url.PathEscape(userId), nil, &res)
}

func (c *Client) GetUserByName(username string) (res User, err error) {
	return res, c.Api("GET", "/org/users/user/"+url.PathEscape(username), nil, &res)
}

func (c *Client) DeleteUser(userId string) (err error) {
	return c.Api("DELETE", "/org/users/"+url.PathEscape(userId), nil, nil)
}

type AccessKey struct {
	Id              string
	DisplayName     string
	UserURN         string
	UserUUID        string
	Expires         string
	AccessKey       string
	SecretAccessKey string
}

func (c *Client) CreateAccessKey(userId string) (res AccessKey, err error) {
	return res, c.Api("POST", "/org/users/"+url.PathEscape(userId)+"/s3-access-keys", map[string]any{}, &res)
}

func (c *Client) GetAccessKeys(userId string) (res []AccessKey, err error) {
	return res, c.Api("GET", "/org/users/"+url.PathEscape(userId)+"/s3-access-keys", nil, &res)
}

func (c *Client) GetAccessKey(userId string, accessKey string) (res AccessKey, err error) {
	return res, c.Api("GET", "/org/users/"+url.PathEscape(userId)+"/s3-access-keys/"+url.PathEscape(accessKey), nil, &res)
}

func (c *Client) DeleteAccessKey(userId string, id string) (err error) {
	return c.Api("DELETE", "/org/users/"+url.PathEscape(userId)+"/s3-access-keys/"+url.PathEscape(id), nil, nil)
}
