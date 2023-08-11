package rest

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/valyala/fasthttp"
)

type Rest struct {
}

func NewRestServer() *Rest {
	return &Rest{}
}

func (o *Rest) Handler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	if handler, ok := handlers[path]; ok {
		handler(ctx)
		return
	}

	err := fmt.Sprintf("(%s) not found", path)
	sentry.CaptureMessage(err)
	ctx.Error(err, fasthttp.StatusNotFound)
}

var handlers = map[string]fasthttp.RequestHandler{
	"/healthz": allowMethods(healthz, fasthttp.MethodGet),
}

func healthz(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte("OK"))
}
