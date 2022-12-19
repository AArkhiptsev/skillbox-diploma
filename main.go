package main

import (
	sms "diploma/lib"
	"log"
	"sort"
)

const smsFileName = "../emul/sms.data"

func init() {
	sort.Strings(sms.SmsProviders) //отсортируем провайдеров, чтобы ускорить поиск по ним
}

func main() {

	log.Println("Старт...")

	log.Println("Подготовленный массив провайдеров:", sms.SmsProviders)
	log.Printf("Разобран файл %v,"+sms.ColorRed+" ошибок разбора %v"+
		sms.ColorReset, smsFileName, sms.FetchSMS(smsFileName))
	log.Println("Результат:")
	sms.LogStorageSMSData()
	log.Printf(sms.ColorCyan+"Обработка %v завершена"+sms.ColorReset, smsFileName)

}
