package lib

import (
	"bufio"
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

	if GetCountryNameByAlpha(s.country) == "" {
		logParseErr(" alpha: %v", s.country)
		return
	}

	if !(found(s.provider, SmsProviders)) {
		logParseErr(" провайдер: %v", s.provider)
		return
	}

	if _, err := strconv.Atoi(s.responseTime); err != nil {
		logParseErr(" среднее время ответа: %v", s.responseTime)
		return
	}

	if _, err := strconv.Atoi(s.bandwith); err != nil {
		logParseErr(" полоса пропускания: %v", s.bandwith)
		return
	}

	result = true

	return

}

func FetchSMS(filename string) (parseErrCount int) {

	lineCounter := 1

	log.Printf(ColorCyan+"Открытие файла %s"+ColorReset, filename)

	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
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
			log.Printf("Ошибка количества элементов. Строка: "+
				ColorRed+"%v"+ColorReset, lineCounter)
			parseErrCount++
		}

		lineCounter++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return

}
