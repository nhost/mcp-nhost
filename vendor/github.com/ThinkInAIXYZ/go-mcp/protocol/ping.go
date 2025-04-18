package protocol

type PingRequest struct{}

type PingResult struct{}

// NewPingRequest creates a new ping request
func NewPingRequest() *PingRequest {
	return &PingRequest{}
}

// NewPingResult creates a new ping response
func NewPingResult() *PingResult {
	return &PingResult{}
}
