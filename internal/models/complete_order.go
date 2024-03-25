package models

type CompleteOrderBody struct {
	Id      int    `json:"id"`
	Captcha string `json:"captcha"`
}
