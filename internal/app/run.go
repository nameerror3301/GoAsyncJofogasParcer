package app

import (
	"GoAsyncJofogasParcer/internal/models"
	"fmt"
	"strconv"
	"sync"

	gojson "github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
)

// Константные имена категорий
const (
	Electronic = "Elec"
	Сlothing   = "Сlothing"
	Hobby      = "Hobby"
	BabyMoM    = "BabyMoM"
	Sport      = "Sport"
)

func Run() {
	var wg sync.WaitGroup

	// i == Кол-во страниц которые будут собраны
	for i := 1; i <= 35; i++ {
		wg.Add(1)
		var (
			urlElectronic = fmt.Sprintf("https://www.jofogas.hu/magyarorszag/muszaki-cikkek-elektronika?o=%s", strconv.Itoa(i))
			urlСlothing   = fmt.Sprintf("https://www.jofogas.hu/magyarorszag/otthon-haztartas?o=%s", strconv.Itoa(i))
			urlHobby      = fmt.Sprintf("https://www.jofogas.hu/magyarorszag/szabadido-sport?o=%s", strconv.Itoa(i))
			urlBabyMoM    = fmt.Sprintf("https://www.jofogas.hu/magyarorszag/baba-mama?o=%s", strconv.Itoa(i))
			urlSport      = fmt.Sprintf("https://www.jofogas.hu/magyarorszag/divat-ruhazat?o=%s", strconv.Itoa(i))
		)

		// Электроника
		go func(urlElectronic string) {
			if err := models.FindProduct(urlElectronic, Electronic); err != nil {
				logrus.Println(err)
			}
			wg.Done()
		}(urlElectronic)

		// Одежда
		go func(urlСlothing string) {
			if err := models.FindProduct(urlСlothing, Сlothing); err != nil {
				logrus.Println(err)
			}
			wg.Done()
		}(urlСlothing)

		// Хобби развлечения
		go func(urlHobby string) {
			if err := models.FindProduct(urlHobby, Hobby); err != nil {
				logrus.Println(err)
			}
			wg.Done()
		}(urlHobby)

		// Спорт
		go func(urlSport string) {
			if err := models.FindProduct(urlSport, Sport); err != nil {
				logrus.Println(err)
			}
			wg.Done()
		}(urlSport)

		// Мать и ребенок
		go func(urlBabyMoM string) {
			if err := models.FindProduct(urlBabyMoM, BabyMoM); err != nil {
				logrus.Println(err)
			}
			wg.Done()
		}(urlBabyMoM)

		wg.Wait()
	}

	/*
		Отправка данных в другой микросервис
	*/
	// fmt.Println(string(MarshalData(models.Elec)))
	// fmt.Println()
	// fmt.Println()
	// fmt.Println(string(MarshalData(models.Сlothing)))
	// fmt.Println()
	// fmt.Println()
	// fmt.Println(string(MarshalData(models.Hobby)))
	// fmt.Println()
	// fmt.Println()
	// fmt.Println(string(MarshalData(models.BabyMoM)))
	// fmt.Println()
	// fmt.Println()
	// fmt.Println(string(MarshalData(models.Sport)))

	// fmt.Printf("{%s:%s}", Electronic, MarshalData(models.Elec))
}

func MarshalData(data interface{}) []byte {
	out, err := gojson.Marshal(data)
	if err != nil {
		logrus.Errorf("Err marshal data in struct - %s", err)
	}
	return out
}
