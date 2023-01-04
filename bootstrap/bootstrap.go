package bootstrap

import "github.com/gin-gonic/gin"

func SetupAll(engine *gin.Engine) {
	SetupConfig()
	SetupLogger()
	SetupServer(engine)
	SetupRequest()
	SetupDatabase()
	SetupMailer()
	SetupIPFS()
}
