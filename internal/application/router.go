package application

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"tix-worker/internal/models"

	"github.com/gin-gonic/gin"
)

func (app *Application) registRouter() {
	api := app.ginEngine.Group("/api")
	api.POST("/order/setCaptcha", app.handleSetCaptcha)
	api.POST("/order/:id/getCaptcha", app.handleGetOrderCaptcha)
	api.POST("/order/:id/complete", app.handleCompleteOrder)
}

func (app *Application) handleSetCaptcha(c *gin.Context) {
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

func (app *Application) handleGetOrderCaptcha(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id should not be empty"})
		return
	}

	orderId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id's type should be number"})
		return
	}

	order := models.Order{}
	if err := app.db.Preload("Creator").Where("id = ?", orderId).First(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("order not found, err: %s", err)})
		return
	}

	if err := app.getOrderCaptcha(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("get order captcha failed, err: %s", err)})
		return
	}

	go func() {
		time.Sleep(4*time.Minute + 30*time.Second)
		if err := app.completeOrder(order); err != nil {
			log.Printf("complete order failed, err: %s", err)
		}
	}()

	c.JSON(http.StatusOK, nil)
}

func (app *Application) handleCompleteOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id should not be empty"})
		return
	}

	orderId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id's type should be number"})
		return
	}

	order := models.Order{}
	if err := app.db.Preload("Creator").Where("id = ?", orderId).First(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("order not found, err: %s", err)})
		return
	}

	if err := app.completeOrder(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("complete order failed, err: %s", err)})
		return
	}

	c.JSON(http.StatusOK, nil)
}
