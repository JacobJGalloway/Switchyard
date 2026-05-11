package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateBOLRequest is the payload sent to POST /api/BillOfLading on the
// Switchyard .NET Logistics API. CustomerFirstName/LastName/City/State come
// from store master data — the service layer must look up the destination
// store before calling CreateBOL.
type CreateBOLRequest struct {
	CustomerFirstName string             `json:"customerFirstName"`
	CustomerLastName  string             `json:"customerLastName"`
	City              string             `json:"city"`
	State             string             `json:"state"`
	LineEntries       []LineEntryRequest `json:"lineEntries"`
}

// LineEntryRequest mirrors the .NET LineEntry shape.
// Quantity: positive = incoming (warehouse pickup), negative = outgoing (store delivery).
type LineEntryRequest struct {
	LocationID string `json:"locationId"`
	SKUMarker  string `json:"sKUMarker"`
	Quantity   int    `json:"quantity"`
}

// ReplaceStopRequest mirrors the .NET ReplaceStopRequest shape.
type ReplaceStopRequest struct {
	OldLocationID string `json:"oldLocationId"`
	NewLocationID string `json:"newLocationId"`
}

// LogisticsClient is the only surface in the Go backend that writes to the
// Switchyard .NET Logistics API. No other package may call that API.
type LogisticsClient interface {
	// CreateBOL submits a validated plan to the .NET Logistics API.
	// Returns the transactionId assigned by the .NET system.
	CreateBOL(ctx context.Context, req *CreateBOLRequest) (string, error)

	// ProcessStop marks all line entries for a location on a committed BOL
	// as processed. Called when a driver logs stop completion.
	ProcessStop(ctx context.Context, transactionID, locationID string) error

	// ReplaceStop atomically moves unprocessed line entries from one location
	// to another on a committed BOL. Emergency dispatcher override only —
	// returns an error if the old location is already processed (409 Conflict).
	ReplaceStop(ctx context.Context, transactionID string, req *ReplaceStopRequest) error
}

type HTTPLogisticsClient struct {
	baseURL       string
	tokenProvider TokenProvider
	httpClient    *http.Client
}

func NewLogisticsClient(baseURL string, tp TokenProvider) *HTTPLogisticsClient {
	return &HTTPLogisticsClient{
		baseURL:       baseURL,
		tokenProvider: tp,
		httpClient:    newHTTPClient(),
	}
}

func (c *HTTPLogisticsClient) CreateBOL(ctx context.Context, req *CreateBOLRequest) (string, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshalling CreateBOL request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/api/BillOfLading", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("CreateBOL returned %d", resp.StatusCode)
	}

	var transactionID string
	if err := json.NewDecoder(resp.Body).Decode(&transactionID); err != nil {
		return "", fmt.Errorf("decoding CreateBOL response: %w", err)
	}

	return transactionID, nil
}

func (c *HTTPLogisticsClient) ProcessStop(ctx context.Context, transactionID, locationID string) error {
	url := fmt.Sprintf("%s/api/BillOfLading/%s/process/%s", c.baseURL, transactionID, locationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("ProcessStop returned %d for transaction %s location %s",
			resp.StatusCode, transactionID, locationID)
	}

	return nil
}

func (c *HTTPLogisticsClient) ReplaceStop(ctx context.Context, transactionID string, req *ReplaceStopRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshalling ReplaceStop request: %w", err)
	}

	url := fmt.Sprintf("%s/api/BillOfLading/%s/replace-stop", c.baseURL, transactionID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusConflict:
		return fmt.Errorf("ReplaceStop conflict: location %s is already processed on transaction %s",
			req.OldLocationID, transactionID)
	default:
		return fmt.Errorf("ReplaceStop returned %d for transaction %s", resp.StatusCode, transactionID)
	}
}

func (c *HTTPLogisticsClient) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.tokenProvider())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}
