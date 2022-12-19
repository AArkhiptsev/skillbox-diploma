package lib

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const smsSeparator = ";"
const convertErrText = "Ошибка конвертации. Параметр:"

var SmsProviders = []string{"Topolo", "Rond", "Kildy"} //отсортируем

type smsData struct {
	country      string
	bandwith     string
	responseTime string
	provider     string
}

var storageSMSData = make([]smsData, 0)

func found(s string, a []string) bool {
	return a[sort.SearchStrings(a, s)] == s
}

func LogStorageSMSData() {
	for _, datum := range storageSMSData {
		log.Println(datum)
	}
}

func (s smsData) Check() (result bool) {

	result = false

	if GetCountryNameByAlpha(s.country) == "" {
		return
	}

	if !(found(s.provider, SmsProviders)) {
		log.Printf(convertErrText+" провайдер: %v\n", s.provider)
		return
	}

	if _, err := strconv.Atoi(s.responseTime); err != nil {
		log.Printf(convertErrText+" среднее время ответа: %v\n", s.responseTime)
		return
	}

	if _, err := strconv.Atoi(s.bandwith); err != nil {
		log.Printf(convertErrText+" полоса пропускания: %v\n", s.bandwith)
		return
	}

	result = true

	return

}

func FetchSMS(filename string) (parseErrCount int) {

	lineCounter := 1

	log.Printf("Открытие файла %s", filename)

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
			log.Printf("Ошибка количества элементов. Строка: %d", lineCounter)
			parseErrCount++
		}

		lineCounter++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return

}
