package shortener

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/reivaj05/GoDB"
	"github.com/reivaj05/GoJSON"
	"github.com/reivaj05/GoServer"
)

func getListHandler(rw http.ResponseWriter, req *http.Request) {
	var urlList []URLShortModel
	dbClient := GoDB.GetDBClient()
	if err := dbClient.List(&urlList); err != nil {
		GoServer.SendResponseWithStatus(rw, "", http.StatusInternalServerError)
		return
	}
	jsonResponse := createJSONResponse(urlList)
	GoServer.SendResponseWithStatus(rw, jsonResponse, http.StatusOK)
}

func createJSONResponse(list []URLShortModel) string {
	json, _ := GoJSON.New("{}")
	json.CreateJSONArrayAtPath("data", convertListToJSON(list))
	return json.ToString()
}

func convertListToJSON(list []URLShortModel) (jsonList []*GoJSON.JSONWrapper) {
	for _, item := range list {
		jsonList = append(jsonList, item.ToJSON())
	}
	return
}

func getItemHandler(rw http.ResponseWriter, req *http.Request) {
	var urlShortInstance URLShortModel
	dbClient := GoDB.GetDBClient()
	params := GoServer.GetQueryParams(req)
	id, _ := strconv.Atoi(params["id"])
	if err := dbClient.Get(&urlShortInstance, id); err != nil {
		GoServer.SendResponseWithStatus(
			rw, GoServer.ResourceNotFound, http.StatusNotFound)
		return
	}
	GoServer.SendResponseWithStatus(
		rw, urlShortInstance.ToJSON().ToString(), http.StatusOK)
}

func postItemhandler(rw http.ResponseWriter, req *http.Request) {
	data, err := getJSONData(req)
	if err != nil {
		sendBadRequestResponse(rw)
		return
	}
	if !isDataValid(data) {
		sendBadRequestResponse(rw)
		return
	}
	if err := saveDataIntoDB(data); err != nil {
		GoServer.SendResponseWithStatus(rw, "", http.StatusInternalServerError)
		return
	}
	GoServer.SendResponseWithStatus(
		rw, GoServer.ResourceCreated, http.StatusCreated)
}

func getJSONData(req *http.Request) (data *GoJSON.JSONWrapper, err error) {
	body, _ := GoServer.ReadBodyRequest(req)
	data, err = GoJSON.New(body)
	return
}

func isDataValid(data *GoJSON.JSONWrapper) bool {
	if data.HasPath("longURL") {
		if longURL, ok := data.GetStringFromPath("longURL"); ok {
			return isURLValid(longURL)
		}
	}
	return false
}

func isURLValid(longURL string) bool {
	_, err := url.ParseRequestURI(longURL)
	return err == nil
}

func sendBadRequestResponse(rw http.ResponseWriter) {
	GoServer.SendResponseWithStatus(
		rw, GoServer.BadRequest, http.StatusBadRequest)
}

func saveDataIntoDB(data *GoJSON.JSONWrapper) error {
	longURL, _ := data.GetStringFromPath("longURL")
	shortURL := createShortURL(longURL)
	dbClient := GoDB.GetDBClient()
	return dbClient.Create(&URLShortModel{
		LongURL:  longURL,
		ShortURL: shortURL,
	})
}

func createShortURL(longURL string) string {
	// TODO: create short_url
	// TODO: check if longURL exists and send appropiate response
	return "short_url_of: " + longURL
}

func putItemHandler(rw http.ResponseWriter, req *http.Request) {
	GoServer.SendResponseWithStatus(
		rw, `{"msg": "implement put"}`, http.StatusOK)
}

func deleteItemHandler(rw http.ResponseWriter, req *http.Request) {
	GoServer.SendResponseWithStatus(
		rw, `{"msg": "implement delete"}`, http.StatusOK)
}
