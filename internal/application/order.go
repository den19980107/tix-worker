package application

import (
	"fmt"
	"log"
	"tix-worker/internal/crawler"
	"tix-worker/internal/models"
)

const tixUrl = "http://192.168.31.149:3003"

func (app *Application) getOrderCaptcha(order models.Order) error {
	c := crawler.Create()
	captcha, jsessionId, err := c.GetCaptchaImageAndJsessionId()
	if err != nil {
		return fmt.Errorf("get captcha image and jsession id failed, err: %s", err)
	}

	err = app.db.Model(&models.Order{}).Where("id = ?", order.Id).Updates(models.Order{Captcha: captcha, JsessionId: jsessionId}).Error
	if err != nil {
		return fmt.Errorf("update order: %d captcha and jsession id failed, err: %s", order.Id, err)
	}

	err = app.mail.Send(order.Creator.Username, fmt.Sprintf("請至 %s/order/thsrc/%d 填寫驗證碼", tixUrl, order.Id))
	if err != nil {
		return err
	}

	app.pool.Set(order.Id, c)
	return nil
}

func (app *Application) completeOrder(order models.Order) error {
	crawler, err := app.pool.Get(order.Id)
	if err != nil {
		return fmt.Errorf("get order %d's crawler failed, err: %s", order.Id, err)
	}
	defer app.pool.Remove(order.Id)

	log.Printf("exec order %+v ...", order)
	err = crawler.CompleteOrder(order)
	if err != nil {
		return err
	}

	err = app.db.Model(&models.Order{}).Where("id = ?", order.Id).Updates(models.Order{Status: models.OrderStatusComplete}).Error
	if err != nil {
		return fmt.Errorf("update order %d's status to %s failed, err: %s", order.Id, models.OrderStatusFailed, err)
	}

	return nil
}
