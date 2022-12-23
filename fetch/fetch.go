package fetch

import (
	"bufio"
	"diploma/lib"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	csvSeparator = ";"
)

var (
	SmsProviders       = []string{"Topolo", "Rond", "Kildy"}
	VoiceCallProviders = []string{"TransparentCalls", "E-Voice", "JustPhone"}
	EmailProviders     = []string{"Gmail", "Yahoo", "Hotmail", "MSN", "Orange",
		"Comcast", "AOL", "Live", "RediffMail", "GMX", "Protonmail", "Yandex", "Mail.ru"}
	AccidentStatus = []string{"closed", "active"}
	MmsProviders   = SmsProviders //работаем с копией, сделано на случай,
	// если в будущем появится иной набор провайдеров для MMS
)

type headerData struct {
	Country      string `json:"country"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
	Provider     string `json:"provider"`
}

type voiceCallData struct {
	header              headerData
	connectionStability float32
	tTFB                int
	voicePurity         int
	medianOfCallsTime   int
}

type emailData struct {
	Country      string
	Provider     string
	DeliveryTime int
}

type BillingData struct {
	CreateCustomer bool
	Purchase       bool
	Payout         bool
	Recurring      bool
	FraudControl   bool
	CheckoutPage   bool
}

type supportData struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}

type AccidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"` // возможные статусы active и	closed
}

type ResultT struct {
	Status bool `json:"status"` // true, если все этапы сбора данных прошли успешно,
	// false во всех остальных случаях
	Data ResultSetT `json:"data"` // заполнен, если все этапы сбора данных прошли успешно,
	// nil во всех остальных случаях
	Error string `json:"error"` // пустая строка если все этапы сбора данных прошли успешно, в случае ошибки заполнено текстом ошибки
	//(детали ниже)
}

type ResultSetT struct {
	SMS       [][]headerData           `json:"sms"`
	MMS       [][]headerData           `json:"mms"`
	VoiceCall []voiceCallData          `json:"voice_call"`
	Email     map[string][][]emailData `json:"email"`
	Billing   BillingData              `json:"billing"`
	Support   []int                    `json:"support"`
	Incidents []AccidentData           `json:"incident"`
}

var (
	StorageSMSData       = []headerData{}
	StorageMMSData       = []headerData{}
	storageVoiceCallData = []voiceCallData{}
	storageEmail         = []emailData{}
	storageBilling       = BillingData{}
	storageSupportData   = []supportData{}
	storageAccidentData  = []AccidentData{}
)

