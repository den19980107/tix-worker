package application

import (
	"log"
	"net/http"
	"tix-worker/internal/models"

	"github.com/gin-gonic/gin"
)

func (app *Application) registRouter() {
	api := app.ginEngine.Group("/api")
	api.POST("/order/setCaptcha", app.setCaptcha)
}

func (app *Application) setCaptcha(c *gin.Context) {
	body := models.CompleteOrderBody{}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	crawler, err := app.pool.Get(body.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("set captcha \"%s\" to crawler", body.Captcha)
	crawler.SetCaptcha(body.Captcha)
	c.JSON(http.StatusOK, nil)
}
