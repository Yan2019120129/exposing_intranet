package contract

// AuthRequest is the wire request used by a client to authenticate with the
// server. It is shared by the inbound server handler and the outbound client
// API adapter.
type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hostname string `json:"hostname"`
}

// AuthResponse is the wire response returned by the server after
// authentication.
type AuthResponse struct {
	Success bool   `json:"success"`
	Symbol  string `json:"symbol"`
	Message string `json:"message"`
}
