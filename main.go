package main

import (
	"diploma/fetch"
	"diploma/lib"
	"fmt"
	"sort"
)

const (
	smsFileName       = "../emul/sms.data"
	voiceCallFileName = "../emul/voice.data"
	emailFileName     = "../emul/email.data"
	mmsDataServer     = "http://127.0.0.1:8383/mms"
)

func init() {
	sort.Strings(fetch.SmsProviders) //отсортируем провайдеров, чтобы ускорить поиск по ним
	sort.Strings(fetch.MmsProviders)
	sort.Strings(fetch.VoiceCallProviders)
	sort.Strings(fetch.EmailProviders)

}

func logSortProviders() {
	lib.LogParseErr(0, "Отсортированные массивы провайдеров:")
	lib.LogParseErr(0, fmt.Sprintf("SMS: %v", fetch.SmsProviders))
	lib.LogParseErr(0, fmt.Sprintf("Voice Calls: %v", fetch.VoiceCallProviders))
	lib.LogParseErr(0, fmt.Sprintf("Email: %v", fetch.EmailProviders))
}

func smsHandler() {

	lib.LogParseErr(1, "Начат разбор файла SMS")
	lib.LogParseErr(2,
		fmt.Sprintf("Разобран файл %v, ошибок разбора %v", smsFileName,
			fetch.FetchSMS(smsFileName)))

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageHeaderData()

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", smsFileName))

}

func mmsHandler() {

	lib.LogParseErr(1, "Запросим данные об MMS "+mmsDataServer)
	fetch.FetchMMS(mmsDataServer)

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageHeaderData()

	lib.LogParseErr(1, "Обработка MMS завершена")

}

func voiceCallHandler() {

	lib.LogParseErr(1, "Начат разбор файла Voice Calls")
	lib.LogParseErr(2,
		fmt.Sprintf("Разобран файл %v, ошибок разбора %v", voiceCallFileName,
			fetch.FetchVoicesCall(voiceCallFileName)))

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageVoicesCallsData()

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", voiceCallFileName))
}

func emailHandler() {

	lib.LogParseErr(1, "Начат разбор файла Email")

	fetch.FetchEmail(emailFileName)

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", emailFileName))
}

func main() {

	lib.LogParseErr(0, "Старт...")

	logSortProviders()

	//smsHandler()
	voiceCallHandler()
	//emailHandler()

	//mmsHandler()

	lib.LogParseErr(0, "Завершение...")
}
