package gateway_api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/app/gateway_api/middlewares/basic_auth"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type httpMiddleware interface {
	HandleAuthorization(next http.Handler) http.Handler
}

type Handler struct {
	Route      func(r *mux.Route)
	Func       http.HandlerFunc
	Middleware httpMiddleware
}

func (h *Handler) AddRoute(r *mux.Router) {
	h.Route(r.NewRoute().Handler(
		h.Middleware.HandleAuthorization(h.Func),
	))
}

func (g *GatewayServer) GetUsersHandler() *Handler {
	return &Handler{
		Route: func(r *mux.Route) {
			r.Path("/users").Methods(http.MethodGet)
		},
		Func:       g.getUsers,
		Middleware: basic_auth.New(&basic_auth.BasicAuthConfig{AuthUser: g.config.AuthUser, AuthPassword: g.config.AuthPassword}),
	}
}

func (g *GatewayServer) FindNumberPositionHandler() *Handler {
	return &Handler{
		Route: func(r *mux.Route) {
			r.Path("find_number_position").Methods(http.MethodGet)
		},
		Func:       g.findNumberPosition,
		Middleware: basic_auth.New(&basic_auth.BasicAuthConfig{AuthUser: g.config.AuthUser, AuthPassword: g.config.AuthPassword}),
	}
}

func (g *GatewayServer) findNumberPosition(writer http.ResponseWriter, request *http.Request) {
	numberParameter, err := validateNumberParameter(request)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "invalid number")

		return
	}

	number, err := g.gatewayService.FindNumberPosition(request.Context(), numberParameter)
	if err != nil {
		log.Errorf("Failed to find number position, %v", err)
		respondWithError(writer, http.StatusInternalServerError, "internal server error")

		return
	}

	respondWithJSON(writer, http.StatusOK, findNumberPositionResponse{Number: number})
}

func (g *GatewayServer) getUsers(writer http.ResponseWriter, request *http.Request) {
	users, err := g.gatewayService.GetUsers(request.Context())
	if err != nil {
		log.Errorf("Failed to get users, %v", err)
		respondWithError(writer, http.StatusInternalServerError, "internal server error")

		return
	}

	respondWithJSON(writer, http.StatusOK, users)
}

func validateNumberParameter(request *http.Request) (int, error) {
	numberParameter := strings.TrimSpace(request.URL.Query().Get("number"))
	if numberParameter == "" {
		return 0, errors.New("number is empty")
	}

	number, err := strconv.Atoi(numberParameter)
	if err != nil {
		return 0, errors.Wrap(err, "strconv Atoi")
	}

	return number, nil
}
