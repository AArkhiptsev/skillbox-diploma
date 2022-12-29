package main

import (
	"diploma/conf"
	"diploma/fetch"
	"diploma/lib"
	"diploma/serve"
	"fmt"
	"sort"
	"time"
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
	line, errCount := fetch.ParseSMS(conf.Config.Files.SmsFileName)
	lib.StdParseMessage(conf.Config.Files.SmsFileName, line, errCount)

	fetch.LogStorageHeaderData(fetch.StorageSMSData)

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", conf.Config.Files.SmsFileName))

}

func mmsHandler() {

	lib.LogParseErr(1, "Запросим данные об MMS "+conf.Config.FetchServers.MmsDataServer)
	line, errCount := fetch.ParseMMS(conf.Config.FetchServers.MmsDataServer)
	lib.StdParseMessage(conf.Config.FetchServers.MmsDataServer, line, errCount)

	fetch.LogStorageHeaderData(fetch.StorageMMSData)

	lib.LogParseErr(1, "Обработка MMS завершена")

}

func supportHandler() {
	lib.LogParseErr(1, "Запросим данные о поддержке "+conf.Config.FetchServers.SupportServer)
	line, errCount := fetch.ParseSupport(conf.Config.FetchServers.SupportServer)
	lib.StdParseMessage(conf.Config.FetchServers.SupportServer, line, errCount)

	fetch.LogSupportData()
	lib.LogParseErr(1, "Обработка данных о поддержке завершена")

}

func voiceCallHandler() {

	lib.LogParseErr(1, "Начат разбор файла Voice Calls")
	line, errCount := fetch.ParseVoicesCall(conf.Config.Files.VoiceCallFileName)
	lib.StdParseMessage(conf.Config.Files.VoiceCallFileName, line, errCount)

	fetch.LogStorageHeaderData(fetch.StorageSMSData)

	fetch.LogStorageVoicesCallsData()

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", conf.Config.Files.VoiceCallFileName))
}

func emailHandler() {

	lib.LogParseErr(1, "Начат разбор файла Email")
	line, errCount := fetch.ParseEmail(conf.Config.Files.EmailFileName)
	lib.StdParseMessage(conf.Config.Files.EmailFileName, line, errCount)

	fetch.LogStorageEmailData()

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", conf.Config.Files.EmailFileName))
}

func billingHandler() {
	lib.LogParseErr(1, "Начат разбор файла Billing")
	line, errCount := fetch.ParseBilling(conf.Config.Files.BillingFileName)
	lib.StdParseMessage(conf.Config.Files.BillingFileName, line, errCount)

	fetch.LogStorageBilling()

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", conf.Config.Files.BillingFileName))
}

func accidentHandler() {

	lib.LogParseErr(1, "Запросим данные об инцидентах "+conf.Config.FetchServers.AccidentServer)
	line, errCount := fetch.ParseAccident(conf.Config.FetchServers.AccidentServer)
	lib.StdParseMessage(conf.Config.FetchServers.AccidentServer, line, errCount)

	fetch.LogStorageAccidentData()

	lib.LogParseErr(1, "Обработка инцидентов завершена")

}

func main() {

	lib.LogParseErr(0, "Старт...")

	conf.ReadConfig(conf.ConfigFileName)

	fmt.Println(conf.Config.FetchServers.AccidentServer)

	logSortProviders()

	smsHandler()
	voiceCallHandler()
	emailHandler()
	billingHandler()

	supportHandler()
	mmsHandler()
	accidentHandler()

	lib.LogParseErr(0, "Сбор всех данных завершен.")

	lib.LogParseErr(1, "Формирование результата")
	serve.GetResultData()

	go lib.Spinner(80 * time.Millisecond)
	serve.ListenAndServeHTTP()

}
