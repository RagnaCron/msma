// Package http
package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ragnacron/msma/internal/model"
)

type HTTPExporter struct {
	Endpoint string
	Client   *http.Client
}

func NewHTTP(endpoint string, timeout int64) *HTTPExporter {
	return &HTTPExporter{
		Endpoint: endpoint,
		Client: &http.Client{
			Timeout: time.Duration(timeout * int64(time.Second)),
		},
	}
}

func (h *HTTPExporter) Export(metrics []model.Metric) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	r := bytes.NewReader(data)
	req, err := http.NewRequest("POST", h.Endpoint, r)
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
