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
	MmsProviders       = SmsProviders //работаем с копией, сделано на случай,
	// если в будущем появится иной набор провайдеров для MMS
)

type smsData struct {
	country      string
	bandwidth    string
	responseTime string
	provider     string
}

type mmsData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

type VoiceCallData struct {
	Country             string
	Bandwidth           string
	ResponseTime        string
	Provider            string
	ConnectionStability float32
	TTFB                int
	VoicePurity         int
	MedianOfCallsTime   int
}

var (
	storageSMSData       = []smsData{}
	storageMMSData       = []mmsData{}
	storageVoiceCallData = []VoiceCallData{}
)

func removeIndex(s []mmsData, index int) []mmsData {
	return append(s[:index], s[index+1:]...)
}

func (s smsData) check() (result bool) {

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

	if _, err := strconv.Atoi(s.bandwidth); err != nil {
		lib.LogParseErr(3, " полоса пропускания: "+s.bandwidth)
		return
	}

	result = true

	return

}

func (m mmsData) check() (result bool) {

	result = false

	if lib.GetCountryNameByAlpha(m.Country) == "" {
		lib.LogParseErr(3, " alpha: "+m.Country)
		return
	}

	if !(lib.Found(m.Provider, MmsProviders)) {
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

func (v VoiceCallData) check(
	ConnectionStability,
	TTFB, VoicePurity, MedianOfCallsTime string) (result bool) {

	result = false

	if lib.GetCountryNameByAlpha(v.Country) == "" {
		lib.LogParseErr(3, " alpha: "+v.Country)
		return
	}

	if !(lib.Found(v.Provider, VoiceCallProviders)) {
		lib.LogParseErr(3, " провайдер: "+v.Provider)
		return
	}

	if a, err := strconv.Atoi(VoicePurity); err != nil {
		lib.LogParseErr(3, " VoicePurity: "+VoicePurity)
		return
	} else {
		v.VoicePurity = a
	}

	if a, err := strconv.Atoi(MedianOfCallsTime); err != nil {
		lib.LogParseErr(3, " MedianOfCallsTime: "+MedianOfCallsTime)
		return
	} else {
		v.MedianOfCallsTime = a
	}

	if b, err := strconv.ParseFloat(ConnectionStability, 32); err != nil {
		lib.LogParseErr(3, " ConnectionStability: "+ConnectionStability)
		return
	} else {
		v.ConnectionStability = float32(b)
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

func LogStorageVoicesCallsData() {
	for _, datum := range storageVoiceCallData {
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

		splittedString := strings.Split(scanner.Text(), csvSeparator)
		//log.Println("Парсинг:", splittedString)

		if len(splittedString) == 4 {

			s := smsData{
				country:      splittedString[0],
				bandwidth:    splittedString[1],
				responseTime: splittedString[2],
				provider:     splittedString[3],
			}

			if s.check() {
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

		if !(mmsData.check(storageMMSData[i])) {
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

		splittedString := strings.Split(scanner.Text(), csvSeparator)
		//log.Println("Парсинг:", splittedString)

		if len(splittedString) == 8 {

			s := VoiceCallData{
				Country:      splittedString[0],
				Bandwidth:    splittedString[1],
				ResponseTime: splittedString[2],
				Provider:     splittedString[3],
			}

			if a, err := strconv.Atoi(splittedString[5]); err != nil {
				lib.LogParseErr(3, " TTFB: "+splittedString[5])
				return
			} else {
				s.TTFB = a
				fmt.Println("!", s.TTFB)
			}

			if s.check(splittedString[4],
				splittedString[5],
				splittedString[6],
				splittedString[7]) {
				fmt.Println(s.TTFB)
				storageVoiceCallData = append(storageVoiceCallData, s)
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
