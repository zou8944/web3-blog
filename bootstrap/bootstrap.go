package bootstrap

import "github.com/gin-gonic/gin"

func SetupAll(engine *gin.Engine) {
	SetupServer(engine)
	SetupDatabase()
	SetupMailer()
}
