package protocol

// PaginatedRequest represents a request that supports pagination
type PaginatedRequest struct {
	Cursor string `json:"cursor,omitempty"`
}

// PaginatedResult represents a response that supports pagination
type PaginatedResult struct {
	NextCursor string `json:"nextCursor,omitempty"`
}
