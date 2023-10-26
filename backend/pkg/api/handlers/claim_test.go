package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"wallemon/pkg/database"
	"wallemon/pkg/docker"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type claimSuit struct {
	suite.Suite
	ctl *gomock.Controller
}

func (s *claimSuit) SetupSuite() {
	fmt.Println("setup suite!!!!!!!!!!!!!!!!!!!!!!")
	host, port, err := docker.Run("postgres")
	s.Require().NoError(err)
	endpoint := fmt.Sprintf("postgres://postgres:postgres@%v:%v/postgres?sslmode=disable", host, port)
	fmt.Println(endpoint)
	database.Initialize(context.Background())

	s.ctl = gomock.NewController(s.T())
}

func (s *claimSuit) TearDownSuite() {
	database.Finalize()
	s.NoError(docker.Remove())
	s.ctl.Finish()
}

func (s *claimSuit) SetupTest() {
}

func TestClaim(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("SIGNER_KEY", "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	endpoint := "/claim"
	router := gin.Default()
	router.POST(endpoint, Claim)

	t.Run("successful claim", func(t *testing.T) {
		claimReq := ClaimReq{
			Address: "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
			RefCode: "wallemon",
		}

		body, _ := json.Marshal(&claimReq)
		req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)

		var claimResp ClaimResp
		err := json.Unmarshal(resp.Body.Bytes(), &claimResp)
		assert.NoError(t, err)

		assert.Equal(t, "f4b6424aebb6e151136076cfe601f8c20cb91603e4bfef367c384f7f5de6fd287c707cdce8aa7757770f756e0b6a4351998861bfb9417de77ab365181625e7af1b", claimResp.Signature)
	})

	t.Run("invalid address", func(t *testing.T) {
		claimReq := ClaimReq{
			Address: "0x1",
			RefCode: "wallemon",
		}

		body, _ := json.Marshal(&claimReq)
		req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid ref code", func(t *testing.T) {
		claimReq := ClaimReq{
			Address: "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
			RefCode: "invalid",
		}

		body, _ := json.Marshal(&claimReq)
		req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestJoinWaitlist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	endpoint := "/joinWaitlist"
	router := gin.Default()
	router.POST(endpoint, JoinWaitlist)

	t.Run("successful join waitlist", func(t *testing.T) {
		joinWaitlistReq := JoinWaitlistReq{
			Email: "user@exmaple.com",
		}

		body, _ := json.Marshal(&joinWaitlistReq)
		req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)

		var joinWaitlistResp JoinWaitlistResp
		err := json.Unmarshal(resp.Body.Bytes(), &joinWaitlistResp)
		assert.NoError(t, err)
		assert.Equal(t, "success", joinWaitlistResp.Message)
	})
}

func (s *claimSuit) migrateUp() {}

func (s *claimSuit) TearDownTest() {}

func TestClaimSuite(t *testing.T) {
	fmt.Println("TestClaimSuite!!!!!!!!!!!!!!!!!!!!!!")
	suite.Run(t, new(claimSuit))
}
