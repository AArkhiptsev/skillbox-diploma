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

var SmsProviders = []string{"Topolo", "Rond", "Kildy"} //отсортируем

type SMSData struct {
	country      string
	bandwith     string
	responseTime string
	provider     string
}

func found(s string, a []string) bool {
	return a[sort.SearchStrings(a, s)] == s
}

func (s SMSData) Check() (result bool) {

	result = false

	if GetCountryNameByAlpha(s.country) == "" {
		return
	}

	if !(found(s.provider, SmsProviders)) {
		log.Printf("Нераспознанный провайдер: %v\n", s.provider)
		return
	}

	if _, err := strconv.Atoi(s.responseTime); err != nil {
		log.Printf("Ошибка конвертации. Параметр: среднее время ответа: %v\n", s.responseTime)
		return
	}

	if _, err := strconv.Atoi(s.bandwith); err != nil {
		log.Printf("Ошибка конвертации. Параметр: полоса пропускания: %v\n", s.bandwith)
		return
	}

	result = true

	return

}

func FetchSMS(filename string) {

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
		log.Println("Парсинг:", splittedString)

		if len(splittedString) > 3 {

			s := SMSData{
				country:      splittedString[0],
				bandwith:     splittedString[1],
				responseTime: splittedString[2],
				provider:     splittedString[3],
			}

			s.Check()

		} else {
			log.Printf("Ошибка в строке: %d", lineCounter)
		}

		lineCounter++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
