package main

import (
	"diploma/fetch"
	"diploma/lib"
	"fmt"
	"sort"
)

const (
	smsFileName   = "../emul/sms.data"
	mmsDataServer = "http://127.0.0.1:8383/mms"
)

func init() {
	sort.Strings(fetch.SmsProviders) //отсортируем провайдеров, чтобы ускорить поиск по ним
}

func smsHandler() {
	lib.LogParseErr(1,
		fmt.Sprintf("Подготовленный массив SMS-провайдеров: %v",
			fetch.SmsProviders))

	lib.LogParseErr(2,
		fmt.Sprintf("Разобран файл %v, ошибок разбора %v", smsFileName,
			fetch.FetchSMS(smsFileName)))

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageSMSData()

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", smsFileName))

}

func mmsHandler() {

	lib.LogParseErr(1, "Запросим данные по MMS "+mmsDataServer)
	fetch.FetchMMS(mmsDataServer)

}

func main() {

	lib.LogParseErr(0, "Старт...")

	smsHandler()
	mmsHandler()

}
