package server

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Basic-CRUD/model"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(key string) *redis.StringCmd {
	args := m.Called(key)
	data, _ := json.Marshal(args.Get(0).(model.Employee))
	return redis.NewStringResult(string(data), args.Error(1))
}

func (m *MockRedisClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(key, value, expiration)
	return redis.NewStatusResult(args.String(0), args.Error(1))
}

func (m *MockRedisClient) Del(key string) *redis.IntCmd {
	args := m.Called(key)
	return redis.NewIntResult(int64(args.Int(0)), args.Error(1))
}

func (m *MockRedisClient) Keys(pattern string) *redis.StringSliceCmd {
	args := m.Called(pattern)
	return redis.NewStringSliceResult(args.Get(0).([]string), args.Error(1))
}

type EmployeeControllerTestSuite struct {
	suite.Suite
	Srv   *Server
	redis *MockRedisClient
}

func TestAllEmployeeControllers(t *testing.T) {
	suite.Run(t, new(EmployeeControllerTestSuite))
}

func (st *EmployeeControllerTestSuite) SetupTest() {
	rc := new(MockRedisClient)
	validator := validator.New()

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	st.Srv = New(rc, logger, validator)

	st.redis = rc
}
