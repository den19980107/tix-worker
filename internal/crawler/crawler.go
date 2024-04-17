package crawler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"tix-worker/internal/models"

	"github.com/gocolly/colly/v2"
)

const (
	BASE_URL           = "https://irs.thsrc.com.tw"
	BOOKING_PAGE_URL   = "https://irs.thsrc.com.tw/IMINT/?locale=tw"
	SUBMIT_FORM_URL    = "https://irs.thsrc.com.tw/IMINT/;jsessionid=%s?wicket:interface=:0:BookingS1Form::IFormSubmitListener"
	CONFIRM_TRAIN_URL  = "https://irs.thsrc.com.tw/IMINT/?wicket:interface=:1:BookingS2Form::IFormSubmitListener"
	CONFIRM_TICKET_URL = "https://irs.thsrc.com.tw/IMINT/?wicket:interface=:2:BookingS3Form::IFormSubmitListener"

	USER_AGENT      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"
	ACCEPT_HTML     = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
	ACCEPT_IMG      = "image/webp,*/*"
	ACCEPT_LANGUAGE = "zh-TW,zh;q=0.8,en-US;q=0.5,en;q=0.3"
	ACCEPT_ENCODING = "gzip, deflate, br"

	BOOKING_PAGE_HOST = "irs.thsrc.com.tw"
)

type Crawler struct {
	collector *colly.Collector
	captcha   *string
}

func Create() Crawler {
	return Crawler{
		collector: colly.NewCollector(
			colly.AllowURLRevisit(),
			colly.MaxDepth(5),
			colly.UserAgent(USER_AGENT),
		),
	}
}

func (c *Crawler) SetCaptcha(captcha string) {
	c.captcha = &captcha
}

func (c *Crawler) CompleteOrder(order models.Order) error {
	if c.captcha == nil {
		return fmt.Errorf("crawler dosent have filled captcha!")
	}

	trainDatas, feedBackErrors, err := c.submitForm(order, *c.captcha, order.JsessionId)
	if err != nil {
		return fmt.Errorf("submit train order failed, err: %s", err)
	}

	if len(feedBackErrors) > 0 {
		errorStr := ""
		for index, err := range feedBackErrors {
			errorStr += err.Error()
			if index != len(feedBackErrors)-1 {
				errorStr += "\n"
			}
		}
		return errors.New(errorStr)
	}

	if len(trainDatas) == 0 {
		return errors.New("no train avaliable")
	}

	validTrain := c.filterValidTrain(order, trainDatas)
	if len(validTrain) == 0 {
		return errors.New("no valid train avaliable")
	}

	err = c.confirmTrain(validTrain[0])
	if err != nil {
		return fmt.Errorf("confirm train failed, err: %s", err)
	}

	err = c.confirmTicket(order)
	if err != nil {
		return fmt.Errorf("confirm ticket failed, err: %s", err)
	}

	return nil
}

func (c *Crawler) GetCaptchaImageAndJsessionId() (imageBase64 string, jsessionId string, err error) {
	var wg sync.WaitGroup
	wg.Add(1)

	collector := cloneCollectorWithBasicEvent(c.collector, "get captcha code and jession id")
	collector.OnHTML("img[id='BookingS1Form_homeCaptcha_passCode']", func(e *colly.HTMLElement) {
		siteCookies := collector.Cookies(e.Request.URL.String())
		for _, cookie := range siteCookies {
			if cookie.Name == "JSESSIONID" {
				jsessionId = cookie.Value
			}
		}

		imgPath := e.Attr("src")
		url := BASE_URL + imgPath
		imageBase64, err = imageUrlToBase64(url)
		if err != nil {
			err = fmt.Errorf("convert image to base64 failed, err: %s", err)
		}

		wg.Done()
	})

	err = collector.Visit(BOOKING_PAGE_URL)
	if err != nil {
		return "", "", fmt.Errorf("visit booking page failed, err: %s", err)
	}

	wg.Wait()

	return imageBase64, jsessionId, err
}

