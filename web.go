package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	env1 string = "VERSION"
)

/*init log*/
func init() {
	logFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("failed to open access.log", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lmicroseconds | log.Ldate)
}

func main() {
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/", access)
	err := http.ListenAndServe("192.168.2.6:8080", nil)
	if err != nil {
		log.Fatal("failed to initial http servers", err)
	}

}

/*healt check handler function*/
func healthz(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			io.WriteString(w, "the system is unvaible now!\n")
			log.Printf("%s  %s---- 500 -----", ClientIP(r), r.RequestURI)
		}
	}()
	w.WriteHeader(200)
	io.WriteString(w, "healthy check is good\n")
	log.Printf("%s  %s ---- 200 -----", ClientIP(r), r.RequestURI)

}

/*root handler function*/
func access(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if p := recover(); p != nil {
			w.WriteHeader(500)
			io.WriteString(w, "the system is unvaible now!\n")
			log.Printf("%s %s ---- 500 -----", ClientIP(r), r.RequestURI)

		}
	}()
	/*add response header*/
	for k, v := range r.Header {
		w.Header().Add("r_"+k, v[0])
	}
	w.Header().Add(env1, getOsVariable(env1))
	w.WriteHeader(200)
	io.WriteString(w, "request header resolve successfully ok\n")
	log.Printf("%s  %s ---- 200 -----", ClientIP(r), r.RequestURI)
}

/*get env varibale*/
func getOsVariable(name string) string {
	if len(name) == 0 {
		return ""
	}
	return os.Getenv(name)

}

/*get clinet IP*/
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
