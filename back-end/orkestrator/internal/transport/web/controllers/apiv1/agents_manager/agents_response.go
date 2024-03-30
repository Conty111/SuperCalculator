package agents_manager

type WorkerResponse struct {
	Status   string                 `json:"worker_response_status"`
	Response map[string]interface{} `json:"worker_response"`
}

type WorkersListResponse struct {
	Status    string           `json:"status"`
	Responses []WorkerResponse `json:"responses"`
}