func (c *Crawler) submitForm(order models.Order, captchaResult string, jsessionId string) ([]TrainData, []FeedBackError, error) {
	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return nil, nil, fmt.Errorf("get location of Asia/Taipei failed, err: %s", err)
	}

	submitForm := SubmitForm{
		TripConTypesoftrip:            "0",
		TrainConTrainRadioGroup:       "0",
		SeatConSeatRadioGroup:         "0",
		BookingMethod:                 "radio31",
		SelectStartStation:            order.From.Code(),
		SelectDestinationStation:      order.To.Code(),
		ToTimeInputField:              order.DepartureDay.In(location).Format("2006/01/02"),
		BackTimeInputField:            order.DepartureDay.In(location).Format("2006/01/02"),
		ToTimeTable:                   order.GetStartTime(),
		TicketPanelRows0TicketAmount:  "1F",
		TicketPanelRows1TicketAmount:  "0H",
		TicketPannelRows2TicketAmount: "0W",
		TicketPanelRows3TicketAmount:  "0E",
		TicketPanelRows4TicketAmount:  "0P",
		HomeCaptchaSecurityCode:       captchaResult,
	}

	submitJsonStr, _ := json.MarshalIndent(submitForm, "", " ")
	log.Printf("form:\n%s", submitJsonStr)
	submitBody := make(map[string]string)
	err = json.Unmarshal(submitJsonStr, &submitBody)
	if err != nil {
		return nil, nil, err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	trainDatas := []TrainData{}
	feedBackErrors := []FeedBackError{}

	collector := cloneCollectorWithBasicEvent(c.collector, "submit form")

	// detect predict captcha failed
	collector.OnHTML("li[class='feedbackPanelERROR']", func(h *colly.HTMLElement) {
		log.Printf("on html get feed back")
		h.ForEach("span", func(i int, h *colly.HTMLElement) {
			feedBackErrors = append(feedBackErrors, GetFeedBackError(h.Text))
		})
	})

	// find the error content message
	collector.OnHTML("div[class='error-content']", func(h *colly.HTMLElement) {
		log.Printf("on html get error content")
		h.ForEachWithBreak("li", func(i int, h *colly.HTMLElement) bool {
			feedBackError := GetFeedBackError(h.Text)
			// if <li> dosent contain the knowed error message, continue
			if feedBackError.Type == Unknow {
				return true
			}

			// if <li> have knowed error message, append to the feed back error and break
			feedBackErrors = append(feedBackErrors, feedBackError)
			return false
		})
	})

	// if predict succes, should get response of list of train
	collector.OnHTML("div[class='result-listing']", func(resultList *colly.HTMLElement) {
		log.Printf("on html get result")
		resultList.ForEach("label[class='uk-flex uk-flex-middle result-item']", func(i int, resultItem *colly.HTMLElement) {
			trainCode := resultItem.ChildText("#QueryCode")
			departureTime := resultItem.ChildText("#QueryDeparture")
			arrivalTime := resultItem.ChildText("#QueryArrival")
			departureDate := resultItem.ChildText("#QueryDepartureDate")
			value := resultItem.ChildAttr(".btn-radio > input", "value")

			trainData := TrainData{
				TrainCode:     trainCode,
				Value:         value,
				Date:          departureDate,
				DepartureTime: departureTime,
				ArrivalTime:   arrivalTime,
			}

			log.Printf("get train data %+v", trainData)
			trainDatas = append(trainDatas, trainData)
		})
	})

	collector.OnScraped(func(r *colly.Response) {
		log.Printf("on scraped")
		defer wg.Done()
	})

	err = collector.Post(fmt.Sprintf(SUBMIT_FORM_URL, jsessionId), submitBody)
	wg.Wait()

	if len(feedBackErrors) > 0 {
		return nil, feedBackErrors, nil
	}

	if err != nil {
		return nil, nil, err
	}

	log.Printf("return trains data: %+v", trainDatas)
	return trainDatas, nil, nil
}

func (c *Crawler) confirmTrain(trainData TrainData) error {
	collector := cloneCollectorWithBasicEvent(c.collector, "confirm train")
	confirmTrainForm := ConfirmTrain{
		BookingS2FormHf0:                  "",
		TrainQueryDataViewPanelTrainGroup: trainData.Value,
	}

	confirmTrainJsonStr, _ := json.Marshal(confirmTrainForm)
	confirmTrainBody := make(map[string]string)
	err := json.Unmarshal(confirmTrainJsonStr, &confirmTrainBody)
	if err != nil {
		return fmt.Errorf("unmarshal confirm train json string failed, err: %s", err)
	}

	return collector.Post(CONFIRM_TRAIN_URL, confirmTrainBody)
}

