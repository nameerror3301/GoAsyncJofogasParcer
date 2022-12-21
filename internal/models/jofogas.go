package models

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"

	gojson "github.com/goccy/go-json"
)

type RequesLast struct {
	User struct {
		PhoneNumber string `json:"user_phone"`
	} `json:"user_data"`
	Products struct {
		ProdName    string `json:"product_name"`
		PhotoUrl    string `json:"photo_url"`
		Price       string `json:"price"`
		Description string `json:"description"`
		Url         string `json:"url"`
	} `json:"product_data"`
}

type PhoneNum struct {
	Phone string `json:"phone"`
}

var (
	Elec     []RequesLast
	Сlothing []RequesLast
	Hobby    []RequesLast
	BabyMoM  []RequesLast
	Sport    []RequesLast
)

func FindProduct(url string, category string) error {
	data, err := Request(url)
	if err != nil {
		return err
	}
	defer data.Close()

	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		logrus.Errorf("Err load data - %s", err)
	}
	defer data.Close()

	doc.Find(".item-title a").Each(func(_ int, s *goquery.Selection) {
		val, _ := s.Attr("href")
		c := colly.NewCollector()

		c.OnHTML("html", func(e *colly.HTMLElement) {
			productID, _ := e.DOM.Find("vi-touch-stone[data-list-id]").Attr("data-list-id")
			productName := e.DOM.Find("title").Text()
			photoUrl, _ := e.DOM.Find("a[class=newGalPopUp]").Attr("href")
			price, _ := e.DOM.Find("meta[itemprop=price]").Attr("content")
			datePublicate := e.DOM.Find("span[class=time]").Text()

			var description string
			e.DOM.Find("meta[property]").Each(func(i int, s *goquery.Selection) {
				val, _ := s.Attr("content")
				if i == 1 {
					description = val
				}
			})

			switch category {
			case "Elec":
				AppendData(&Elec, FindPhone(productID), productName, photoUrl, price, description, datePublicate, val)
			case "Сlothing":
				AppendData(&Сlothing, FindPhone(productID), productName, photoUrl, price, description, datePublicate, val)
			case "Hobby":
				AppendData(&Hobby, FindPhone(productID), productName, photoUrl, price, description, datePublicate, val)
			case "BabyMoM":
				AppendData(&BabyMoM, FindPhone(productID), productName, photoUrl, price, description, datePublicate, val)
			case "Sport":
				AppendData(&Sport, FindPhone(productID), productName, photoUrl, price, description, datePublicate, val)
			}
		})

		c.OnError(func(r *colly.Response, err error) {
			logrus.Errorf("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
			time.Sleep(time.Second * 5)
		})

		c.Visit(val)
	})
	return nil
}

func AppendData(data *[]RequesLast, phone string, respData ...string) {
	*data = append(*data, RequesLast{
		User: struct {
			PhoneNumber string "json:\"user_phone\""
		}{
			PhoneNumber: phone,
		},
		Products: struct {
			ProdName    string "json:\"product_name\""
			PhotoUrl    string "json:\"photo_url\""
			Price       string "json:\"price\""
			Description string "json:\"description\""
			Url         string "json:\"url\""
		}{
			ProdName:    respData[0],
			PhotoUrl:    respData[1],
			Price:       respData[2],
			Description: respData[3],
			Url:         respData[5],
		},
	})
}

func FindPhone(id string) string {
	var phone PhoneNum

	url := fmt.Sprintf("https://apiv2.jofogas.hu/v2/items/getPhone?list_id=%s", id)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("api_key", "jofogas-web-eFRv9myucHjnXFbj")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Err parce body - %s", err)
		return ""
	}

	if err = gojson.Unmarshal(data, &phone); err != nil {
		logrus.Warnf("Err unmarshal json - %s", err)
		return ""
	}
	return phone.Phone
}

func Request(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		logrus.Fatalf("Err request to %s - %s", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Err responce - %d %s", resp.StatusCode, resp.Status)
		time.Sleep(time.Second * 5)
	}

	fmt.Println(resp.StatusCode)
	return resp.Body, nil
}
