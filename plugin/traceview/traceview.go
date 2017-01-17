package traceview

import (
	"golang.org/x/net/context"

	"github.com/tracelytics/go-traceview/v1/tv"
	gcontext "gopkg.in/h2non/gentleman.v1/context"
	"gopkg.in/h2non/gentleman.v1/plugin"
)

// New creates gentleman plugin for TraceView http client.
func New(ctx context.Context) plugin.Plugin {
	p := plugin.New()
	p.SetHandler("request", func(gctx *gcontext.Context, h gcontext.Handler) {
		l := tv.BeginHTTPClientLayer(ctx, gctx.Request)
		defer func() {
			l.AddHTTPResponse(gctx.Response, gctx.Error)
			l.End()
		}()
		h.Next(gctx)
	})
	return p
}
