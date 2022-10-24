package bootstrap

import "github.com/gin-gonic/gin"

func SetupAll(engine *gin.Engine) {
	SetupConfig()
	SetupServer(engine)
	SetupDatabase()
	SetupMailer()
	SetupLogger()
}
