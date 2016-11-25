package shortener

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/reivaj05/GoDB"
	"github.com/reivaj05/GoJSON"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ControllersTestSuite struct {
	suite.Suite
	assert   *assert.Assertions
	router   *mux.Router
	dbClient *GoDB.DBClient
}

func (suite *ControllersTestSuite) SetupSuite() {
	suite.assert = assert.New(suite.T())
	suite.router = createTestRouter()
}

func (suite *ControllersTestSuite) SetupTest() {
	GoDB.Init(&GoDB.DBOptions{
		DBEngine: "sqlite3",
		DBName:   "test.db",
	}, Models...)

	suite.dbClient = GoDB.GetDBClient()
}

func (suite *ControllersTestSuite) TearDownTest() {
	os.Remove("test.db")
}

func createTestRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/shortener/", getListHandler)
	router.HandleFunc("/shortener/{id:[0-9]+}", getItemHandler)
	return router
}

func (suite *ControllersTestSuite) TestGetListHandler() {

	suite.dbClient.Create(&URLShortModel{})
	suite.dbClient.Create(&URLShortModel{})
	response, status := suite.makeRequest("GET", "http://localhost/shortener/")
	suite.assert.Equal(http.StatusOK, status)
	jsonResponse, _ := GoJSON.New(response)
	suite.assert.True(jsonResponse.HasPath("data"))
	dataArray := jsonResponse.GetArrayFromPath("data")
	suite.assert.Equal(2, len(dataArray))
}

func (suite *ControllersTestSuite) TestGetItemHandler() {
	suite.dbClient.Create(&URLShortModel{})
	_, status := suite.makeRequest("GET", "http://localhost/shortener/1")
	suite.assert.Equal(http.StatusOK, status)
}

func (suite *ControllersTestSuite) TestGetItemHandlerNotFound() {
	_, status := suite.makeRequest("GET", "http://localhost/shortener/1")
	suite.assert.Equal(http.StatusNotFound, status)
}

func (suite *ControllersTestSuite) makeRequest(method, url string) (string, int) {
	rw := httptest.NewRecorder()
	request, _ := http.NewRequest(method, url, nil)
	suite.router.ServeHTTP(rw, request)
	res, _ := ioutil.ReadAll(rw.Body)
	return string(res), rw.Code
}

func TestControllers(t *testing.T) {
	suite.Run(t, new(ControllersTestSuite))
}
