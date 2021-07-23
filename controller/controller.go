package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/buaazp/fasthttprouter"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"

	"github.com/hellerox/parrot/model"
	"github.com/hellerox/parrot/service"
)

const statusOK = "ok"

type response struct {
	Status  string `json:"status,omitempty"`
	OrderID int    `json:"orderID,omitempty"`
}

// Controller controller
type Controller struct {
	Router  *fasthttprouter.Router
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
	c.Router.POST("/order", c.BasicAuth(c.createOrder))
	c.Router.POST("/report", c.BasicAuth(c.generateReport))
}

func (c *Controller) createUser(ctx *fasthttp.RequestCtx) {
	var u model.User
	if err := json.Unmarshal(ctx.Request.Body(), &u); err != nil {
		respondError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}

	err := c.Service.CreateUser(u)
	if err != nil {
		respondError(ctx, fasthttp.StatusInternalServerError, err.Error())
		return
	}

	respondInterface(ctx,
		http.StatusCreated,
		response{
			Status: statusOK,
		},
	)
}

func (c *Controller) createOrder(ctx *fasthttp.RequestCtx) {
	var m model.Order
	if err := json.Unmarshal(ctx.Request.Body(), &m); err != nil {
		respondError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}

	oID, err := c.Service.CreateOrder(m)
	if err != nil {
		respondError(ctx, fasthttp.StatusInternalServerError, err.Error())
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
	var req model.GenerateReportRequest
	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		respondError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}

	response, err := c.Service.GenerateReport(req)
	if err != nil {
		respondError(ctx, fasthttp.StatusInternalServerError, err.Error())
		return
	}

	respondInterface(ctx,
		http.StatusCreated,
		response)
}

const authPrefix = "Basic "

func (c *Controller) BasicAuth(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	fn := func(ctx *fasthttp.RequestCtx) {
		auth := ctx.Request.Header.Peek("Authorization")

		if bytes.HasPrefix(auth, []byte(authPrefix)) {
			payload, err := base64.StdEncoding.DecodeString(string(auth[len([]byte(authPrefix)):]))
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				user := pair[0]
				hash := c.Service.GetUserHash(string(user))
				log.Println(string(pair[1]))
				u := model.User{
					Email:        string(user),
					Password:     string(pair[1]),
					PasswordHash: hash,
				}

				errHash := compareHash(&u)
				if errHash != nil {
					log.Errorln("hash error: ", errHash)
				} else {
					next(ctx)

					return
				}
			}
		}

		c.unauthorized(ctx)
	}

	return fn
}

func compareHash(user *model.User) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(user.Password+service.Pepper))
	if err != nil {
		return err
	}

	return nil
}
