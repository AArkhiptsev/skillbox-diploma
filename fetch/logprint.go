package fetch

import "log"

func LogStorageHeaderData(v []HeaderData) {
	for _, datum := range v {
		log.Println(datum)
	}
}

func LogStorageVoicesCallsData() {
	for _, datum := range StorageVoiceCallData {
		log.Println(datum)
	}
}

func LogStorageEmailData() {
	for _, datum := range StorageEmail {
		log.Println(datum)
	}
}

func LogStorageBilling() {
	log.Println("CreateCustomer :", StorageBilling.CreateCustomer)
	log.Println("Purchase       :", StorageBilling.Purchase)
	log.Println("Payout         :", StorageBilling.Payout)
	log.Println("Recurring      :", StorageBilling.Recurring)
	log.Println("FraudControl   :", StorageBilling.FraudControl)
	log.Println("CheckoutPage   :", StorageBilling.CheckoutPage)
}

func LogSupportData() {
	for _, datum := range StorageSupportData {
		log.Println(datum)
	}
}

func LogStorageAccidentData() {
	for _, datum := range StorageAccidentData {
		log.Println(datum)
	}
}
