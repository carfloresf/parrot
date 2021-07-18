package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/hellerox/parrot/model"
	"github.com/hellerox/parrot/service"
	"github.com/hellerox/parrot/storage"
)

const statusOK = "ok"

type response struct {
	Status  string `json:"status,omitempty"`
	OrderID int    `json:"orderID,omitempty"`
}

// Controller controller
type Controller struct {
	Router  *fasthttprouter.Router
	Storage storage.Storage
	Service service.Service
}

// InitializeRoutes route initialize
func (c *Controller) InitializeRoutes() {
	c.Router.HandleMethodNotAllowed = true
	c.Router.NotFound = c.notFound
	c.Router.MethodNotAllowed = c.methodNotAllowed
	c.Router.PanicHandler = c.panic

	c.Router.GET("/healthcheck", c.healthcheck)
	c.Router.POST("/user", c.createUser)
	c.Router.POST("/order", c.createOrder)
	c.Router.GET("/report", c.generateReport)
}

func (c *Controller) basicAuth(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
	}
}

func (c *Controller) createUser(ctx *fasthttp.RequestCtx) {
	var m model.User
	if err := json.Unmarshal(ctx.Request.Body(), &m); err != nil {
		respond(ctx, fasthttp.StatusBadRequest, fmt.Sprintf(`{"error":"%s"}`, err))
		return
	}

	err := c.Service.CreateUser(m)
	if err != nil {
		respond(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err))
		return
	}

	respondInterface(ctx,
		http.StatusCreated,
		response{
			Status: statusOK,
		})
}

func (c *Controller) createOrder(ctx *fasthttp.RequestCtx) {
	var m model.Order
	if err := json.Unmarshal(ctx.Request.Body(), &m); err != nil {
		respond(ctx, fasthttp.StatusBadRequest, fmt.Sprintf(`{"error":"%s"}`, err))
		return
	}

	oID, err := c.Service.CreateOrder(m)
	if err != nil {
		respond(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err))
		return
	}

	respondInterface(ctx,
		http.StatusCreated,
		response{
			Status:  statusOK,
			OrderID: oID,
		})
}

func (c *Controller) generateReport(ctx *fasthttp.RequestCtx) {

}
