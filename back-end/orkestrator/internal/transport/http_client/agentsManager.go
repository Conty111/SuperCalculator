package http_client

import (
	"bytes"
	"encoding/json"
	"errors"
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

type GetInfoResponse struct {
	ID   int32             `json:"id"`
	Info *models.AgentInfo `json:"info"`
}

type SetSettingsResponse struct {
	Status  string `json:"attr,status"`
	Message string `json:"message"`
}

type ErrResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func (s *AgentHTTPClient) GetAgentInfo(agent *models.AgentConfig) (*models.AgentInfo, error) {
	body, status, err := sendHTTPRequest(
		s.HTTPClient,
		nil,
		fmt.Sprintf("%s/status", agent.Address+strconv.Itoa(agent.HttpPort)),
		http.MethodGet,
	)
	if err != nil {
		return &models.AgentInfo{}, err
	}
	if status != http.StatusOK {
		decoded, err := unmarshalErrorResponse(body)
		if err != nil {
			log.Error().Err(err).Msg("failed to unmarshal err response")
			return nil, err
		}
		log.Error().
			Str("bodyStatus", decoded.Status).
			Str("error", decoded.Error).
			Msg("got not OK response status from GetInfo request")
		return nil, errors.New(decoded.Error)
	}
	var info GetInfoResponse
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal OK response")
		return nil, err
	}
	return info.Info, nil
}

func (s *AgentHTTPClient) SetAgentSettings(settings *models.Settings, agent *models.AgentConfig) error {
	reqBody, err := json.Marshal(settings.DurationSettings)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal duration settings")
		return err
	}
	body, status, err := sendHTTPRequest(
		s.HTTPClient,
		bytes.NewReader(reqBody),
		fmt.Sprintf("%s/calculator", agent.Address+strconv.Itoa(agent.HttpPort)),
		http.MethodPut,
	)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		decoded, err := unmarshalErrorResponse(body)
		if err != nil {
			log.Error().Err(err).Msg("failed to unmarshal err response")
			return err
		}
		log.Error().
			Str("bodyStatus", decoded.Status).
			Str("error", decoded.Error).
			Msg("got not OK response status from SetSettings request")
		return errors.New(decoded.Error)
	}
	var decoded SetSettingsResponse
	err = json.Unmarshal(body, &decoded)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal body data")
		return err
	}
	log.Info().
		Str("httpStatus", http.StatusText(status)).
		Str("message", decoded.Message).
		Msg("response from agent")

	return nil
}

func unmarshalErrorResponse(body []byte) (*ErrResponse, error) {
	var resp ErrResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// sendHTTPRequestToAgent sends request and returns body and status
func sendHTTPRequest(
	client *http.Client,
	reqBody io.Reader,
	addr string,
	method string) ([]byte, int, error) {

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
	return data, resp.StatusCode, nil
}
