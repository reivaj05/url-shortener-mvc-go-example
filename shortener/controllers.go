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
	json.CreateJSONArrayAtPathWithArray("data", convertListToJSON(list))
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
	shortURL := params["shortUrl"]
	if err := dbClient.GetWhere(
		&urlShortInstance, "short_url", shortURL); err != nil {

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
	item, err := saveDataIntoDB(data)
	if err != nil {
		GoServer.SendResponseWithStatus(rw, "", http.StatusInternalServerError)
		return
	}
	GoServer.SendResponseWithStatus(
		rw, item.ToJSON().ToString(), http.StatusCreated)
}

func putItemHandler(rw http.ResponseWriter, req *http.Request) {
	data, err := getJSONData(req)
	if err != nil {
		sendBadRequestResponse(rw)
		return
	}
	params := GoServer.GetQueryParams(req)
	id, _ := strconv.Atoi(params["id"])
	if !isDataValid(data) {
		sendBadRequestResponse(rw)
		return
	}
	data.SetValueAtPath("id", id)
	item, err := saveDataIntoDB(data)
	if err != nil {
		GoServer.SendResponseWithStatus(rw, "", http.StatusInternalServerError)
		return
	}
	GoServer.SendResponseWithStatus(
		rw, item.ToJSON().ToString(), http.StatusOK)
}

func getJSONData(req *http.Request) (data *GoJSON.JSONWrapper, err error) {
	body, _ := GoServer.ReadBodyRequest(req)
	data, err = GoJSON.New(body)
	return
}

func isDataValid(data *GoJSON.JSONWrapper) bool {
	if data.HasPath("longUrl") {
		if longURL, ok := data.GetStringFromPath("longUrl"); ok {
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

func saveDataIntoDB(data *GoJSON.JSONWrapper) (*URLShortModel, error) {
	longURL, _ := data.GetStringFromPath("longUrl")
	shortURL := createShortURL(longURL)
	dbClient := GoDB.GetDBClient()
	if data.HasPath("id") {
		id, _ := data.GetIntFromPath("id")
		return updateURL(longURL, shortURL, id, dbClient)
	}
	return createURL(longURL, shortURL, dbClient)
}

func updateURL(longURL, shortURL string, id int,
	dbClient *GoDB.DBClient) (*URLShortModel, error) {

	var item URLShortModel
	if err := dbClient.Get(&item, id); err != nil {
		return nil, err
	}
	item.LongURL = longURL
	item.ShortURL = shortURL
	return &item, dbClient.Update(&item)
}

func createURL(longURL, shortURL string,
	dbClient *GoDB.DBClient) (*URLShortModel, error) {

	item := URLShortModel{
		LongURL:  longURL,
		ShortURL: shortURL,
	}
	return &item, dbClient.Create(&item)
}

func createShortURL(longURL string) string {
	// TODO: create short_url
	// TODO: check if longURL exists and send appropiate response
	return longURL
}

func deleteItemHandler(rw http.ResponseWriter, req *http.Request) {
	params := GoServer.GetQueryParams(req)
	id, _ := strconv.Atoi(params["id"])
	var url URLShortModel
	dbClient := GoDB.GetDBClient()
	notDeletedMsg := `{"error": "The resource was not deleted"}`
	if err := dbClient.Get(&url, id); err != nil {
		GoServer.SendResponseWithStatus(rw, GoServer.ResourceNotFound,
			http.StatusInternalServerError)
		return
	}
	if err := dbClient.Delete(&url); err != nil {
		GoServer.SendResponseWithStatus(rw, notDeletedMsg,
			http.StatusInternalServerError)
		return
	}
	GoServer.SendResponseWithStatus(rw, "", http.StatusOK)
}
