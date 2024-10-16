package depserve

import "github.com/mrpawan-gupta/depserve/server"

func Run() {
	apiServer := server.NewAPIServer()
	apiServer.InitRouter()
	apiServer.InitMiddleWare()
	//apiServer.InitServer()
	//apiServer.InitDomain()
	//apiServer.RunServer()
}
