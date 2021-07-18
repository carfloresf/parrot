package controller

import (
	"encoding/json"
	"fmt"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func respond(ctx *fasthttp.RequestCtx, code int, payload string) {
	ctx.SetStatusCode(code)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetBodyString(payload)
}

func respondInterface(ctx *fasthttp.RequestCtx, code int, payload interface{}) {
	ctx.SetStatusCode(code)
	ctx.SetContentType("application/json; charset=utf-8")

	response, err := json.Marshal(payload)
	if err != nil {
		log.Errorln(err.Error())
	}

	ctx.SetBodyString(string(response))
}

// methodNotAllowed handles the response when a method call is not allowed from
// the client.
func (c *Controller) methodNotAllowed(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetBodyString(fmt.Sprintf(`{"error": "Method not allowed calling %s"}`,
		ctx.Path()))
}

// notFound handles the response when a path have not been found.
func (c *Controller) notFound(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetBodyString(fmt.Sprintf(`{"error": "Path %s not found"}`,
		ctx.Path()))
}

// notFound handles the response when a path have not been found.
func (c *Controller) unauthorized(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusUnauthorized)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetBodyString(`{"er1ror": "Unauthorized"}`)
}

// healthcheck handler.
func (c *Controller) healthcheck(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetBodyString(`{"status": "ok"}`)
}

// panic handles when the server have a fatal error.
func (c *Controller) panic(ctx *fasthttp.RequestCtx, from interface{}) {
	log.Errorf(string(debug.Stack()))
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetBodyString(fmt.Sprintf(`{"error": "API fatal error calling %s"}`,
		ctx.Path()))
}
