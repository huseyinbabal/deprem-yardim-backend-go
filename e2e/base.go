package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"io"
	"net/http"
	"strings"
	"time"
)

// TestSuite base test suite
type TestSuite struct {
	compose *testcontainers.LocalDockerCompose
	Client  *Client
}

// SetupBaseSuite base setup site for tests
func (t *TestSuite) SetupBaseSuite() error {
	composeFilePaths := []string{"../resources/docker-compose.yml"}
	identifier := strings.ToLower(uuid.New().String())
	t.compose = testcontainers.NewLocalDockerCompose(composeFilePaths, identifier)
	execError := t.compose.WithCommand([]string{"up", "-d"}).Invoke()
	if execError.Error != nil {
		return execError.Error
	}
	time.Sleep(5 * time.Second) // No health check for scratch images.
	t.Client = NewClient("http://localhost:80")
	return nil
}

// TearDownBaseSuite tear down method for tests
func (t *TestSuite) TearDownBaseSuite() error {
	execError := t.compose.
		WithCommand([]string{"down"}).
		Invoke()
	err := execError.Error
	if err != nil {
		return fmt.Errorf("could not shutdown compose stack: %v", err)
	}
	return nil
}

func NewClient(baseURL string) *Client {
	return &Client{
		client:  &http.Client{},
		BaseURL: baseURL,
	}
}

type Client struct {
	client  *http.Client
	BaseURL string
}

func (c *Client) MakeRequest(r *http.Request, v interface{}) error {
	resp, err := c.client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	switch status := resp.StatusCode; {
	case status >= 200 && status < 300:
		if err := json.Unmarshal(data, v); err != nil {
			return fmt.Errorf("error unmarshalling response: %w", err)
		}
		return nil
	case status >= 400 && status < 500: // 4xx
		var apiError error
		if err := json.Unmarshal(data, &apiError); err != nil {
			return fmt.Errorf("error unmarshalling response: %w", err)
		}
		return apiError
	default:
		return fmt.Errorf("error response: %s", string(data))
	}
}
