package fetch

import (
	"bufio"
	"diploma/lib"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	smsSeparator  = ";"
	mmsDataServer = "localhost:8383"
	x
)

var SmsProviders = []string{"Topolo", "Rond", "Kildy"} //отсортируем

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

var storageSMSData = make([]smsData, 0)

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

func LogStorageSMSData() {
	for _, datum := range storageSMSData {
		log.Println(datum)
	}
}

func FetchSMS(filename string) (parseErrCount int) {

	lineCounter := 1

	lib.LogParseErr(1, "Открытие файла: "+filename)

	file, err := os.Open(filename)
	if err != nil {
		lib.LogParseErr(3, err.Error())
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
			lib.LogParseErr(3, "Ошибка количества элементов. Строка: "+
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
