package checker

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/YOUR_GITHUB_USERNAME/go-api-monitor/internal/config"
)

type Result struct {
	Name           string        `json:"name"`
	Method         string        `json:"method"`
	URL            string        `json:"url"`
	ExpectedStatus int           `json:"expected_status"`
	Status         int           `json:"status"`
	LatencyMS      int64         `json:"latency_ms"`
	OK             bool          `json:"ok"`
	Error          string        `json:"error,omitempty"`
	CheckedAt      time.Time     `json:"checked_at"`
}

type Checker struct {
	client      *http.Client
	concurrency int
}

func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}

func NewChecker(client *http.Client, concurrency int) *Checker {
	if concurrency <= 0 {
		concurrency = 5
	}
	return &Checker{client: client, concurrency: concurrency}
}

func (c *Checker) CheckAll(ctx context.Context, endpoints []config.Endpoint) []Result {
	results := make([]Result, len(endpoints))
	sem := make(chan struct{}, c.concurrency)

	var wg sync.WaitGroup
	wg.Add(len(endpoints))

	for i := range endpoints {
		i := i
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			results[i] = c.checkOne(ctx, endpoints[i])
		}()
	}

	wg.Wait()
	return results
}

func (c *Checker) checkOne(ctx context.Context, ep config.Endpoint) Result {
	start := time.Now()
	r := Result{
		Name:           ep.Name,
		Method:         ep.Method,
		URL:            ep.URL,
		ExpectedStatus: ep.ExpectedStatus,
		CheckedAt:      time.Now().UTC(),
	}

	req, err := http.NewRequestWithContext(ctx, ep.Method, ep.URL, nil)
	if err != nil {
		r.Error = err.Error()
		return r
	}

	resp, err := c.client.Do(req)
	lat := time.Since(start)
	r.LatencyMS = lat.Milliseconds()

	if err != nil {
		// context cancellations and timeouts end up here too
		r.Error = err.Error()
		return r
	}
	defer resp.Body.Close()

	r.Status = resp.StatusCode
	r.OK = (resp.StatusCode == ep.ExpectedStatus)
	if !r.OK {
		r.Error = errors.New("unexpected status").Error()
	}
	return r
}
