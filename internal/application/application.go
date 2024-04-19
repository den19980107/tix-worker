package application

import (
	"fmt"
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
	tixUrl    string
}

func New(db *gorm.DB) Application {
	r := gin.Default()
	r.Use(cors.Default())

	mailUser := os.Getenv("MAIL_USER")
	mailPwd := os.Getenv("MAIL_PASSWORD")
	mailSmtpHost := os.Getenv("MAIL_SMTP_HOST")
	mailSmtpPort := os.Getenv("MAIL_SMTP_PORT")
	tixUrl := os.Getenv("NEXTAUTH_URL")

	return Application{
		db:        db,
		ginEngine: r,
		pool:      crawlerpool.New(),
		cron:      cron.New(),
		mail:      mail.New(mailUser, mailPwd, mailSmtpHost, mailSmtpPort),
		tixUrl:    tixUrl,
	}
}

func (app *Application) Run() {
	app.registerJob()
	app.registRouter()
	err := app.ginEngine.Run()
	if err != nil {
		panic(fmt.Sprintf("run gin engine failed, err: %s", err))
	}
}
