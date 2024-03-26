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

	return app.getOrderInDate(tomorrowNow.Year(), tomorrowNow.Month(), tomorrowNow.Day())
}

func (app *Application) getOrderInDate(year int, month time.Month, day int) []models.Order {
	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Printf("get location of Asia/Taipei failed, err: %s", err)
		return []models.Order{}
	}

	dateStart := time.Date(year, month, day, 0, 0, 0, 0, location)
	dateEnd := dateStart.Add(24 * time.Hour)

	orders := []models.Order{}
	app.db.Preload("Creator").Where("\"execDay\" >= ? AND \"execDay\" < ? AND \"status\" = ?", dateStart, dateEnd, models.OrderStatusPending).Find(&orders)

	log.Printf("get %d order between %s ~ %s", len(orders), dateStart, dateEnd)
	return orders
}

func (app *Application) completeOrder() {
	log.Printf("running complete order job ...")

	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Printf("get location of Asia/Taipei failed, err: %s", err)
		return
	}

	now := time.Now().In(location)
	orders := app.getOrderInDate(now.Year(), now.Month(), now.Day())
	for _, order := range orders {
		crawler, err := app.pool.Get(order.Id)
		if err != nil {
			log.Printf("get order %d's crawler failed, err: %s", order.Id, err)
			err := app.db.Model(&models.Order{}).Where("id = ?", order.Id).Updates(models.Order{Status: models.OrderStatusFailed}).Error
			if err != nil {
				log.Printf("update order %d's status to %s failed, err: %s", order.Id, models.OrderStatusFailed, err)
			}
			continue
		}

		log.Printf("exec order %+v ...", orders)
		err = crawler.CompleteOrder(order)
		if err != nil {
			log.Printf("complete order %+v failed, err: %s", order, err)
			err := app.db.Model(&models.Order{}).Where("id = ?", order.Id).Updates(models.Order{Status: models.OrderStatusFailed}).Error
			if err != nil {
				log.Printf("update order %d's status to %s failed, err: %s", order.Id, models.OrderStatusFailed, err)
			}
			continue
		}

		err = app.db.Model(&models.Order{}).Where("id = ?", order.Id).Updates(models.Order{Status: models.OrderStatusComplete}).Error
		if err != nil {
			log.Printf("update order %d's status to %s failed, err: %s", order.Id, models.OrderStatusFailed, err)
		}
	}
}