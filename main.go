package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env File")
	}

	remote, err := url.Parse(os.Getenv("TARGET_HOST"))
	if err != nil {
		panic(err)
	}

	log.Printf("Forwarding Target : %s://%s", remote.Scheme, remote.Host)

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ModifyResponse = corsHeaderModify

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	http.HandleFunc("/", handler(proxy))

	err = http.ListenAndServeTLS(":"+os.Getenv("LOCAL_PORT"), os.Getenv("SSL_CERT_PATH"), os.Getenv("SSL_KEY_PATH"), nil)
	if err != nil {
		panic(err)
	}
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, r *http.Request) {
		//  If, current request is pre-flight
		if r.Method == "OPTIONS" {
			resp.Header().Set("Access-Control-Allow-Origin", os.Getenv("ACCESS_CONTROL_ALLOWS_ORIGIN"))
			resp.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, content-type")
			resp.Header().Set("Access-Control-Allow-Methods", "*")
			resp.Header().Set("Access-Control-Expose-Headers", "Set-Cookie, Access-Control-Allow-Origin, Access-Control-Allow-Methods, Access-Control-Allow-Credential, Authorization")
			resp.Header().Set("Vary", "Origin")
			resp.Header().Set("Vary", "Access-Control-Request-Method")
			resp.Header().Set("Vary", "Access-Control-Request-Headers")
			resp.Header().Set("Access-Control-Allow-Credentials", "true")
			return
		} else {
			log.Printf("%s %s -> Cros_Proxy -> %s", r.Method, r.Host+r.URL.RequestURI(), os.Getenv("TARGET_HOST")+r.URL.RequestURI())
			p.ServeHTTP(resp, r)
		}
	}
}

func corsHeaderModify(resp *http.Response) error {
	// Set Basic Cors related header
	resp.Header.Set("Access-Control-Allow-Origin", os.Getenv("ACCESS_CONTROL_ALLOWS_ORIGIN"))
	resp.Header.Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, content-type")
	resp.Header.Set("Access-Control-Allow-Methods", "*")
	resp.Header.Set("Access-Control-Expose-Headers", "Set-Cookie, Access-Control-Allow-Origin, Access-Control-Allow-Methods, Access-Control-Allow-Credential, Authorization")
	resp.Header.Set("Vary", "Origin")
	resp.Header.Set("Vary", "Access-Control-Request-Method")
	resp.Header.Set("Vary", "Access-Control-Request-Headers")
	resp.Header.Set("Access-Control-Allow-Credentials", "true")

	// Parsing cookie in header
	for _, value := range strings.Split(resp.Header.Get("Set-Cookie"), ";") {
		// If remove the domain value, the client host information is automatically set to the domain value by the browser.
		if strings.Contains(value, "Domain=") {
			var newCookie = strings.Replace(resp.Header.Get("Set-Cookie"), value, "", 1)
			resp.Header.Set("Set-Cookie", newCookie)
		}
	}
	return nil
}