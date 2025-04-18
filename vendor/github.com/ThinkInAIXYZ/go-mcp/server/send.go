package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func (server *Server) sendMsgWithRequest(ctx context.Context, sessionID string, requestID protocol.RequestID,
	method protocol.Method, params protocol.ServerRequest,
) error { //nolint:whitespace
	if requestID == nil {
		return fmt.Errorf("requestID can't is nil")
	}

	req := protocol.NewJSONRPCRequest(requestID, method, params)

	message, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if err := server.transport.Send(ctx, sessionID, message); err != nil {
		return fmt.Errorf("sendRequest: transport send: %w", err)
	}
	return nil
}

func (server *Server) sendMsgWithResponse(ctx context.Context, sessionID string, requestID protocol.RequestID, result protocol.ServerResponse) error {
	resp := protocol.NewJSONRPCSuccessResponse(requestID, result)

	message, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	if err := server.transport.Send(ctx, sessionID, message); err != nil {
		return fmt.Errorf("sendResponse: transport send: %w", err)
	}
	return nil
}

func (server *Server) sendMsgWithNotification(ctx context.Context, sessionID string, method protocol.Method, params protocol.ServerNotify) error {
	notify := protocol.NewJSONRPCNotification(method, params)

	message, err := json.Marshal(notify)
	if err != nil {
		return err
	}

	if err := server.transport.Send(ctx, sessionID, message); err != nil {
		return fmt.Errorf("sendNotification: transport send: %w", err)
	}
	return nil
}

func (server *Server) sendMsgWithError(ctx context.Context, sessionID string, requestID protocol.RequestID, code int, msg string) error {
	if requestID == nil {
		return fmt.Errorf("requestID can't is nil")
	}

	resp := protocol.NewJSONRPCErrorResponse(requestID, code, msg)

	message, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	if err := server.transport.Send(ctx, sessionID, message); err != nil {
		return fmt.Errorf("sendResponse: transport send: %w", err)
	}
	return nil
}
