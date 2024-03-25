package application

import (
	"fmt"
	"tix-worker/internal/crawler"
	"tix-worker/internal/models"
)

const tixUrl = "http://192.168.31.149:3003"

func (app *Application) GetOrderCaptcha(order models.Order) error {
	c := crawler.Create()
	captcha, jsessionId, err := c.GetCaptchaImageAndJsessionId()
	if err != nil {
		return fmt.Errorf("get captcha image and jsession id failed, err: %s", err)
	}

	err = app.db.Model(&models.Order{}).Where("id = ?", order.Id).Updates(models.Order{Captcha: captcha, JsessionId: jsessionId}).Error
	if err != nil {
		return fmt.Errorf("update order captcha and jsession id failed, err: %s", err)
	}

	err = app.mail.Send(order.Creator.Username, fmt.Sprintf("請至 %s/order/thsrc/%d 填寫驗證碼", tixUrl, order.Id))
	if err != nil {
		return err
	}

	app.pool.Set(order.Id, c)
	return nil
}
