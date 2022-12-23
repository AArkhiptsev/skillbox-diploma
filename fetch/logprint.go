package fetch

import "log"

func LogStorageHeaderData(v []headerData) {
	for _, datum := range v {
		log.Println(datum)
	}
}

func LogStorageVoicesCallsData() {
	for _, datum := range storageVoiceCallData {
		log.Println(datum)
	}
}

func LogStorageEmailData() {
	for _, datum := range storageEmail {
		log.Println(datum)
	}
}

func LogStorageBilling() {
	log.Println("CreateCustomer :", storageBilling.CreateCustomer)
	log.Println("Purchase       :", storageBilling.Purchase)
	log.Println("Payout         :", storageBilling.Payout)
	log.Println("Recurring      :", storageBilling.Recurring)
	log.Println("FraudControl   :", storageBilling.FraudControl)
	log.Println("CheckoutPage   :", storageBilling.CheckoutPage)
}

func LogSupportData() {
	for _, datum := range storageSupportData {
		log.Println(datum)
	}
}

func LogStorageAccidentData() {
	for _, datum := range storageAccidentData {
		log.Println(datum)
	}
}
