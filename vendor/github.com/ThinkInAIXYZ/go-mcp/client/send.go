package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func (client *Client) sendMsgWithRequest(ctx context.Context, requestID protocol.RequestID, method protocol.Method, params protocol.ClientRequest) error {
	if requestID == nil {
		return fmt.Errorf("requestID can't is nil")
	}

	req := protocol.NewJSONRPCRequest(requestID, method, params)

	message, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if err := client.transport.Send(ctx, message); err != nil {
		return fmt.Errorf("sendRequest: transport send: %w", err)
	}
	return nil
}

func (client *Client) sendMsgWithResponse(ctx context.Context, requestID protocol.RequestID, result protocol.ClientResponse) error {
	if requestID == nil {
		return fmt.Errorf("requestID can't is nil")
	}

	resp := protocol.NewJSONRPCSuccessResponse(requestID, result)

	message, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	if err := client.transport.Send(ctx, message); err != nil {
		return fmt.Errorf("sendResponse: transport send: %w", err)
	}
	return nil
}

func (client *Client) sendMsgWithNotification(ctx context.Context, method protocol.Method, params protocol.ClientNotify) error {
	notify := protocol.NewJSONRPCNotification(method, params)

	message, err := json.Marshal(notify)
	if err != nil {
		return err
	}

	if err := client.transport.Send(ctx, message); err != nil {
		return fmt.Errorf("sendNotification: transport send: %w", err)
	}
	return nil
}

func (client *Client) sendMsgWithError(ctx context.Context, requestID protocol.RequestID, code int, msg string) error {
	if requestID == nil {
		return fmt.Errorf("requestID can't is nil")
	}

	resp := protocol.NewJSONRPCErrorResponse(requestID, code, msg)

	message, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	if err := client.transport.Send(ctx, message); err != nil {
		return fmt.Errorf("sendResponse: transport send: %w", err)
	}
	return nil
}
