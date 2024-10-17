package server

//
//import (
//	"fmt"
//	"github.com/go-chi/chi"
//	"github.com/jwalton/gchalk"
//	"golang.org/x/mod/modfile"
//	"math"
//	"net/http"
//	"os"
//	"path/filepath"
//	"reflect"
//	"runtime"
//	"strings"
//)
//
//func (server *APIServer) InitDomain() {
//	server.initVersion()
//	server.initAPIRoute()
//	//server.PrintAllRegisteredRoutes()
//}
//
//func (server *APIServer) initVersion() {
//	versionRouter := chi.NewRouter()
//	//versionRouter.Route("/v1/version", func(router chi.Router) {
//	//	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
//	//		utils.CreateJsonResponse(w, http.StatusOK, utils.ApiResponse{
//	//			Data:    server.Version,
//	//			Success: true,
//	//			Status:  http.StatusOK,
//	//		})
//	//	})
//	//})
//	server.Router.Mount("/api/", versionRouter)
//}
//func (server *APIServer) initAPIRoute() {
//	//api.InitAPIRoute(server.Database(), server.Router)
//}
//
//func strPad(input string, padLength int, padString string, padType string) string {
//	var output string
//
//	inputLength := len(input)
//	padStringLength := len(padString)
//
//	if inputLength >= padLength {
//		return input
//	}
//
//	repeat := math.Ceil(float64(1) + (float64(padLength-padStringLength))/float64(padStringLength))
//
//	switch padType {
//	case "RIGHT":
//		output = input + strings.Repeat(padString, int(repeat))
//		output = output[:padLength]
//	case "LEFT":
//		output = strings.Repeat(padString, int(repeat)) + input
//		output = output[len(output)-padLength:]
//	case "BOTH":
//		length := (float64(padLength - inputLength)) / float64(2)
//		repeat = math.Ceil(length / float64(padStringLength))
//		output = strings.Repeat(padString, int(repeat))[:int(math.Floor(float64(length)))] + input + strings.Repeat(padString, int(repeat))[:int(math.Ceil(float64(length)))]
//	}
//
//	return output
//}
//
//func getHandler(projectName string, handler http.Handler) (funcName string) {
//	// https://github.com/go-chi/chi/issues/424
//	funcName = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
//	base := filepath.Base(funcName)
//
//	nameSplit := strings.Split(funcName, "")
//	names := nameSplit[len(projectName):]
//	path := strings.Join(names, "")
//
//	pathSplit := strings.Split(path, "/")
//	path = strings.Join(pathSplit[:len(pathSplit)-1], "/")
//
//	sFull := strings.Split(base, ".")
//	s := sFull[len(sFull)-1:]
//
//	s = strings.Split(s[0], "")
//	if len(s) <= 4 && len(sFull) >= 3 {
//		s = sFull[len(sFull)-3 : len(sFull)-2]
//		return "@" + gchalk.Blue(strings.Join(s, ""))
//	}
//	s = s[:len(s)-3]
//	funcName = strings.Join(s, "")
//
//	return path + "@" + gchalk.Blue(funcName)
//}
//
//func getModName() string {
//	goModBytes, err := os.ReadFile("go.mod")
//	if err != nil {
//		os.Exit(0)
//	}
//	return modfile.ModulePath(goModBytes)
//}
//
//func (server *APIServer) PrintAllRegisteredRoutes(exceptions ...string) {
//	exceptions = append(exceptions, "/swagger")
//
//	walkFunc := func(method string, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
//
//		for _, val := range exceptions {
//			if strings.HasPrefix(path, val) {
//				return nil
//			}
//		}
//
//		switch method {
//		case "GET":
//			fmt.Printf("%s", gchalk.Green(fmt.Sprintf("%-8s", method)))
//		case "POST", "PUT", "PATCH":
//			fmt.Printf("%s", gchalk.Yellow(fmt.Sprintf("%-8s", method)))
//		case "DELETE":
//			fmt.Printf("%s", gchalk.Red(fmt.Sprintf("%-8s", method)))
//		default:
//			fmt.Printf("%s", gchalk.White(fmt.Sprintf("%-8s", method)))
//		}
//		fmt.Printf("%s", strPad(path, 25, "-", "RIGHT"))
//		fmt.Printf("%s\n", strPad(getHandler(getModName(), handler), 60, "-", "LEFT"))
//
//		return nil
//	}
//	if err := chi.Walk(server.Router, walkFunc); err != nil {
//		fmt.Print(err)
//	}
//
//}
