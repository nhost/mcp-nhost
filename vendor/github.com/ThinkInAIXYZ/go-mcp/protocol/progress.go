package protocol

// ProgressNotification represents a progress notification for a long-running request
type ProgressNotification struct {
	ProgressToken ProgressToken `json:"progressToken"`
	Progress      float64       `json:"progress"`
	Total         float64       `json:"total,omitempty"`
}

// ProgressToken represents a token used to associate progress notifications with the original request
type ProgressToken interface{} // can be string or integer

// NewProgressNotification creates a new progress notification
func NewProgressNotification(token ProgressToken, progress float64, total float64) *ProgressNotification {
	return &ProgressNotification{
		ProgressToken: token,
		Progress:      progress,
		Total:         total,
	}
}
