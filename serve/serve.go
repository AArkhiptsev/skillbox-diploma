package serve

import (
	"diploma/fetch"
	"diploma/lib"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
)

const serveraddr = "localhost:8282"

func sortByCountry(a []fetch.HeaderData) (rs []fetch.HeaderData) {

	sort.Slice(a, func(i, j int) bool {
		return a[i].Country < a[j].Country
	})
	rs = a
	return
}

func sortByProvider(a []fetch.HeaderData) (rs []fetch.HeaderData) {

	sort.Slice(a, func(i, j int) bool {
		return a[i].Provider < a[j].Provider
	})
	rs = a
	return
}

func prepareSMS() {

	a := fetch.StorageSMSData

	for i, _ := range a {
		a[i].Country = lib.GetCountryNameByAlpha(a[i].Country)
	}

	fetch.ResultSet.SMS = append(fetch.ResultSet.SMS, sortByCountry(a))

	lib.LogParseErr(0, "SMS с сортировкой по полным названия стран:")
	for _, data := range fetch.ResultSet.SMS[0] {
		fmt.Println(data)
	}

	fetch.ResultSet.SMS = append(fetch.ResultSet.SMS, sortByProvider(a))

	lib.LogParseErr(0, "SMS с сортировкой по провайдерам:")
	for _, data := range fetch.ResultSet.SMS[1] {
		fmt.Println(data)
	}
}

func prepareMMS() {

	a := fetch.StorageMMSData

	for i, _ := range a {
		a[i].Country = lib.GetCountryNameByAlpha(a[i].Country)
	}

	fetch.ResultSet.MMS = append(fetch.ResultSet.MMS, sortByCountry(a))

	lib.LogParseErr(0, "MMS с сортировкой по полным названия стран:")
	for _, data := range fetch.ResultSet.MMS[0] {
		fmt.Println(data)
	}

	fetch.ResultSet.MMS = append(fetch.ResultSet.MMS, sortByProvider(a))

	lib.LogParseErr(0, "MMS с сортировкой по провайдерам:")
	for _, data := range fetch.ResultSet.MMS[1] {
		fmt.Println(data)
	}
}

func GetResultData() {

	prepareSMS()
	prepareMMS()
	fetch.ResultSet.VoiceCall = fetch.StorageVoiceCallData
	fetch.ResultSet.Billing = fetch.StorageBilling

	return
}

func handleServer(w http.ResponseWriter, r *http.Request) {
	return
}

func ListenAndServeHTTP() {

	router := mux.NewRouter()

	router.HandleFunc("/", handleServer)

	lib.LogParseErr(1,
		fmt.Sprintf("Запускаю сервер %s", serveraddr))

	http.ListenAndServe(serveraddr, router)
}
