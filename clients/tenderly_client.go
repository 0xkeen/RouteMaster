package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const TenderlyAPIEndpoint = "https://api.tenderly.co/api/v1"

type TenderlyClient struct {
	accessKey string
	user      string
	project   string
	client    *http.Client
}

type SimulateRequest struct {
	NetworkID   int    `json:"network_id"`
	From        string `json:"from"`
	To          string `json:"to"`
	Input       string `json:"input"`
	Save        bool   `json:"save"`
	SaveIfFails bool   `json:"save_if_fails"`
}

type SimulateResponse struct {
	Transaction struct {
		Hash string `json:"hash"`
	} `json:"transaction"`
	Simulation struct {
		Status bool   `json:"status"`
		Error  string `json:"error"`
	} `json:"simulation"`
}

func NewTenderlyClient(accessKey, user, project string) *TenderlyClient {
	return &TenderlyClient{
		accessKey: accessKey,
		user:      user,
		project:   project,
		client:    &http.Client{},
	}
}

func (c *TenderlyClient) Simulate(req *SimulateRequest) (*SimulateResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/account/%s/project/%s/simulate", TenderlyAPIEndpoint, c.user, c.project),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Access-Key", c.accessKey)

	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do request failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var result SimulateResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	if !result.Simulation.Status {
		return nil, fmt.Errorf("simulation failed: %s", result.Simulation.Error)
	}

	return &result, nil
}
