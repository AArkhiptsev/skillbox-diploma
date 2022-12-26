package serve

import (
	"diploma/fetch"
	"diploma/lib"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"strconv"
)

const (
	serveraddr             = "localhost:8282"
	averageSupportBandwith = 18
	lowLoad                = 9
	middleLoad             = 16
	//specCount              = 7
)

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

func sortByAccident(a []fetch.AccidentData) (rs []fetch.AccidentData) {

	sort.Slice(a, func(i, j int) bool {
		return a[i].Status < a[j].Status
	})
	rs = a
	return
}

func sortByCountryAndSpeed(a []fetch.EmailData) (rs []fetch.EmailData) {

	sort.Slice(a, func(i, j int) bool {
		iv, jv := a[i], a[j]
		switch {
		case iv.Country != jv.Country:
			return iv.Country < jv.Country
		case iv.DeliveryTime != jv.DeliveryTime:
			return iv.DeliveryTime < jv.DeliveryTime
		default:
			return iv.Country < jv.Country
		}
	})
	rs = a
	return
}

func prepareSMS() (result int) {

	a := fetch.StorageSMSData

	if len(a) == 0 {
		lib.LogParseErr(2, "SMS нет")
		result++
		return
	}

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
	return
}

func prepareMMS() (result int) {

	a := fetch.StorageMMSData

	if len(a) == 0 {
		lib.LogParseErr(2, "MMS нет")
		result++
		return
	}

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

	return
}

func printFastAndSlow(x []fetch.EmailData) (result map[string][][]fetch.EmailData) {

	v := map[string][][]fetch.EmailData{}
	alpha := x[0].Country
	start := 0

	for i, data := range x {

		if alpha != data.Country {
			alpha = data.Country

			fastProviders := x[start : start+3]
			slowProviders := x[i-3 : i]

			fmt.Println(fastProviders)

			v[data.Country] = append(v[data.Country], fastProviders)
			v[data.Country] = append(v[data.Country], slowProviders)

			//fmt.Println(v)

			start = i

		}

	}

	result = v
	return

}

func prepareEmail() (result int) {

	a := fetch.StorageEmail
	if len(a) == 0 {
		lib.LogParseErr(2, "MMS нет")
		result++
		return
	}

	x := sortByCountryAndSpeed(a)
	fetch.ResultSet.Email = printFastAndSlow(x)

	return
}

func prepareAccident() (result int) {

	if len(fetch.StorageAccidentData) == 0 {
		lib.LogParseErr(2, "Инцидентов нет")
		result++
		return
	}
	x := sortByAccident(fetch.StorageAccidentData)

	fetch.ResultSet.Incidents = x

	//fmt.Println(fetch.ResultSet)
	return
}

func sumActiveTicket(a []fetch.SupportData) (ticketCount int) {
	for _, data := range a {
		ticketCount += data.ActiveTickets
	}
	return
}

func prepareSupport() (result int) {

	a := fetch.StorageSupportData

	if len(fetch.StorageAccidentData) == 0 {
		lib.LogParseErr(2, "Тикетов нет")
		result++
		return
	}

	activeTicket := sumActiveTicket(a)
	supportLoad := 0

	switch {

	case activeTicket < lowLoad:
		{
			supportLoad = 1
		}
	case activeTicket < middleLoad:
		{
			supportLoad = 2
		}
	default:
		{
			supportLoad = 3
		}

	}

	//fmt.Println(activeTicket)
	//fmt.Println(supportLoad)

	fetch.ResultSet.Support = append(fetch.ResultSet.Support, supportLoad)

	timeToResolveTicket := 60 / averageSupportBandwith //* specCount

	fetch.ResultSet.Support = append(fetch.ResultSet.Support,
		supportLoad*timeToResolveTicket)

	return
}

func GetResultData() { //11.1

	errCount := 0

	errCount += prepareSMS()                               //11.2
	errCount += prepareMMS()                               //11.3
	fetch.ResultSet.VoiceCall = fetch.StorageVoiceCallData //11.4
	errCount += prepareEmail()                             //11.5
	fetch.ResultSet.Billing = fetch.StorageBilling         //11.6
	errCount += prepareSupport()                           //11.7
	errCount += prepareAccident()                          //11.8

	if errCount > 0 {
		fetch.Result.Error = "Ошибок сбора ResultSet:" + strconv.Itoa(errCount)
		fetch.Result.Status = true
	} else {
		fetch.Result.Error = ""
		fetch.Result.Status = false
	}

	fetch.Result.Data = fetch.ResultSet

	lib.LogParseErr(2, fetch.Result.Error)

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
