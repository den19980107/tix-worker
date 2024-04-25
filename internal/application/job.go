package application

import (
	"log"
	"time"
	"tix-worker/internal/models"
)

func (app *Application) registerJob() {
	// running at every day 11:59 pm at utc time
	log.Println("add get order captcha job at 12:00 pm")
	_, err := app.cron.AddFunc("00 16 * * *", func() {
		app.getOrdersCaptcha()
	})

	if err != nil {
		panic(err)
	}

	app.cron.Start()
}

func (app *Application) getOrdersCaptcha() {
	log.Printf("running get order captcha job ...")
	orders := app.getTodayOrders()
	for _, order := range orders {
		go func(order models.Order) {
			log.Printf("get %s's order captcha ...", order.Creator.Username)
			err := app.getOrderCaptcha(order)
			if err != nil {
				log.Printf("get %s's order captcha failed, err: %s", order.Creator.Username, err)
			}
		}(order)
	}
}

func (app *Application) getTodayOrders() []models.Order {
	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Printf("get location of Asia/Taipei failed, err: %s", err)
		return []models.Order{}
	}

	now := time.Now().In(location)

	return app.getOrderInDate(now.Year(), now.Month(), now.Day())
}

func (app *Application) getOrderInDate(year int, month time.Month, day int) []models.Order {
	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Printf("get location of Asia/Taipei failed, err: %s", err)
		return []models.Order{}
	}

	dateStart := time.Date(year, month, day, 0, 0, 0, 0, location).UTC()
	dateEnd := dateStart.Add(24 * time.Hour)

	orders := []models.Order{}
	app.db.Preload("Creator").Where("\"execDay\" >= ? AND \"execDay\" < ? AND \"status\" = ?", dateStart, dateEnd, models.OrderStatusPending).Find(&orders)

	log.Printf("get %d order between %s ~ %s", len(orders), dateStart, dateEnd)
	return orders
}
