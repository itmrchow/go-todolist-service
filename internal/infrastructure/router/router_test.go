package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// RouterTestSuite 定義測試套件
type RouterTestSuite struct {
	suite.Suite
	router *RouterImpl
	engine *gin.Engine
}

// SetupSuite 在所有測試開始前執行一次
func (suite *RouterTestSuite) SetupSuite() {
	// 設定 Gin 為測試模式
	gin.SetMode(gin.TestMode)
}

// SetupTest 在每個測試前執行
func (suite *RouterTestSuite) SetupTest() {
	suite.router = &RouterImpl{}
	suite.engine = suite.router.SetupRoutes()
}

func (suite *RouterTestSuite) TestSetupRoutes() {
	assert.NotNil(suite.T(), suite.engine)
	assert.IsType(suite.T(), &gin.Engine{}, suite.engine)
}

func (suite *RouterTestSuite) TestHealthEndpoint() {
	// 測試 /health 端點
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "healthy")
	assert.Contains(suite.T(), w.Body.String(), "todolist-service")
}

func (suite *RouterTestSuite) TestVersionEndpoint() {
	// 測試 /version 端點
	req, _ := http.NewRequest("GET", "/version", nil)
	w := httptest.NewRecorder()
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "version")
	assert.Contains(suite.T(), w.Body.String(), "todolist-service")
}

func (suite *RouterTestSuite) TestV1RouteGroup() {
	// 測試 /api/v1 路由群組是否存在
	// 這裡測試一個基本的 404，確保路由系統工作正常
	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	w := httptest.NewRecorder()
	suite.engine.ServeHTTP(w, req)

	// 應該是 404 因為我們還沒有定義 /api/v1/test 路由
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *RouterTestSuite) TestRegisterV1Routes() {
	engine := gin.New()
	v1Group := engine.Group("/api/v1")

	// 這個測試確保 RegisterV1Routes 方法不會出錯
	assert.NotPanics(suite.T(), func() {
		suite.router.RegisterV1Routes(v1Group)
	})
}

// TestRouterSuite 執行測試套件
func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}
