package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/enums"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/helpers"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type AgentManager struct {
	ApiType         enums.ApiType
	Agents          []models.AgentConfig
	TimeoutResponse time.Duration
	HTTPClient      http.Client
}

func (s *AgentManager) SetSettings(settings *models.Settings) []*helpers.AgentResponse {
	//TODO implement me
	panic("implement me")
}

func (s *AgentManager) GetWorkersInfo() []*helpers.AgentResponse {
	wg := sync.WaitGroup{}
	wg.Add(len(s.Agents))

	responses := make([]*helpers.AgentResponse, len(s.Agents))
	for i, agent := range s.Agents {
		agent := agent
		i := i
		go func() {
			defer wg.Done()
			res, status, err := s.GetAgentInfo(agent)
			if err != nil {
				log.Error().Err(err).Msg("failed to send request to the agent")
				responses[i] = &helpers.AgentResponse{
					Body:   nil,
					Status: http.StatusInternalServerError,
				}
			}
			responses[i] = &helpers.AgentResponse{
				Body:   res,
				Status: status,
			}
		}()
	}
	wg.Wait()
	return responses
	//return responses
	//panic("implement me")
}

func NewAgentManager(
	apiType enums.ApiType,
	agents []models.AgentConfig,
	timeout time.Duration) *AgentManager {

	return &AgentManager{
		ApiType:         apiType,
		Agents:          agents,
		TimeoutResponse: timeout,
	}
}

func (s *AgentManager) GetAgentInfo(
	agent models.AgentConfig,
) (map[string]interface{}, int, error) {

	switch s.ApiType {
	case enums.GrpcApi:
		return nil, 0, nil
	default:
		client := http.Client{
			Timeout: s.TimeoutResponse,
		}
		body, status, err := sendHTTPRequest(
			&client,
			nil,
			fmt.Sprintf("%s/api/v1/status", agent.Address+":"+strconv.Itoa(agent.HttpPort)),
			http.MethodGet,
		)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return body, status, nil
	}
}

func (s *AgentManager) SetAgentsSettings(
	settings models.Settings,
	agent models.AgentConfig) (map[string]interface{}, int, error) {

	switch s.ApiType {
	case enums.GrpcApi:
		return nil, 0, nil
	default:
		client := http.Client{
			Timeout: s.TimeoutResponse,
		}
		reqBody, err := json.Marshal(settings.DurationSettings)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal duration settings")
			return nil, http.StatusInternalServerError, err
		}
		body, status, err := sendHTTPRequest(
			&client,
			bytes.NewReader(reqBody),
			fmt.Sprintf("%s/calculator", agent.Address+":"+strconv.Itoa(agent.HttpPort)),
			http.MethodPut,
		)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return body, status, nil
	}
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
