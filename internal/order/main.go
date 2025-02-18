package main

import (
	"github.com/Hana-bii/gorder-v2/common/config"
	"github.com/spf13/viper"
	"log"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	//log.Println("Listening :8082")
	//mux := http.NewServeMux()
	//mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	//	log.Println("%v", r.RequestURI)
	//	_, _ = io.WriteString(w, "pong")
	//})
	//if err := http.ListenAndServe(":8082", mux); err != nil {
	//	log.Fatal(err)
	//}
	log.Println("%v", viper.Get("order"))
}