func (s headerData) check(providers []string, lineNumber int) (result bool) {

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

func (s *emailData) check(providers []string, deliveryTime string, lineNumber int) (result bool) {
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

func ParseSMS(filename string) (lineCounter, parseErrCount int) {

	lib.LogParseErr(1, "Открытие файла: "+filename)

	file, err := os.Open(filename)
	if err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		splitString := strings.Split(scanner.Text(), csvSeparator)
		//log.Println("Парсинг:", splitString)

		if len(splitString) != 4 {
			lib.LogParseErr(4,
				fmt.Sprintf("Ошибка количества элементов. Строка: %v", lineCounter))
			parseErrCount++
			continue
		}

		s := headerData{
			Country:      splitString[0],
			Bandwidth:    splitString[1],
			ResponseTime: splitString[2],
			Provider:     splitString[3],
		}

		if !(s.check(SmsProviders, lineCounter)) {
			parseErrCount++
			continue
		}
		StorageSMSData = append(StorageSMSData, s)

		lineCounter++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return

}

func ParseVoicesCall(filename string) (lineCounter, parseErrCount int) {

	lib.LogParseErr(1, "Открытие файла: "+filename)

	file, err := os.Open(filename)
	if err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		splitString := strings.Split(scanner.Text(), csvSeparator)

		if len(splitString) != 8 {
			lib.LogParseErr(4,
				fmt.Sprintf("Ошибка количества элементов. Строка: %v", lineCounter))
			parseErrCount++
			continue
		}

		s := voiceCallData{
			header: headerData{
				Country:      splitString[0],
				Bandwidth:    splitString[1],
				ResponseTime: splitString[2],
				Provider:     splitString[3],
			},
		}

		if !(s.header.check(VoiceCallProviders, lineCounter)) {
			parseErrCount++
			continue
		}

		if !(s.check(splitString[4:8], lineCounter)) {
			parseErrCount++
			continue
		}

		storageVoiceCallData = append(storageVoiceCallData, s)

		lineCounter++

	}

	lib.LogParseErr(0,
		fmt.Sprintf("Обработано строк без ошибок: %v", lineCounter))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return

}

func ParseEmail(filename string) (lineCounter, parseErrCount int) {

	lib.LogParseErr(1, "Открытие файла: "+filename)

	file, err := os.Open(filename)
	if err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		splitString := strings.Split(scanner.Text(), csvSeparator)

		//log.Println(splitString)

		if len(splitString) != 3 {
			lib.LogParseErr(4,
				fmt.Sprintf("Ошибка количества элементов. Строка: %v", lineCounter))
			parseErrCount++
			continue
		}

		s := emailData{
			Country:  splitString[0],
			Provider: splitString[1],
		}

		if !(s.check(EmailProviders, splitString[2], lineCounter)) {
			parseErrCount++
			continue
		}

		storageEmail = append(storageEmail, s)
		lineCounter++

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return

}

func ParseBilling(filename string) (lineCounter, parseErrCount int) {

	lib.LogParseErr(1, "Открытие файла: "+filename)

	file, err := os.Open(filename)
	if err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		//	fmt.Println(scanner.Text())

		i, err := strconv.ParseInt(scanner.Text(), 2, 64)

		if err != nil {
			lib.LogParseErr(4, "Ошибка конвертации строки. "+err.Error())
			parseErrCount++
			return
		}

		lib.LogParseErr(0,
			fmt.Sprintf("Значение в dec- формате: %d, в bin- формате:  %b", i, i))

		storageBilling.parse(i)

		lineCounter++
	}

	return
}

func ParseMMS(URL string) {

	content, err := lib.RequestContent(URL)

	if err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	if err := json.Unmarshal(content, &StorageMMSData); err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	lib.LogParseErr(1,
		fmt.Sprintf("Разбор JSON произведен. Записей %v", len(StorageMMSData)))

	lib.LogParseErr(1, "Проверка на корректность значений...")

	k := len(StorageMMSData)
	errCount := 0

	for i := 0; i < k; i++ {

		if !(headerData.check(StorageMMSData[i], MmsProviders, i)) {
			StorageMMSData = append(StorageMMSData[:i],
				StorageMMSData[i+1:]...)
			k--
			errCount++
		}

	}

	lib.LogParseErr(1,
		fmt.Sprintf("Проверка корректности произведена. Записей %v", len(StorageMMSData)))

}

func ParseAccident(URL string) {

	content, err := lib.RequestContent(URL)

	if err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	if err := json.Unmarshal(content, &storageAccidentData); err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	lib.LogParseErr(1,
		fmt.Sprintf("Разбор JSON произведен. Записей %v", len(storageAccidentData)))

	lib.LogParseErr(1, "Проверка на корректность значений...")

	k := len(storageAccidentData)
	errCount := 0

	for i := 0; i < k; i++ {

		if !(AccidentData.check(storageAccidentData[i], AccidentStatus, i)) {
			storageAccidentData = append(storageAccidentData[:i],
				storageAccidentData[i+1:]...)
			k--
			errCount++
		}

	}

	lib.LogParseErr(1,
		fmt.Sprintf("Проверка корректности произведена. Записей %v", len(storageAccidentData)))

	return
}

func ParseSupport(URL string) {

	content, err := lib.RequestContent(URL)

	if err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	if err := json.Unmarshal(content, &storageSupportData); err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	lib.LogParseErr(1,
		fmt.Sprintf("Разбор JSON произведен. Записей %v", len(storageSupportData)))

	lib.LogParseErr(1, "Проверка на корректность значений...")

}
