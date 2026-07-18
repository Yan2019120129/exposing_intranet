package apis

import (
	"context"

	"my-base/code/contract"
)

// PortManager is the client-side outbound port management boundary.
type PortManager interface {
	ManagePort(context.Context, contract.PortRequest) (contract.PortResponse, error)
}

func (a *HTTPAPI) ManagePort(ctx context.Context, req contract.PortRequest) (contract.PortResponse, error) {
	var result contract.PortResponse
	if err := a.post(ctx, "/api/client/port", req, &result); err != nil {
		return result, err
	}
	if !result.Success {
		return result, apiError(result.Message)
	}
	return result, nil
}

type apiError string

func (e apiError) Error() string { return string(e) }
