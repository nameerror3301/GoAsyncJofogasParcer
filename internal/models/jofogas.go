package models

import (
	"GoAsyncJofogasParcer/internal/config"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"

	gojson "github.com/goccy/go-json"
)

type RequesLast struct {
	User struct {
		Name        string `json:"user_name"`
		DateRegistr string `json:"user_date_reg"`
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

type ProxyData struct {
	Type string `json:"type"`
	Data struct {
		IP   string `json:"ip"`
		Port string `json:"port"`
	} `json:"data"`
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
	conf := config.ReadConfig()
	data, err := RequestFromParce(url, conf.Data.JwtToken)
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

		c.OnRequest(func(r *colly.Request) {
			r.ProxyURL = FindProxy(conf.Data.JwtToken)
		})

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
				AppendData(&Elec, FindPhone(productID, FindProxy(conf.Data.JwtToken)), productName, photoUrl, price, description, datePublicate, val)
			case "Сlothing":
				AppendData(&Сlothing, FindPhone(productID, FindProxy(conf.Data.JwtToken)), productName, photoUrl, price, description, datePublicate, val)
			case "Hobby":
				AppendData(&Hobby, FindPhone(productID, FindProxy(conf.Data.JwtToken)), productName, photoUrl, price, description, datePublicate, val)
			case "BabyMoM":
				AppendData(&BabyMoM, FindPhone(productID, FindProxy(conf.Data.JwtToken)), productName, photoUrl, price, description, datePublicate, val)
			case "Sport":
				AppendData(&Sport, FindPhone(productID, FindProxy(conf.Data.JwtToken)), productName, photoUrl, price, description, datePublicate, val)
			}
		})
		c.Visit(val)
	})
	return nil
}

func AppendData(data *[]RequesLast, phone string, respData ...string) {
	*data = append(*data, RequesLast{
		User: struct {
			Name        string "json:\"user_name\""
			DateRegistr string "json:\"user_date_reg\""
			PhoneNumber string "json:\"user_phone\""
		}{
			Name:        "No Data",
			DateRegistr: "No Data",
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

func FindPhone(id string, proxy string) string {
	var p PhoneNum

	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		logrus.Errorf("Err parce proxy url - %s", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://apiv2.jofogas.hu/v2/items/getPhone?list_id=%s", id), nil)
	if err != nil {
		logrus.Errorf("Err generate request from phone - %s", err)
		return ""
	}
	req.Header.Add("api_key", "jofogas-web-eFRv9myucHjnXFbj")

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Err request from finding phone - %s", err)
		return ""
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Err parce phone body - %s", err)
		return ""
	}

	if err = gojson.Unmarshal(data, &p); err != nil {
		logrus.Errorf("Err unmarshal data to struct phone - %s", err)
		return ""
	}
	return p.Phone
}

func RequestFromParce(urlFromParce string, token string) (io.ReadCloser, error) {
	proxyUrl, err := url.Parse(FindProxy(token))
	if err != nil {
		logrus.Errorf("Err parce proxy url - %s", err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	resp, err := client.Get(urlFromParce)
	if err != nil {
		logrus.Fatalf("Err request to %s - %s", urlFromParce, err)
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Err responce - %d %s", resp.StatusCode, resp.Status)
		time.Sleep(time.Second * 5)
	}

	fmt.Println(resp.StatusCode)
	return resp.Body, nil
}

func FindProxy(token string) string {
	conf := config.ReadConfig()
	var p ProxyData

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, conf.Data.OutProxyAddr, nil)
	if err != nil {
		logrus.Errorf("Err generate request from finding proxy - %s", err)
		return ""
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Err request to proxy service - %s", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Err read body - %s", err)
		return ""
	}
	if err := gojson.Unmarshal(body, &p); err != nil {
		logrus.Errorf("Err unmarshal data to struct - %s", err)
		return ""
	}
	return fmt.Sprintf("http://%s:%s", p.Data.IP, p.Data.Port)
}
