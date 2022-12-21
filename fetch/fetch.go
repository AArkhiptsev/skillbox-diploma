package fetch

import (
	"bufio"
	"diploma/lib"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	MmsProviders = SmsProviders //работаем с копией, сделано на случай,
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

var (
	storageSMSData       = []headerData{}
	storageMMSData       = []headerData{}
	storageVoiceCallData = []voiceCallData{}
	storageEmail         = []emailData{}
)

func removeIndex(s []headerData, index int) []headerData {
	return append(s[:index], s[index+1:]...)
}

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

func LogStorageHeaderData() {
	for _, datum := range storageSMSData {
		log.Println(datum)
	}
}

func LogStorageVoicesCallsData() {
	for _, datum := range storageVoiceCallData {
		log.Println(datum)
	}
}

func FetchSMS(filename string) (parseErrCount int) {

	lineCounter := 0

	lib.LogParseErr(1, "Открытие файла: "+filename)

	file, err := os.Open(filename)
	if err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		splittedString := strings.Split(scanner.Text(), csvSeparator)
		//log.Println("Парсинг:", splittedString)

		if len(splittedString) == 4 {

			s := headerData{
				Country:      splittedString[0],
				Bandwidth:    splittedString[1],
				ResponseTime: splittedString[2],
				Provider:     splittedString[3],
			}

			if s.check(SmsProviders, lineCounter) {
				storageSMSData = append(storageSMSData, s)
			} else {
				parseErrCount++
			}

		} else {
			lib.LogParseErr(4,
				fmt.Sprintf("Ошибка количества элементов. Строка: %v", lineCounter))
			parseErrCount++
		}

		lineCounter++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	lib.LogParseErr(0,
		fmt.Sprintf("Обработано строк: %v", lineCounter))

	return

}

func FetchMMS(URL string) {

	resp, err := http.Get(URL)

	if err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		lib.LogParseErr(4,
			fmt.Sprintf("Код ответа сервера: %v", resp.StatusCode))
		return
	}

	lib.LogParseErr(0,
		fmt.Sprintf("Код ответа сервера: %v", resp.StatusCode))

	lib.LogParseErr(1, "Произведем JSON разбор...")
	content, err := io.ReadAll(resp.Body)

	if err != nil {
		lib.LogParseErr(0,
			fmt.Sprintf("Ошибка чтения Body: %v", err.Error()))
		return
	}

	if err := json.Unmarshal(content, &storageMMSData); err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	lib.LogParseErr(1,
		fmt.Sprintf("Разбор JSON произведен. Записей %v", len(storageMMSData)))

	lib.LogParseErr(1, "Проверка на корректность значений...")

	k := len(storageMMSData)
	errCount := 0

	for i := 0; i < k; i++ {

		if !(headerData.check(storageMMSData[i], MmsProviders, i)) {
			//fmt.Println("DELETE...", i)
			storageMMSData = removeIndex(storageMMSData, 2)
			k--
			errCount++
		}

	}

	lib.LogParseErr(1,
		fmt.Sprintf("Проверка корректности произведена. Записей %v", len(storageMMSData)))

}

func FetchVoicesCall(filename string) (parseErrCount int) {
	lineCounter := 0

	lib.LogParseErr(1, "Открытие файла: "+filename)

	file, err := os.Open(filename)
	if err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		splittedString := strings.Split(scanner.Text(), csvSeparator)

		if len(splittedString) != 8 {
			lib.LogParseErr(4,
				fmt.Sprintf("Ошибка количества элементов. Строка: %v", lineCounter))
			parseErrCount++
			continue
		}

		s := voiceCallData{
			header: headerData{
				Country:      splittedString[0],
				Bandwidth:    splittedString[1],
				ResponseTime: splittedString[2],
				Provider:     splittedString[3],
			},
		}

		if !(s.header.check(VoiceCallProviders, lineCounter)) {
			parseErrCount++
			continue
		}

		b, err := strconv.ParseFloat(splittedString[4], 32)
		if err != nil {
			lib.LogParseErr(3,
				fmt.Sprintf(" connectionStability: %v. строка %v",
					splittedString[4], lineCounter))
			parseErrCount++
			continue
		}

		s.connectionStability = float32(b)

		s.tTFB, err = strconv.Atoi(splittedString[5])
		if err != nil {
			lib.LogParseErr(3,
				fmt.Sprintf(" tTFB: %v. строка %v",
					splittedString[5], lineCounter))
			parseErrCount++
			continue
		}

		s.voicePurity, err = strconv.Atoi(splittedString[6])
		if err != nil {
			lib.LogParseErr(3,
				fmt.Sprintf(" voicePurity:  %v. строка %v",
					splittedString[6], lineCounter))
			parseErrCount++
			continue
		}

		s.medianOfCallsTime, err = strconv.Atoi(splittedString[7])
		if err != nil {
			lib.LogParseErr(3,
				fmt.Sprintf(" medianOfCallsTime: %v. строка %v",
					splittedString[7], lineCounter))
			parseErrCount++
			continue
		}

		lineCounter++

	}

	lib.LogParseErr(0,
		fmt.Sprintf("Обработано строк без ошибок: %v", lineCounter))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return

}

func FetchEmail(filename string) (parseErrCount int) {
	lineCounter := 0

	lib.LogParseErr(1, "Открытие файла: "+filename)

	file, err := os.Open(filename)
	if err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		/* splittedString := strings.Split(scanner.Text(), csvSeparator)

		log.Println(splittedString) */
		lineCounter++

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return

}
