package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NMEJIA93/go_lib_response/response"
	"net/http"
	"strconv"

	"github.com/NMEJIA93/gocourse_user/src/user"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewUserHTTPServer(cxt context.Context, endpoints user.Endpoints) http.Handler {

	r := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Handle("/user", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateUserRequest,
		encodeResponse,
		options...,
	)).Methods("POST")

	r.Handle("/user/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetUSer,
		encodeResponse,
		options...,
	)).Methods("GET")

	r.Handle("/user", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllUSer,
		encodeResponse,
		options...,
	)).Methods("GET")

	return r
}

func decodeCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req user.CreateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprint("invalid request format:  '%v'", err.Error()))
	}
	return req, nil

}

func encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	//json.NewEncoder(w).Encode(resp)
	_ = json.NewEncoder(w).Encode(resp)

}

func decodeGetUSer(_ context.Context, r *http.Request) (interface{}, error) {
	p := mux.Vars(r)
	req := user.GetReq{
		ID: p["id"],
	}

	return req, nil
}

func decodeGetAllUSer(_ context.Context, r *http.Request) (interface{}, error) {
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := user.GetAllReq{
		FirstName: v.Get("first_name"),
		LastName:  v.Get("last_name"),
		Limit:     limit,
		Page:      page,
	}

	return req, nil
}
