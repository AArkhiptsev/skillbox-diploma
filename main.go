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
	mmsDataServer     = "http://127.0.0.1:8383/mms"
)

func init() {
	sort.Strings(fetch.SmsProviders) //отсортируем провайдеров, чтобы ускорить поиск по ним
	sort.Strings(fetch.MmsProviders)
	sort.Strings(fetch.VoiceCallProviders)

}

func smsHandler() {

	lib.LogParseErr(1, "Начат разбр файла SMS")
	lib.LogParseErr(2,
		fmt.Sprintf("Разобран файл %v, ошибок разбора %v", smsFileName,
			fetch.FetchSMS(smsFileName)))

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageSMSData()

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", smsFileName))

}

func mmsHandler() {

	lib.LogParseErr(1, "Запросим данные об MMS "+mmsDataServer)
	fetch.FetchMMS(mmsDataServer)

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageMMSData()

	lib.LogParseErr(1, "Обработка MMS завершена")

}

func voiceCallHandler() {

	lib.LogParseErr(1, "Начат разбр файла Voice Calls")
	lib.LogParseErr(2,
		fmt.Sprintf("Разобран файл %v, ошибок разбора %v", voiceCallFileName,
			fetch.FetchVoicesCall(voiceCallFileName)))

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageVoicesCallsData()

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", voiceCallFileName))
}

func main() {

	lib.LogParseErr(0, "Старт...")

	lib.LogParseErr(0, "Подготовленный массив провайдеров:")
	lib.LogParseErr(0, fmt.Sprintf("SMS: %v\n", fetch.SmsProviders))
	lib.LogParseErr(0, fmt.Sprintf("Voice Calls: %v\n", fetch.VoiceCallProviders))

	smsHandler()
	mmsHandler()
	voiceCallHandler()

}
