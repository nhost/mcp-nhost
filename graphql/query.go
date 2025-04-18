package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrQueryingGraphqlEndpoint = errors.New("error querying graphql endpoint")
	ErrGraphqlContainErrors    = errors.New("graphql response contains errors")
)

func Query[T any](
	ctx context.Context,
	graphqlURL string,
	query string,
	variables map[string]any,
	response *Response[T],
	requestInterceptor ...func(ctx context.Context, req *http.Request) error,
) error {
	requestBody, err := json.Marshal(map[string]any{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		graphqlURL,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	for _, interceptor := range requestInterceptor {
		if err := interceptor(ctx, request); err != nil {
			return fmt.Errorf("failed to intercept request: %w", err)
		}
	}

	client := &http.Client{} //nolint:exhaustruct
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%w: %s\n%s", ErrQueryingGraphqlEndpoint, resp.Status, b)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, response); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if len(response.Errors) > 0 {
		return fmt.Errorf("%w: %s", ErrGraphqlContainErrors, body)
	}

	return nil
}
