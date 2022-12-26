package fetch

import (
	"diploma/lib"
	"fmt"
	"strconv"
)

func (s HeaderData) check(providers []string, lineNumber int) (result bool) {

	result = false

	if lib.GetCountryNameByAlpha(s.Country) == "" {
		lib.LogParseErr(3,
			fmt.Sprintf(" alpha: %v, строка: %v", s.Country, lineNumber))
		return
	}

	if !(lib.Found(s.Provider, providers)) {
		lib.LogParseErr(3,
			fmt.Sprintf(" провайдер: %v, строка: %v", s.Provider, lineNumber))
		return
	}

	if _, err := strconv.Atoi(s.ResponseTime); err != nil {
		lib.LogParseErr(3,
			fmt.Sprintf(" среднее время ответа: %v, строка: %v",
				s.ResponseTime, lineNumber))
		return
	}

	if _, err := strconv.Atoi(s.Bandwidth); err != nil {
		lib.LogParseErr(3,
			fmt.Sprintf(" полоса пропускания: %v, строка: %v",
				s.Bandwidth, lineNumber))
		return
	}

	result = true

	return

}

func (s *voiceCallData) check(val []string, lineNumber int) (result bool) {

	result = false

	b, err := strconv.ParseFloat(val[0], 32)
	if err != nil {
		lib.LogParseErr(3,
			fmt.Sprintf(" connectionStability: %v. строка %v",
				val[0], lineNumber))
		return
	}
	s.connectionStability = float32(b)

	s.tTFB, err = strconv.Atoi(val[1])
	if err != nil {
		lib.LogParseErr(3,
			fmt.Sprintf(" tTFB: %v. строка %v",
				val[1], lineNumber))
		return
	}

	s.voicePurity, err = strconv.Atoi(val[2])
	if err != nil {
		lib.LogParseErr(3,
			fmt.Sprintf(" voicePurity:  %v. строка %v",
				val[2], lineNumber))
		return
	}

	s.medianOfCallsTime, err = strconv.Atoi(val[3])
	if err != nil {
		lib.LogParseErr(3,
			fmt.Sprintf(" medianOfCallsTime: %v. строка %v",
				val[3], lineNumber))
		return
	}

	result = true

	return
}

func (s *EmailData) check(providers []string, deliveryTime string, lineNumber int) (result bool) {
	result = false

	if lib.GetCountryNameByAlpha(s.Country) == "" {
		lib.LogParseErr(3,
			fmt.Sprintf(" alpha: %v, строка: %v", s.Country, lineNumber))
		return
	}

	if !(lib.Found(s.Provider, providers)) {
		lib.LogParseErr(3,
			fmt.Sprintf(" провайдер: %v, строка: %v", s.Provider, lineNumber))
		return
	}

	b, err := strconv.Atoi(deliveryTime)
	if err != nil {
		lib.LogParseErr(3,
			fmt.Sprintf(" среднее время ответа: %v, строка: %v",
				deliveryTime, lineNumber))
		return
	}
	s.DeliveryTime = b

	result = true
	return

}

func (s AccidentData) check(statuses []string, lineNumber int) (result bool) {

	result = false

	if !(lib.Found(s.Status, statuses)) {
		lib.LogParseErr(3,
			fmt.Sprintf(" статус: %v, строка: %v", s.Status, lineNumber))
		return
	}

	result = true

	return

}

func (b *BillingData) parse(a int64) {

	bits := []byte(strconv.FormatInt(int64(a), 2))

	b.CreateCustomer = lib.CheckBit(bits[0])
	b.Purchase = lib.CheckBit(bits[1])
	b.Payout = lib.CheckBit(bits[2])
	b.Recurring = lib.CheckBit(bits[3])
	b.FraudControl = lib.CheckBit(bits[4])
	b.CheckoutPage = lib.CheckBit(bits[5])

	return
}
