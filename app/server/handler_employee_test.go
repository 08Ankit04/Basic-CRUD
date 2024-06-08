package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/Basic-CRUD/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
)

func (st *EmployeeControllerTestSuite) TestGetEmployeeHandler() {
	st.redis.On("Get", "employee:1").Return(model.Employee{Id: 1, Name: "John Doe", Position: "Developer", Salary: 50000}, nil)

	req, err := http.NewRequest("GET", "/employee/1", nil)
	st.Require().Nil(err)

	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Get("/employee/{id}", st.Srv.HandleGetEmployee)

	r.ServeHTTP(rr, req)

	st.Assert().Equal(http.StatusOK, rr.Code)

	expected := `{"id":1,"name":"John Doe","position":"Developer","salary":50000}`
	actual := strings.TrimSuffix(rr.Body.String(), "\n")

	st.Assert().Equal(expected, actual)
}

func (st *EmployeeControllerTestSuite) TestCreateEmployeeHandler() {
	employee := model.Employee{Id: 1, Name: "John Doe", Position: "Developer", Salary: 50000}
	employeeJSON, _ := json.Marshal(employee)

	st.redis.On("Get", "employee:1").Return(model.Employee{}, redis.Nil)

	st.redis.On("Set", "employee:1", employeeJSON, time.Duration(0)).Return("", nil)

	req, err := http.NewRequest("POST", "/employee", bytes.NewBuffer(employeeJSON))
	st.Require().Nil(err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(st.Srv.HandleCreateEmployee)

	handler.ServeHTTP(rr, req)

	st.Assert().Equal(http.StatusCreated, rr.Code)

	expected := string(employeeJSON)
	actual := strings.TrimSuffix(rr.Body.String(), "\n")

	st.Assert().Equal(expected, actual)
}

func (st *EmployeeControllerTestSuite) TestUpdateEmployeeHandler() {
	employee := model.Employee{Id: 1, Name: "John Doe", Position: "Developer", Salary: 50000}
	employeeJSON, _ := json.Marshal(employee)

	st.redis.On("Set", "employee:1", employeeJSON, time.Duration(0)).Return("", nil)

	req, err := http.NewRequest("PUT", "/employee/1", bytes.NewBuffer(employeeJSON))
	st.Require().Nil(err)

	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Put("/employee/{id}", st.Srv.HandleUpdateEmployee)

	r.ServeHTTP(rr, req)

	st.Assert().Equal(http.StatusOK, rr.Code)

	expected := string(employeeJSON)
	actual := strings.TrimSuffix(rr.Body.String(), "\n")

	st.Assert().Equal(expected, actual)
}

func (st *EmployeeControllerTestSuite) TestDeleteEmployeeHandler() {
	st.redis.On("Del", "employee:1").Return(1, nil)

	req, err := http.NewRequest("DELETE", "/employee/1", nil)
	st.Require().Nil(err)

	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Delete("/employee/{id}", st.Srv.HandleDeleteEmployee)

	r.ServeHTTP(rr, req)

	st.Assert().Equal(http.StatusOK, rr.Code)
}
