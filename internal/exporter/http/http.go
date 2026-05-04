// Package exporter
package exporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ragnacron/msma/internal/model"
)

type HTTPExporter struct {
	Endpoint string
	Client   *http.Client
}

func New(endpoint string) *HTTPExporter {
	return &HTTPExporter{
		Endpoint: endpoint,
		Client:   &http.Client{},
	}
}

func (h *HTTPExporter) Export(metrics []model.Metric) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(metrics); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", h.Endpoint, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := h.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("error status code: %d", res.StatusCode)
	}

	return nil
}
