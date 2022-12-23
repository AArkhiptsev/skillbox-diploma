package main

import (
	"diploma/fetch"
	"diploma/lib"
	"diploma/serve"
	"fmt"
	"sort"
	"time"
)

const (
	smsFileName       = "../emul/sms.data"
	voiceCallFileName = "../emul/voice.data"
	emailFileName     = "../emul/email.data"
	billingFileName   = "../emul/billing.data"
	mmsDataServer     = "http://127.0.0.1:8383/mms"
	supportServer     = "http://127.0.0.1:8383/support"
	accidentServer    = "http://127.0.0.1:8383/accendent"
)

func init() {
	sort.Strings(fetch.SmsProviders) //отсортируем провайдеров, чтобы ускорить поиск по ним
	sort.Strings(fetch.MmsProviders)
	sort.Strings(fetch.VoiceCallProviders)
	sort.Strings(fetch.EmailProviders)
	sort.Strings(fetch.AccidentStatus)

}

func logSortProviders() {
	lib.LogParseErr(0, "Отсортированные массивы провайдеров:")
	lib.LogParseErr(0, fmt.Sprintf("SMS: %v", fetch.SmsProviders))
	lib.LogParseErr(0, fmt.Sprintf("Voice Calls: %v", fetch.VoiceCallProviders))
	lib.LogParseErr(0, fmt.Sprintf("Email: %v", fetch.EmailProviders))
	lib.LogParseErr(0, fmt.Sprintf("Accident Status: %v", fetch.AccidentStatus))
}

func smsHandler() {

	lib.LogParseErr(1, "Начат разбор файла SMS")
	line, errCount := fetch.ParseSMS(smsFileName)
	lib.StdParseMessage(smsFileName, line, errCount)

	fetch.LogStorageHeaderData(fetch.StorageSMSData)

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", smsFileName))

}

func mmsHandler() {

	lib.LogParseErr(1, "Запросим данные об MMS "+mmsDataServer)
	fetch.ParseMMS(mmsDataServer)

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageHeaderData(fetch.StorageMMSData)

	lib.LogParseErr(1, "Обработка MMS завершена")

}

func supportHandler() {
	lib.LogParseErr(1, "Запросим данные об MMS "+supportServer)
	fetch.ParseSupport(supportServer)

	lib.LogParseErr(0, "Результат:")
	fetch.LogSupportData()
	lib.LogParseErr(1, "Обработка MMS завершена")

}

func voiceCallHandler() {

	lib.LogParseErr(1, "Начат разбор файла Voice Calls")
	line, errCount := fetch.ParseVoicesCall(voiceCallFileName)
	lib.StdParseMessage(voiceCallFileName, line, errCount)

	fetch.LogStorageHeaderData(fetch.StorageSMSData)

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageVoicesCallsData()

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", voiceCallFileName))
}

func emailHandler() {

	lib.LogParseErr(1, "Начат разбор файла Email")
	line, errCount := fetch.ParseEmail(emailFileName)
	lib.StdParseMessage(emailFileName, line, errCount)

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageEmailData()

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", emailFileName))
}

func billingHandler() {
	lib.LogParseErr(1, "Начат разбор файла Billing")
	line, errCount := fetch.ParseBilling(billingFileName)
	lib.StdParseMessage(billingFileName, line, errCount)

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageBilling()

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", billingFileName))
}

func accidentHandler() {

	lib.LogParseErr(1, "Запросим данные об инцидентах "+accidentServer)
	fetch.ParseAccident(accidentServer)

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageAccidentData()

	lib.LogParseErr(1, "Обработка инцидентов завершена")

}

func main() {

	lib.LogParseErr(0, "Старт...")

	logSortProviders()

	smsHandler()
	//voiceCallHandler()
	//emailHandler()
	//billingHandler()

	//supportHandler()
	//mmsHandler()
	//accidentHandler()

	lib.LogParseErr(0, "Сбор данных завершен.")

	go lib.Spinner(80 * time.Millisecond)
	serve.ListenAndServeHTTP()

}
