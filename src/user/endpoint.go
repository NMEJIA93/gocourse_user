package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/NMEJIA93/go_lib_response/response"
	"github.com/NMEJIA93/gocourse_meta/meta"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
		//Get    Controller
		//GetAll Controller
		//Update Controller
		//Delete Controller
	}
	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	UpdateReq struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	ErrorResp struct {
		Error string `json:"error"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Error  string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}
	Config struct {
		LimitPageDef string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		//Get:    makeGetEndpoint(s),
		//GetAll: makeGetAllEndpoint(s, config),
		//Update: makeUpdateEndpoint(s),
		//Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(CreateReq)

		if req.FirstName == "" {
			return nil, response.BadRequest(fmt.Sprint(errors.New("first name is required")))
		}

		if req.LastName == "" {
			return nil, response.BadRequest(fmt.Sprint("last name is required"))
		}

		userRequest := CreateUserDTO{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
		}

		user, serviceErr := s.Create(ctx, userRequest)
		if serviceErr != nil {
			return nil, response.InternalServerError(fmt.Sprint(serviceErr))
		}

		userResponse := ResponseUserDto{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Phone:     user.Phone,
		}

		fmt.Println("Create user: ", userResponse.ID)

		return response.Created("success", userResponse, nil), nil
	}
}

/*

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		err := s.Delete(id)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Error:  "User not found",
			})
			return
		}

		fmt.Println("Delete user")
		json.NewEncoder(w).Encode(&Response{
			Status: 200,
			Data:   "Success",
		})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		path := mux.Vars(r)
		id := path["id"]
		user, err := s.Get(id)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Error:  "User not found",
			})
			return
		}
		fmt.Println("Get user")
		json.NewEncoder(w).Encode(&Response{
			Status: 200,
			Data:   user,
		})
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		v := r.URL.Query()
		filters := Filters{
			FirstName: v.Get("first_name"),
			LastName:  v.Get("last_name"),
		}
		limit, _ := strconv.Atoi(v.Get("limit"))
		page, _ := strconv.Atoi(v.Get("page"))

		count, err := s.Count(filters)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{
				Status: 500,
				Error:  err.Error(),
			})
		}
		meta, err := meta.New(page, limit, count, config.LimitPageDef)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{
				Status: 500,
				Error:  err.Error(),
			})
		}

		users, err := s.GetAll(filters, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Error:  err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(&Response{
			Status: 200,
			Data:   users,
			Meta:   meta,
		})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Update user")

		var req UpdateReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Error:  "Invalid request",
			})
			return
		}

		if req.FirstName != nil && *req.FirstName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Error:  "first name is required",
			})
			return
		}

		if req.LastName != nil && *req.LastName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Error:  "last name is required",
			})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		err = s.Update(id, req.FirstName, req.LastName, req.Email, req.Phone)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Error:  "user not found",
			})
			return
		}

		json.NewEncoder(w).Encode(&Response{
			Status: 200,
			Data:   "Success",
		})
	}
}

*/
