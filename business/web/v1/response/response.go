package response

// Error is the structure used by the API to respond to the client
// when a failure happens.

type Error struct {
	SecurityToken  string `json:"securityToken,omitempty" mask:"filled"`
	Message        string `json:"message"`
	TraceID        string `json:"traceID"`
	AdditionalInfo any    `json:"data,omitempty"`
}
