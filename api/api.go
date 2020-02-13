package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/InVisionApp/interview-test/api/dal"
	"github.com/InVisionApp/rye"
	"github.com/gorilla/mux"
)

const (
	callingService = "example-bff"
)

type API struct {
	listenAddress string
	MWHandler     *rye.MWHandler
	apiDAL        dal.DAL
}

func New(addr string) *API {
	if addr == "" {
		addr = ":31337"
	}
	return &API{
		listenAddress: addr,
		MWHandler:     rye.NewMWHandler(rye.Config{}),
		apiDAL:        dal.NewMockDAL(),
	}
}

func (a *API) Run() error {
	log.Println("Starting API server...")

	routes := mux.NewRouter().StrictSlash(true)

	a.MWHandler.Use(mwRouteLogger())
	a.MWHandler.Use(mwCallingService(callingService))

	routes.Handle("/api/v1/login",
		a.MWHandler.Handle([]rye.Handler{a.loginHandler})).Methods("POST")
	routes.Handle("/api/v1/teams/{TeamID}",
		a.MWHandler.Handle([]rye.Handler{a.getTeamHandler})).Methods("GET")

	log.Printf("API server running on %s\n", a.listenAddress)

	return http.ListenAndServe(a.listenAddress, routes)
}

func mwCallingService(service string) func(rw http.ResponseWriter, req *http.Request) *rye.Response {
	return func(rw http.ResponseWriter, r *http.Request) *rye.Response {
		callingService := r.Header.Get("Calling-Service")
		if callingService != service {
			return &rye.Response{
				Err:        errors.New("Caller not whitelisted"),
				StatusCode: http.StatusForbidden,
			}
		}
		return nil
	}
}

func mwRouteLogger() func(rw http.ResponseWriter, req *http.Request) *rye.Response {
	return func(rw http.ResponseWriter, r *http.Request) *rye.Response {
		log.Printf("%s \"%s %s %s\"", r.RemoteAddr, r.Method, r.RequestURI, r.Proto)
		return nil
	}
}

func respondAsJSON(rw http.ResponseWriter, code int, v interface{}) *rye.Response {
	body, err := json.Marshal(v)
	if err != nil {
		return &rye.Response{
			Err:        fmt.Errorf("Unable to generate response JSON: %v", err),
			StatusCode: http.StatusInternalServerError,
		}
	}

	rye.WriteJSONResponse(rw, code, body)

	return nil
}
