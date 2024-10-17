package depserve

import (
	"github.com/mrpawan-gupta/depserve/internal/server"
)

func RunServer() {
	apiServer := server.NewAPIServer()
	apiServer.InitRouter()
	apiServer.InitMiddleWare()
	//apiServer.InitServer()
	//apiServer.InitDomain()
	//apiServer.RunServer()
}