func (c *Crawler) confirmTicket(order models.Order) error {
	collector := cloneCollectorWithBasicEvent(c.collector, "confirm ticket")
	confirmTicketForm := ConfirmTicket{
		BookingS3FormSpHf0: "",
		DiffOver:           "1",
		IsSPromotion:       "1",
		PassengerCount:     "1",
		IsGoBackM:          "",
		BackHome:           "",
		TgoError:           "1",
		IdInputRadio:       "0",
		DummyId:            order.Creator.IdNumber,
		DummyPhone:         order.Creator.PhoneNumber,
		Email:              order.Creator.Username,
		Agree:              "on",
		TicketMemberSystemInputPanelTakerMemberSystemDataViewMemberSystemRadioGroup: "radio44",
	}
	confirmTicketJsonStr, _ := json.Marshal(confirmTicketForm)
	confirmTicketBody := make(map[string]string)
	err := json.Unmarshal(confirmTicketJsonStr, &confirmTicketBody)
	if err != nil {
		return fmt.Errorf("unmarshal confirm ticket json string failed, err: %s", err)
	}

	return collector.Post(CONFIRM_TICKET_URL, confirmTicketBody)
}

func cloneCollectorWithBasicEvent(collector *colly.Collector, execName string) *colly.Collector {
	clone := collector.Clone()
	clone.OnRequest(func(r *colly.Request) {
		r.Headers.Add("HOST", BOOKING_PAGE_HOST)
		r.Headers.Add("Accept", ACCEPT_HTML)
		r.Headers.Add("Accept-Language", ACCEPT_LANGUAGE)
		r.Headers.Add("Accept-Encoding", ACCEPT_ENCODING)
		r.Headers.Add("Connection", "keep-alive")
		r.Headers.Add("Content-Type", "application/x-www-form-urlencoded")
		// log.Println("Visiting", r.URL)
	})

	clone.OnError(func(r *colly.Response, err error) {
		log.Printf("%s", err.Error())
	})

	clone.OnResponse(func(r *colly.Response) {
		// log.Printf("%s on response with status code: %d", execName, r.StatusCode)
		if r.StatusCode == 200 {
			_ = os.WriteFile(fmt.Sprintf("%s.html", execName), r.Body, 0644)
		}
	})

	return clone
}

func imageUrlToBase64(URL string) (string, error) {
	//Get the response bytes from the url
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Accept", ACCEPT_HTML)
	req.Header.Set("Accept-Language", ACCEPT_LANGUAGE)
	req.Header.Set("Accept-Encoding", ACCEPT_ENCODING)

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", errors.New("Received non 200 response code")
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return "", err
	}
	encodedStr := base64.StdEncoding.EncodeToString(buf.Bytes())
	return encodedStr, nil
}

func (c *Crawler) filterValidTrain(order models.Order, trainDatas []TrainData) []TrainData {
	validTrains := []TrainData{}
	for _, trainData := range trainDatas {
		trainDepartureHour, err := trainData.getHour()
		if err != nil {
			log.Printf("failed to get train departure hour, err: %s", err)
			continue
		}

		trainDepartureMin, err := trainData.getMin()
		if err != nil {
			log.Printf("failed to get train departue min, err: %s", err)
			continue
		}

		trainDepartureTime := time.Date(order.DepartureDay.Year(), order.DepartureDay.Month(), order.DepartureDay.Day(), trainDepartureHour, trainDepartureMin, 0, 0, order.DepartureDay.Location())
		validStartTime := time.Date(order.DepartureDay.Year(), order.DepartureDay.Month(), order.DepartureDay.Day(), order.GetStartHour(), order.GetStartMin(), 0, 0, order.DepartureDay.Location())
		validEndTime := time.Date(order.DepartureDay.Year(), order.DepartureDay.Month(), order.DepartureDay.Day(), order.GetEndHour(), order.GetEndMin(), 0, 0, order.DepartureDay.Location())

		if trainDepartureTime.Equal(validStartTime) || trainDepartureTime.After(validStartTime) && trainDepartureTime.Before(validEndTime) || trainDepartureTime.Equal(validEndTime) {
			validTrains = append(validTrains, trainData)
		}
	}

	return validTrains
}
