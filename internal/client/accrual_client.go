package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

type AccrualClient struct {
	baseURL    string
	httpClient *http.Client
	retryAfter time.Duration
}

func NewAccrualClient(baseURL string) *AccrualClient {
	return &AccrualClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		retryAfter: 60 * time.Second,
	}
}

func (c *AccrualClient) FetchAccrual(ctx context.Context, orderNumber string) (*AccrualResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, orderNumber)

	retryDelay := time.Second
	const maxRetries = 5

	for attempt := 0; attempt < maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("request creation failed: %w", err)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			time.Sleep(retryDelay)
			retryDelay *= 2
			continue
		}
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.Printf("failed to close response body: %v", err)
				return
			}
		}()
		switch resp.StatusCode {
		case http.StatusOK:
			var result AccrualResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return nil, fmt.Errorf("json decode failed: %w", err)
			}
			return &result, nil

		case http.StatusNoContent:
			return nil, utils.ErrOrderNotFoundInAccrual

		case http.StatusTooManyRequests:
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if sec, err := strconv.Atoi(ra); err == nil {
					c.retryAfter = time.Duration(sec) * time.Second
				}
			}
			time.Sleep(c.retryAfter)
			continue

		case http.StatusInternalServerError:
			time.Sleep(retryDelay)
			retryDelay *= 2
			continue

		default:
			return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}
	}
	return nil, fmt.Errorf("max retries exceeded")
}
