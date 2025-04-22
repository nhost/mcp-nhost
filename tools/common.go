package tools

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrInvalidRequestBody = errors.New("invalid request")

type GraphqlQueryRequest struct {
	Query     string
	Variables map[string]any
}

func QueryRequestFromParams(params map[string]any) (GraphqlQueryRequest, error) {
	query, ok := params["query"].(string)
	if !ok {
		return GraphqlQueryRequest{}, fmt.Errorf("%w: query is required", ErrInvalidRequestBody)
	}

	var variables map[string]any
	if v, ok := params["variables"]; ok {
		switch v := v.(type) {
		case string:
			if err := json.Unmarshal([]byte(v), &variables); err != nil {
				return GraphqlQueryRequest{}, fmt.Errorf("failed to unmarshal variables: %w", err)
			}
		case map[string]any:
			variables = v
		default:
			return GraphqlQueryRequest{},
				fmt.Errorf("%w: variables must be a string or map[string]any", ErrInvalidRequestBody)
		}
	}

	return GraphqlQueryRequest{
		Query:     query,
		Variables: variables,
	}, nil
}

type GraphqlQueryWithRoleRequest struct {
	Query     string
	Variables map[string]any
	Role      string
}

func RoleFromParams(params map[string]any) (string, error) {
	var role string
	r, ok := params["role"]
	if ok {
		switch r := r.(type) {
		case string:
			role = r
		default:
			return "", fmt.Errorf("%w: role must be a string", ErrInvalidRequestBody)
		}
	}

	if role == "" {
		return "", fmt.Errorf("%w: role is required", ErrInvalidRequestBody)
	}

	return role, nil
}

func QueryRequestWithRoleFromParams(params map[string]any) (GraphqlQueryWithRoleRequest, error) {
	request, err := QueryRequestFromParams(params)
	if err != nil {
		return GraphqlQueryWithRoleRequest{}, err
	}

	role, err := RoleFromParams(params)
	if err != nil {
		return GraphqlQueryWithRoleRequest{}, err
	}

	return GraphqlQueryWithRoleRequest{
		Query:     request.Query,
		Variables: request.Variables,
		Role:      role,
	}, nil
}

func ProjectFromParams(params map[string]any) (string, error) {
	var project string
	p, ok := params["projectSubdomain"]
	if ok {
		switch r := p.(type) {
		case string:
			project = r
		default:
			return "", fmt.Errorf("%w: project must be a string", ErrInvalidRequestBody)
		}
	}

	if project == "" {
		return "", fmt.Errorf("%w: project is required", ErrInvalidRequestBody)
	}

	return project, nil
}

func FromParams[T any](params map[string]any, name string) (T, error) { //nolint:ireturn
	var res T
	p, ok := params[name]
	if ok {
		switch r := p.(type) {
		case T:
			res = r
		default:
			return res, fmt.Errorf("%w: %s must be a %T", ErrInvalidRequestBody, name, res)
		}
	}

	return res, nil
}
