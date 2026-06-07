package dto

// TestPayload is the request body for creating and updating Test records.
type TestPayload struct {
	Name string `json:"name"`
}
