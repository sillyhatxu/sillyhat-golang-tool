package main

import (
	"sillyhat-golang-tool/sillyhat_enviroment"
	"log"
)

func init() {
	client := sillyhat_enviroment.NewEnviromentClient("dt","testMOduleName","sillyhat",0600)
	client.InitialEnviromentConfig()
}

func main() {
	client := sillyhat_enviroment.NewEnviromentClient("dt","testMOduleName","sillyhat",0600)
	log.Println(client.Get(""))

}
