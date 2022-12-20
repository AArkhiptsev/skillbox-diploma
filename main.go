package main

import (
	"diploma/fetch"
	"diploma/lib"
	"fmt"
	"sort"
)

const smsFileName = "../emul/sms.data"

func init() {
	sort.Strings(fetch.SmsProviders) //отсортируем провайдеров, чтобы ускорить поиск по ним
}

func main() {

	lib.LogParseErr(0, "Старт...")

	lib.LogParseErr(1,
		fmt.Sprint("Подготовленный массив провайдеров: %s", fetch.SmsProviders))

	lib.LogParseErr(2,
		fmt.Sprintf("Разобран файл %v, ошибок разбора %v", smsFileName,
			fetch.FetchSMS(smsFileName)))

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageSMSData()

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", smsFileName))

}
