package driving

import (
	"fmt"
	"net/http"
)

type authorControllerDaemon struct{}

func (authorControllerDaemon) create(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusAccepted)
}

func (authorControllerDaemon) query(writer http.ResponseWriter, request *http.Request) {
	authorName := request.URL.Query().Get("name")

	response := fmt.Appendf(nil, `{"name":"%s"}`, authorName)

	writer.WriteHeader(http.StatusAccepted)
	writer.Write(response)
}

func NewAuthorController() *http.ServeMux {
	var router http.ServeMux
	controller := authorControllerDaemon{}

	router.HandleFunc("GET /author", controller.query)
	router.HandleFunc("POST /author", controller.create)

	return &router
}
