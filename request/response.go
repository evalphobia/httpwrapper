package request

import (
	"fmt"

	"gopkg.in/h2non/gentleman.v1"
)

// Response is wrapper struct of *gentleman.Response
type Response struct {
	*gentleman.Response
}

// HasStatusCodeError returns error if HTTP status code is not 2xx code.
func (r Response) HasStatusCodeError() error {
	if r.Response.Ok {
		return nil
	}
	return fmt.Errorf("%d %s", r.StatusCode, r.String())
}
