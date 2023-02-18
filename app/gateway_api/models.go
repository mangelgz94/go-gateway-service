package gateway_api

import "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/models"

type findNumberPositionResponse struct {
	Number int    `json:"number,omitempty"`
	Error  string `json:"error,omitempty"`
}

type getUsersResponse struct {
	Users []*models.User `json:"users,omitempty"`
	Error string         `json:"error,omitempty"`
}
