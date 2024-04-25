package application

import (
	"fmt"
	"tix-worker/internal/crawler"
	"tix-worker/internal/models"
)

func (app *Application) getOrderCaptcha(order models.Order) error {
	c := crawler.Create(order)
	captcha, jsessionId, err := c.GetCaptchaImageAndJsessionId()
	if err != nil {
		return fmt.Errorf("get captcha image and jsession id failed, err: %s", err)
	}

	err = app.db.Model(&models.Order{}).Where("id = ?", order.Id).Updates(models.Order{Captcha: captcha, JsessionId: jsessionId}).Error
	if err != nil {
		return fmt.Errorf("update order: %d captcha and jsession id failed, err: %s", order.Id, err)
	}

	err = app.mail.Send(order.Creator.Username, fmt.Sprintf("請至 %s/order/thsrc/%d 填寫驗證碼", app.tixUrl, order.Id))
	if err != nil {
		return err
	}

	app.pool.Set(order.Id, c)
	return nil
}
