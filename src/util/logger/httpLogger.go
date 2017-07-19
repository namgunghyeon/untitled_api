package logger

import (
	"net/http"
	"fmt"
	"log"
	"os"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello World</h1>")
}

func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.URL, r.Body)
		fmt.Println(r.RemoteAddr, r.Method, r.URL, r.Body)
		handler.ServeHTTP(w, r)
	})
}

func OpenLogFile(logfile string) {
  log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

		if err != nil {
			log.Fatal("OpenLogfile: os.OpenFile:", err)
		}

		log.SetOutput(lf)
	}
}
