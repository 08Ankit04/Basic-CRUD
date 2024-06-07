package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Basic-CRUD/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
)

const (
	defaultPageLimit  = 5
	defaultPageNumber = 1
)

func handlePaginationQuery(r *http.Request) (int64, int64) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = defaultPageLimit
	}

	pageNumber, err := strconv.Atoi(r.URL.Query().Get("pageNumber"))
	if err != nil {
		pageNumber = defaultPageNumber
	}

	if limit < 0 {
		limit = defaultPageLimit
	}

	if pageNumber < defaultPageNumber {
		pageNumber = defaultPageNumber
	}

	offset := (pageNumber - defaultPageNumber) * limit

	return int64(limit), int64(offset)
}

// HandleListEmployee godoc
//
//	@summary		list employee
//	@description	Get employee list
//	@tags			employee
//
//	@router			/employee [GET]
//	@accept			json
//	@produce		json
//	@param			limit	query	string	false	"number of employee details in one page" default(5)
//	@param			pageNumber	query	string	false	"page number of employee" default(1)
//
//	@success		200				{array}	    model.Employee
//	@failure		500				"Internal Server Error"
func (srv *Server) HandleListEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset := handlePaginationQuery(r)

	keys, err := srv.Redis.Keys(ctx, "employee:*").Result()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if offset >= int64(len(keys)) {
		http.Error(w, "No more items", http.StatusNotFound)
		return
	}

	end := offset + limit

	if end >= int64(len(keys)) {
		end = int64(len(keys))
	}

	employees := make([]model.Employee, 0, limit)
	for _, key := range keys[offset:end] {
		data, err := srv.Redis.Get(ctx, key).Result()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var employee model.Employee
		if err := json.Unmarshal([]byte(data), &employee); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		employees = append(employees, employee)
	}

	if err := json.NewEncoder(w).Encode(employees); err != nil {
		srv.Logger.Print(err)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// HandleGetEmployee godoc
//
//	@summary		get employee details
//	@description	Get employee by id
//	@tags			employee
//
//	@router			/employee/{id} [GET]
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"employee ID in int format"
//
//	@success		200				{object}	model.Employee
//	@failure		404				"Not Found"
//	@failure		422				"Unprocessable entity"
//	@failure		500				"Internal Server Error"
func (srv *Server) HandleGetEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Unprocessable entity", http.StatusUnprocessableEntity)
		return
	}

	data, err := srv.Redis.Get(ctx, "employee:"+strconv.Itoa(id)).Result()
	if err != nil && err != redis.Nil {
		srv.Logger.Print(err)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if data == "" {
		http.Error(w, "Status Not Found", http.StatusNotFound)
		return
	}

	var employee model.Employee
	if err := json.Unmarshal([]byte(data), &employee); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&employee); err != nil {
		srv.Logger.Print(err)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// HandleCreateEmployee godoc
//
//	@summary		create employee details
//	@description	create employee
//	@tags			employee
//
//	@router			/employee [POST]
//	@accept			json
//	@produce		json
//	@param			model.Employee	body	model.Employee	true	"employee details in json format"
//
//	@success		201				{object}	model.Employee
//	@failure		404				"Not Found"
//	@failure		409				"Conflict"
//	@failure		500				"Internal Server Error"
func (srv *Server) HandleCreateEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var employee model.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(employee)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	existingEmployee, err := srv.Redis.Get(ctx, fmt.Sprintf("employee:%d", employee.Id)).Result()
	if err != nil && err != redis.Nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if existingEmployee != "" {
		http.Error(w, "Conflict", http.StatusConflict)
		return
	}

	err = srv.Redis.Set(ctx, fmt.Sprintf("employee:%d", employee.Id), data, 0).Err()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(employee)
}

// HandleUpdateEmployee godoc
//
//	@summary		update employee details
//	@description	update employee details by id
//	@tags			employee
//
//	@router			/employee/{id} [PUT]
//	@accept			json
//	@produce		json
//	@param			model.Employee	body	model.Employee	true	"employee details in json format"
//	@param			id	path		string	true	"employee ID in int format"
//
//	@success		200				{object}	model.Employee
//	@failure		404				"Not Found"
//	@failure		422				"Unprocessable entity"
//	@failure		500				"Internal Server Error"
func (srv *Server) HandleUpdateEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Unprocessable entity", http.StatusUnprocessableEntity)
		return
	}

	var employee model.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(employee)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = srv.Redis.Set(ctx, fmt.Sprintf("employee:%d", id), data, 0).Err()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(employee)
}

// HandleDeleteEmployee godoc
//
//	@summary		delete employee details
//	@description	delete employee details by id
//	@tags			employee
//
//	@router			/employee/{id} [DELETE]
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"employee ID in int format"
//
//	@success		200				"Delete Successful"
//	@failure		422				"Unprocessable entity"
//	@failure		500				"Internal Server Error"
func (srv *Server) HandleDeleteEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Unprocessable entity", http.StatusUnprocessableEntity)
		return
	}

	err = srv.Redis.Del(ctx, fmt.Sprintf("employee:%d", id)).Err()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
