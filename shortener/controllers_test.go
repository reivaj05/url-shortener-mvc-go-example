package shortener

import (
	"bytes"
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
	router.Methods("GET").Path("/shortener/").HandlerFunc(getListHandler)
	router.Methods("GET").Path("/shortener/{id:[0-9]+}").HandlerFunc(getItemHandler)
	router.Methods("POST").Path("/shortener/").HandlerFunc(postItemhandler)
	return router
}

func (suite *ControllersTestSuite) TestGetListHandler() {

	suite.dbClient.Create(&URLShortModel{LongURL: "long1", ShortURL: "short1"})
	suite.dbClient.Create(&URLShortModel{LongURL: "long2", ShortURL: "short2"})
	response, status := suite.makeRequest("GET", "http://localhost/shortener/", nil)
	suite.assert.Equal(http.StatusOK, status)
	jsonResponse, _ := GoJSON.New(response)
	suite.assert.True(jsonResponse.HasPath("data"))
	dataArray := jsonResponse.GetArrayFromPath("data")
	suite.assert.Equal(2, len(dataArray))
}

func (suite *ControllersTestSuite) TestGetItemHandler() {
	suite.dbClient.Create(&URLShortModel{})
	_, status := suite.makeRequest("GET", "http://localhost/shortener/1", nil)
	suite.assert.Equal(http.StatusOK, status)
}

func (suite *ControllersTestSuite) TestGetItemHandlerNotFound() {
	_, status := suite.makeRequest("GET", "http://localhost/shortener/1", nil)
	suite.assert.Equal(http.StatusNotFound, status)
}

func (suite *ControllersTestSuite) TestPostItemHandler() {
	const body = `{"longURL": "http://www.example.com"}`
	_, status := suite.makeRequest("POST", "http://localhost/shortener/", []byte(body))
	suite.assert.Equal(http.StatusCreated, status)
}

func (suite *ControllersTestSuite) TestPostItemHandlerWithWrongJSON() {
	_, status := suite.makeRequest("POST", "http://localhost/shortener/", []byte(""))
	suite.assert.Equal(http.StatusBadRequest, status)
}

func (suite *ControllersTestSuite) TestPostItemHandlerWithoutCorrectFields() {
	_, status := suite.makeRequest("POST", "http://localhost/shortener/", []byte("{}"))
	suite.assert.Equal(http.StatusBadRequest, status)
}

func (suite *ControllersTestSuite) TestPostItemHandlerWithInvalidData() {
	const body = `{"longURL": "not a url"}`
	_, status := suite.makeRequest("POST", "http://localhost/shortener/", []byte(body))
	suite.assert.Equal(http.StatusBadRequest, status)
}

func (suite *ControllersTestSuite) TestPostItemHandlerWithRepeatedData() {
	suite.dbClient.Create(&URLShortModel{LongURL: "http://www.example.com", ShortURL: ""})
	const body = `{"longURL": "http://www.example.com"}`
	_, status := suite.makeRequest("POST", "http://localhost/shortener/", []byte(body))
	suite.assert.Equal(http.StatusInternalServerError, status)
}

func (suite *ControllersTestSuite) makeRequest(
	method, url string, body []byte) (string, int) {

	rw := httptest.NewRecorder()
	var request *http.Request
	if body != nil {
		request, _ = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		request, _ = http.NewRequest(method, url, nil)
	}

	suite.router.ServeHTTP(rw, request)
	res, _ := ioutil.ReadAll(rw.Body)
	return string(res), rw.Code
}

func TestControllers(t *testing.T) {
	suite.Run(t, new(ControllersTestSuite))
}
