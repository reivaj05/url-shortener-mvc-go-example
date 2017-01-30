package shortener

import (
	"github.com/reivaj05/GoServer"
)

var Endpoints = []*GoServer.Endpoint{
	&GoServer.Endpoint{
		Method:  "GET",
		Path:    "/shortener/",
		Handler: getListHandler,
	},
	&GoServer.Endpoint{
		Method: "GET",
		// Path:    "/shortener/{id:[0-9]+}",
		Path:    "/shortener/{shortUrl}",
		Handler: getItemHandler,
	},
	&GoServer.Endpoint{
		Method:  "POST",
		Path:    "/shortener/",
		Handler: postItemhandler,
	},
	&GoServer.Endpoint{
		Method:  "PUT",
		Path:    "/shortener/{id:[0-9]+}",
		Handler: putItemHandler,
	},
	&GoServer.Endpoint{
		Method:  "DELETE",
		Path:    "/shortener/{id:[0-9]+}",
		Handler: deleteItemHandler,
	},
}
