package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type flags struct {
	listenAddress string
	listenPort    int

	userHeader     string
	forwardHeaders map[string]bool
	fwdfwd         string
}

type controller struct {
	healthy int64
	logger  *log.Logger
}

var opts = &flags{}

func main() {

	// Load command line flags
	flag.StringVar(&opts.listenAddress, "listen-address", "0.0.0.0", "HTTP listen address")
	flag.IntVar(&opts.listenPort, "listen-port", 8000, "HTTP listen port")
	flag.StringVar(&opts.userHeader, "auth-user-header", "X-Forwarded-User", "Request header name that will hold authenticated users info")
	fwdFlag := flag.String("forward-headers", "X-Forwarded-User", "Which headers should be kept/forwarded (comma-separated)")
	flag.Parse()

	// convert comma-separated flag to a map of [header key, true]
	fwdHeaders := strings.Split(*fwdFlag, ",")
	opts.forwardHeaders = make(map[string]bool)
	for i := 0; i < len(fwdHeaders); i++ {
		opts.forwardHeaders[fwdHeaders[i]] = true
	}

	serve()
}

func serve() {

	serveAddress := opts.listenAddress + ":" + strconv.Itoa(opts.listenPort)
	logger := log.New(os.Stdout, "", log.LstdFlags)
	c := &controller{logger: logger}

	logger.Printf("--- starting to listen on %s ---", serveAddress)

	http.HandleFunc("/", c.handler)
	http.HandleFunc("/health", c.healthCheck)

	atomic.StoreInt64(&c.healthy, time.Now().UnixNano())
	//set healthcheck down with: atomic.StoreInt64(&c.healthy, 0)

	log.Fatal(http.ListenAndServe(serveAddress, nil))

}

func (c *controller) handler(w http.ResponseWriter, r *http.Request) {

	hs := make(map[string][]string)

	if opts.forwardHeaders["Host"] {
		hs["Host"] = []string{r.Host}
	}
	for k, v := range r.Header {
		if opts.forwardHeaders[k] {
			hs[k] = v
		}
	}
	users, hasHeader := r.Header[opts.userHeader]

	if !hasHeader {
		c.reject(w)
	} else {
		c.doAuth(w, users, hs, r.Host, r.RemoteAddr)
	}

}

func (c *controller) doAuth(w http.ResponseWriter, users []string, respHeaders map[string][]string, host string, remote string) {
	c.logger.Printf(`msg="incoming request", host=%q, from=%q, users=%s`, host, remote, users)

	// prepare headers
	for k, v := range respHeaders {
		for _, header := range v {
			c.logger.Printf(`msg="setting header", key=%q, value=%q`, k, header)
			w.Header().Add(k, header)
		}
	}
	// TODO: check if any (should be just one, at most) user has a service account and map it to a ServiceAccount bearer token

	// if no service account, accept auth but do not enrich with bearer token (i.e. skip)
	c.skip(w)
}

func (c *controller) reject(w http.ResponseWriter) {
	c.logger.Print(`msg="request blocked"`)
	w.WriteHeader(http.StatusUnauthorized)
}

func (c *controller) skip(w http.ResponseWriter) {
	c.logger.Print(`msg="skipping"`)
	fmt.Fprint(w, "skip")
}

func (c *controller) healthCheck(w http.ResponseWriter, req *http.Request) {
	if h := atomic.LoadInt64(&c.healthy); h == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		fmt.Fprintf(w, "uptime: %s\n", time.Since(time.Unix(0, h)))
	}
}
