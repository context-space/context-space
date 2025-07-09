package supabase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/context-space/context-space/backend/internal/shared/utils"
)

type Client struct {
	client  http.Client
	baseURL string
	apiKey  string
	token   string
}

func NewClient(projectReference string, serviceRole string) *Client {
	baseURL := fmt.Sprintf("https://%s.supabase.co/auth/v1", projectReference)
	return &Client{
		client: http.Client{
			Timeout: time.Second * 10,
		},
		baseURL: baseURL,
		apiKey:  serviceRole,
		token:   serviceRole,
	}
}

func (c *Client) newRequest(path string, method string, body io.Reader) (*http.Request, error) {
	url := utils.StringsBuilder(c.baseURL, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("apiKey", c.apiKey)
	if c.token != "" {
		req.Header.Add("Authorization", utils.StringsBuilder("Bearer ", c.token))
	}
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

const adminUsersPath = "/admin/users"

// AdminGetUser
//
// GET /admin/users/{user_id}
//
// Get a user by their user_id.
func (c *Client) AdminGetUser(req AdminGetUserRequest) (*AdminGetUserResponse, error) {
	path := fmt.Sprintf("%s/%s", adminUsersPath, req.UserID)
	r, err := c.newRequest(path, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("response status code %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("response status code %d: %s", resp.StatusCode, fullBody)
	}

	var res AdminGetUserResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
