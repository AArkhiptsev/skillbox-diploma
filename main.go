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
	billingFileName   = "../emul/billing.data"
	mmsDataServer     = "http://127.0.0.1:8383/mms"
	supportServer     = "http://127.0.0.1:8383/support"
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
	line, errCount := fetch.ParseSMS(smsFileName)
	lib.LogParseErr(1, fmt.Sprintf("Разобран файл %v", smsFileName))
	lib.LogParseErr(2,
		fmt.Sprintf("Разобрано строк: %v, ошибок: %v", line, errCount))

	lib.LogParseErr(0, "Результат:")
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
	lib.LogParseErr(2,
		fmt.Sprintf("Разобран файл %v, ошибок разбора %v", voiceCallFileName,
			fetch.PatchVoicesCall(voiceCallFileName)))

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageVoicesCallsData()

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", voiceCallFileName))
}

func emailHandler() {

	lib.LogParseErr(1, "Начат разбор файла Email")
	lib.LogParseErr(2,
		fmt.Sprintf("Разобран файл %v, ошибок разбора %v", voiceCallFileName,
			fetch.ParseEmail(emailFileName)))

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageEmailData()

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", emailFileName))
}

func billingHandler() {
	lib.LogParseErr(1, "Начат разбор файла Billing")
	lib.LogParseErr(2,
		fmt.Sprintf("Разобран файл %v, ошибок разбора %v", billingFileName,
			fetch.ParseBilling(billingFileName)))

	lib.LogParseErr(0, "Результат:")
	fetch.LogStorageBilling()

	lib.LogParseErr(1,
		fmt.Sprintf("Обработка %v завершена", billingFileName))
}

func main() {

	lib.LogParseErr(0, "Старт...")

	logSortProviders()

	smsHandler()
	voiceCallHandler()
	emailHandler()
	//billingHandler()

	supportHandler()
	mmsHandler()

	lib.LogParseErr(0, "Завершение...")
}
