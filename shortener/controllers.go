package shortener

import (
	"net/http"
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
			rw, `{"error": "Resource not found"}`, http.StatusNotFound)
		return
	}
	GoServer.SendResponseWithStatus(
		rw, urlShortInstance.ToJSON().ToString(), http.StatusOK)
}

func postItemhandler(rw http.ResponseWriter, req *http.Request) {
	GoServer.SendResponseWithStatus(
		rw, `{"msg": "implement post"}`, http.StatusOK)
}

func putItemHandler(rw http.ResponseWriter, req *http.Request) {
	GoServer.SendResponseWithStatus(
		rw, `{"msg": "implement put"}`, http.StatusOK)
}

func deleteItemHandler(rw http.ResponseWriter, req *http.Request) {
	GoServer.SendResponseWithStatus(
		rw, `{"msg": "implement delete"}`, http.StatusOK)
}
