package request

import "gopkg.in/h2non/gentleman.v1"

// Response is wrapper struct of *gentleman.Response
type Response struct {
	*gentleman.Response
}
