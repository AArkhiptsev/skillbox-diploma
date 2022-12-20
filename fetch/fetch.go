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
	smsSeparator = ";"

	x
)

var (
	SmsProviders = []string{"Topolo", "Rond", "Kildy"}
	mmsProviders = SmsProviders //работаем с копией, сделано на случай,
	// если в будущем появится иной набор провайдеров для MMS
)

type smsData struct {
	country      string
	bandwith     string
	responseTime string
	provider     string
}

type mmsData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

var (
	storageSMSData = []smsData{}
	storageMMSData = []mmsData{}
)

func (s smsData) Check() (result bool) {

	result = false

	if lib.GetCountryNameByAlpha(s.country) == "" {
		lib.LogParseErr(3, " alpha: "+s.country)
		return
	}

	if !(lib.Found(s.provider, SmsProviders)) {
		lib.LogParseErr(3, " провайдер: "+s.provider)
		return
	}

	if _, err := strconv.Atoi(s.responseTime); err != nil {
		lib.LogParseErr(3, " среднее время ответа: "+s.responseTime)
		return
	}

	if _, err := strconv.Atoi(s.bandwith); err != nil {
		lib.LogParseErr(3, " полоса пропускания: "+s.bandwith)
		return
	}

	result = true

	return

}

func (m mmsData) Check() (result bool) {

	result = false

	if lib.GetCountryNameByAlpha(m.Country) == "" {
		lib.LogParseErr(3, " alpha: "+m.Country)
		return
	}

	if !(lib.Found(m.Provider, mmsProviders)) {
		lib.LogParseErr(3, " провайдер: "+m.Provider)
		return
	}
	if _, err := strconv.Atoi(m.ResponseTime); err != nil {
		lib.LogParseErr(3, " среднее время ответа: "+m.ResponseTime)
		return
	}

	if _, err := strconv.Atoi(m.Bandwidth); err != nil {
		lib.LogParseErr(3, " полоса пропускания: "+m.Bandwidth)
		return
	}

	result = true

	return
}

func LogStorageSMSData() {
	for _, datum := range storageSMSData {
		log.Println(datum)
	}
}

func LogStorageMMSData() {
	for _, datum := range storageMMSData {
		log.Println(datum)
	}
}

func FetchSMS(filename string) (parseErrCount int) {

	lineCounter := 1

	lib.LogParseErr(1, "Открытие файла: "+filename)

	file, err := os.Open(filename)
	if err != nil {
		lib.LogParseErr(4, err.Error())
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		splittedString := strings.Split(scanner.Text(), smsSeparator)
		//log.Println("Парсинг:", splittedString)

		if len(splittedString) > 3 {

			s := smsData{
				country:      splittedString[0],
				bandwith:     splittedString[1],
				responseTime: splittedString[2],
				provider:     splittedString[3],
			}

			if s.Check() {
				storageSMSData = append(storageSMSData, s)
			} else {
				parseErrCount++
			}

		} else {
			lib.LogParseErr(4, "Ошибка количества элементов. Строка: "+
				strconv.Itoa(lineCounter))
			parseErrCount++
		}

		lineCounter++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return

}

func removeIndex(s []mmsData, index int) []mmsData {
	return append(s[:index], s[index+1:]...)
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

		if !(mmsData.Check(storageMMSData[i])) {
			//fmt.Println("DELETE...", i)
			storageMMSData = removeIndex(storageMMSData, 2)
			k--
			errCount++
		}

	}

	lib.LogParseErr(1,
		fmt.Sprintf("Проверка корректности произведена. Записей %v", len(storageMMSData)))

}
