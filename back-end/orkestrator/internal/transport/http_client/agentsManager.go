package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
	"time"
)

type AgentHTTPClient struct {
	HTTPClient *http.Client
}

func NewAgentHTTPClient(timeout time.Duration) *AgentHTTPClient {
	return &AgentHTTPClient{
		HTTPClient: &http.Client{Timeout: timeout},
	}
}

func (s *AgentHTTPClient) GetAgentInfo(agent models.AgentConfig) (map[string]interface{}, int, error) {

	body, status, err := sendHTTPRequest(
		s.HTTPClient,
		nil,
		fmt.Sprintf("%s/status", agent.Address+strconv.Itoa(agent.HttpPort)),
		http.MethodGet,
	)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return body, status, nil
}

func (s *AgentHTTPClient) SetAgentSettings(settings models.Settings, agent models.AgentConfig) (map[string]interface{}, int, error) {
	reqBody, err := json.Marshal(settings.DurationSettings)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal duration settings")
		return nil, http.StatusInternalServerError, err
	}
	body, status, err := sendHTTPRequest(
		s.HTTPClient,
		bytes.NewReader(reqBody),
		fmt.Sprintf("%s/calculator", agent.Address+strconv.Itoa(agent.HttpPort)),
		http.MethodPut,
	)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return body, status, nil
}

// sendHTTPRequestToAgent sends request and returns body and status
func sendHTTPRequest(
	client *http.Client,
	reqBody io.Reader,
	addr string,
	method string) (map[string]interface{}, int, error) {

	body := make(map[string]interface{})

	req, err := http.NewRequest(method, fmt.Sprintf("http://%s", addr), reqBody)
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return nil, 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to send request")
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
		return nil, 0, err
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal body data")
		return nil, 0, err
	}
	return body, resp.StatusCode, nil
}
