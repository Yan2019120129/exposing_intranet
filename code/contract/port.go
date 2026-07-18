package contract

// PortRequest is the wire request used to manage a client's port mappings.
type PortRequest struct {
	Symbol     string `json:"symbol" binding:"required"`
	Action     string `json:"action" binding:"required"`
	ServerPort string `json:"server_port,omitempty"`
	LocalAddr  string `json:"local_addr,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

// PortResponse is the wire response returned by the port management API.
type PortResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    []PortMappingInfo `json:"data,omitempty"`
}

// PortMappingInfo describes one persisted or active port mapping.
type PortMappingInfo struct {
	ServerPort string `json:"server_port"`
	LocalAddr  string `json:"local_addr"`
	Comment    string `json:"comment"`
	Status     string `json:"status"`
}
