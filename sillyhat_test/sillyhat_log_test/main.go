package main

import log "sillyhat-golang-tool/sillyhat_log/logrus"


func main() {
	log.Println("start")
	log.Debug("Debug")
	log.Info("info")
	log.Error("Error")
	log.Panic("Panic")
	log.Fatal("Fatal")
	log.Println("end")
}
