package application

import (
	"log"
	"time"
	"tix-worker/internal/models"
)

func (app *Application) registerJob() {
	// running at every day 11:55 pm at utc time
	_, err := app.cron.AddFunc("55 15 * * *", func() {
		app.getOrderCaptcha()
	})

	if err != nil {
		panic(err)
	}

	// running at every day 12:00 am at utc time
	_, err = app.cron.AddFunc("00 16 * * *", func() {
		app.completeOrder()
	})

	if err != nil {
		panic(err)
	}

	app.cron.Start()
}

func (app *Application) getOrderCaptcha() {
	log.Printf("running get order captcha job ...")
	orders := app.getTomorrowOrders()
	for _, order := range orders {
		log.Printf("get order %+v captcha ...", orders)
		err := app.GetOrderCaptcha(order)
		if err != nil {
			log.Printf("get order %+v captcha failed, err: %s", order, err)
		}
	}
}

func (app *Application) getTomorrowOrders() []models.Order {
	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Printf("get location of Asia/Taipei failed, err: %s", err)
		return []models.Order{}
	}

	now := time.Now().In(location)
	tomorrowNow := now.Add(24 * time.Hour)
	tomorrowStart := time.Date(tomorrowNow.Year(), tomorrowNow.Month(), tomorrowNow.Day(), 0, 0, 0, 0, tomorrowNow.Location()).UTC()
	tomorrowEnd := tomorrowStart.Add(24 * time.Hour)

	orders := []models.Order{}
	app.db.Preload("Creator").Where("\"execDay\" >= ? AND \"execDay\" < ?", tomorrowStart, tomorrowEnd).Find(&orders)

	log.Printf("get %d order for tomorrow %s ~ %s", len(orders), tomorrowStart, tomorrowEnd)
	return orders
}

func (app *Application) completeOrder() {
	log.Printf("running complete order job ...")
	orders := app.getTomorrowOrders()
	for _, order := range orders {
		crawler, err := app.pool.Get(order.Id)
		if err != nil {
			log.Printf("get order %d's crawler failed, err: %s", order.Id, err)
			continue
		}

		log.Printf("exec order %+v ...", orders)
		err = crawler.CompleteOrder(order)
		if err != nil {
			log.Printf("complete order %+v failed, err: %s", order, err)
			continue
		}
	}
}
