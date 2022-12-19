package main

import (
	sms "diploma/lib"
	"log"
	"sort"
)

const smsFileName = "../emul/sms.data"

func main() {

	log.Println("Старт...")
	log.Println("Подготовим массив провайдеров.")
	sort.Strings(sms.SmsProviders)
	log.Println(sms.SmsProviders)
	sms.FetchSMS(smsFileName)

}
