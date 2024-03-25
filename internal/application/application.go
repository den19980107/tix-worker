package application

import (
	"os"
	crawlerpool "tix-worker/internal/crawler-pool"
	"tix-worker/internal/mail"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Application struct {
	db        *gorm.DB
	ginEngine *gin.Engine
	pool      crawlerpool.CrawlerPool
	cron      *cron.Cron
	mail      mail.Mail
}

func New(db *gorm.DB) Application {
	r := gin.Default()
	r.Use(cors.Default())

	mailUser := os.Getenv("MAIL_USER")
	mailPwd := os.Getenv("MAIL_PASSWORD")
	mailSmtpHost := os.Getenv("MAIL_SMTP_HOST")
	mailSmtpPort := os.Getenv("MAIL_SMTP_PORT")

	return Application{
		db:        db,
		ginEngine: r,
		pool:      crawlerpool.New(),
		cron:      cron.New(),
		mail:      mail.New(mailUser, mailPwd, mailSmtpHost, mailSmtpPort),
	}
}

func (app *Application) Run() {
	app.getOrderCaptcha()
	app.registerJob()
	app.registRouter()
	app.ginEngine.Run()
}
