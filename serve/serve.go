package serve

import (
	"diploma/lib"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const serveraddr = "localhost:8282"

func handleMMS(w http.ResponseWriter, r *http.Request) {
	return
}

func ListenAndServeHTTP() {

	router := mux.NewRouter()

	router.HandleFunc("/", handleMMS)

	lib.LogParseErr(1,
		fmt.Sprintf("Запускаю сервер %s", serveraddr))

	http.ListenAndServe(serveraddr, router)
}
