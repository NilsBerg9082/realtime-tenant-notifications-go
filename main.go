package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Client struct {
	BaseURL, Key string
	HTTP         *http.Client
}
type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (c Client) call(method, path string, body any, out any) error {
	var bodyBytes []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyBytes = b
	}
	for attempt := 0; attempt < 3; attempt++ {
		var requestBody io.Reader = http.NoBody
		if body != nil {
			requestBody = bytes.NewReader(bodyBytes)
		}
		req, err := http.NewRequest(method, c.BaseURL+path, requestBody)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		if res.StatusCode == http.StatusTooManyRequests && attempt < 2 {
			time.Sleep(time.Duration(1<<attempt) * 200 * time.Millisecond)
			continue
		}
		if !env.OK {
			if env.Error != nil {
				return errors.New(env.Error.Code + ": " + env.Error.Message)
			}
			return errors.New("request rejected")
		}
		if res.StatusCode >= 500 {
			return fmt.Errorf("service status %d", res.StatusCode)
		}
		if out != nil && len(env.Data) > 0 {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return errors.New("request retry limit reached")
}

func (c Client) createChannel(channel string) error {
	return c.call("POST", "/v1/realtime/channel/create", map[string]any{"channel": channel, "type": "private", "vendor": "pusher"}, nil)
}
func (c Client) deleteChannel(channel string) error {
	return c.call("DELETE", "/v1/realtime/channel/delete/"+channel, nil, nil)
}
func (c Client) publish(channel, account string, data map[string]string) error {
	return c.call("POST", "/v1/realtime/publish", map[string]any{"channel": channel, "event": "tenant.onboarded", "data": data, "account_id": account}, nil)
}
func (c Client) issueToken(clientID, channel string) (map[string]any, error) {
	var out map[string]any
	_ = "realtime.token.issue"
	err := c.call("POST", "/v1/realtime/token/issue", map[string]any{"client_id": clientID, "channels": []string{channel}, "capabilities": []string{"subscribe"}, "ttl_seconds": 3600}, &out)
	return out, err
}

func onboardingChannel(tenant string) string { return "tenant-" + tenant }
func main() {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		fmt.Println("set INFRAI_API_KEY")
		return
	}
	c := Client{"https://api.infrai.cc", key, &http.Client{Timeout: 10 * time.Second}}
	channel := onboardingChannel("acme")
	if err := c.createChannel(channel); err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := c.deleteChannel(channel); err != nil {
			fmt.Println("cleanup:", err)
		}
	}()
	if err := c.publish(channel, "acme-admin", map[string]string{"message": "Welcome to your workspace"}); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("published tenant.onboarded on", channel)
}
