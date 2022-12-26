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

type (
	HeaderData struct {
		Country      string `json:"country"`
		Bandwidth    string `json:"bandwidth"`
		ResponseTime string `json:"response_time"`
		Provider     string `json:"provider"`
	}

	voiceCallData struct {
		header              HeaderData
		connectionStability float32
		tTFB                int
		voicePurity         int
		medianOfCallsTime   int
	}

	EmailData struct {
		Country      string
		Provider     string
		DeliveryTime int
	}

	BillingData struct {
		CreateCustomer bool
		Purchase       bool
		Payout         bool
		Recurring      bool
		FraudControl   bool
		CheckoutPage   bool
	}

	supportData struct {
		Topic         string `json:"topic"`
		ActiveTickets int    `json:"active_tickets"`
	}

	AccidentData struct {
		Topic  string `json:"topic"`
		Status string `json:"status"` // возможные статусы active и	closed
	}
)

type ResultT struct {
	Status bool `json:"status"` // true, если все этапы сбора данных прошли успешно,
	// false во всех остальных случаях
	Data ResultSetT `json:"data"` // заполнен, если все этапы сбора данных прошли успешно,
	// nil во всех остальных случаях
	Error string `json:"error"` // пустая строка если все этапы сбора данных прошли успешно, в случае ошибки заполнено текстом ошибки
}

type ResultSetT struct {
	SMS       [][]HeaderData           `json:"sms"`
	MMS       [][]HeaderData           `json:"mms"`
	VoiceCall []voiceCallData          `json:"voice_call"`
	Email     map[string][][]EmailData `json:"email"`
	Billing   BillingData              `json:"billing"`
	Support   []int                    `json:"support"`
	Incidents []AccidentData           `json:"incident"`
}

var (
	StorageSMSData       = []HeaderData{}
	StorageMMSData       = []HeaderData{}
	StorageVoiceCallData = []voiceCallData{}
	StorageEmail         = []EmailData{}
	StorageBilling       = BillingData{}
	storageSupportData   = []supportData{}
	StorageAccidentData  = []AccidentData{}
	ResultSet            = ResultSetT{}
	ResultEmail          = map[string][][]EmailData{}
)

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

		s := HeaderData{
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
			header: HeaderData{
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

		StorageVoiceCallData = append(StorageVoiceCallData, s)

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

		s := EmailData{
			Country:  splitString[0],
			Provider: splitString[1],
		}

		if !(s.check(EmailProviders, splitString[2], lineCounter)) {
			parseErrCount++
			continue
		}

		StorageEmail = append(StorageEmail, s)
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

		StorageBilling.parse(i)

		lineCounter++
	}

	return
}

func ParseMMS(URL string) (lineCounter, parseErrCount int) {

	content, err := lib.RequestContent(URL)

	if err != nil {
		lib.LogParseErr(4, err.Error())
		parseErrCount++
		return
	}

	if err := json.Unmarshal(content, &StorageMMSData); err != nil {
		lib.LogParseErr(4, err.Error())
		parseErrCount++
		return
	}

	k := len(StorageMMSData)

	lineCounter = k

	for i := 0; i < k; i++ {

		if !(HeaderData.check(StorageMMSData[i], MmsProviders, i)) {
			StorageMMSData = append(StorageMMSData[:i],
				StorageMMSData[i+1:]...)
			k--
			parseErrCount++
		}

	}
	return
}

func ParseAccident(URL string) (lineCounter, parseErrCount int) {

	content, err := lib.RequestContent(URL)

	if err != nil {
		lib.LogParseErr(4, err.Error())
		parseErrCount++
		return
	}

	if err := json.Unmarshal(content, &StorageAccidentData); err != nil {
		lib.LogParseErr(4, err.Error())
		parseErrCount++
		return
	}

	k := len(StorageAccidentData)
	lineCounter = k

	for i := 0; i < k; i++ {

		if !(AccidentData.check(StorageAccidentData[i], AccidentStatus, i)) {
			StorageAccidentData = append(StorageAccidentData[:i],
				StorageAccidentData[i+1:]...)
			k--
			parseErrCount++
		}

	}

	return
}

func ParseSupport(URL string) (lineCounter, parseErrCount int) {

	content, err := lib.RequestContent(URL)

	if err != nil {
		lib.LogParseErr(4, err.Error())
		parseErrCount++
		return
	}

	if err := json.Unmarshal(content, &storageSupportData); err != nil {
		lib.LogParseErr(4, err.Error())
		parseErrCount++
		return
	}

	lineCounter = len(storageSupportData)

	return
}
