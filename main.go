package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
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

// Binding type to be loaded in the configuration
// (must be exported for yaml unmarshalling)
type Binding struct {
	User           string `yaml:"user"`
	ServiceAccount string `yaml:"service-account"`
	UserAccount    string `yaml:"account"`
}

// HeadersConfig type to be loaded in the configuration
// (must be exported for yaml unmarshalling)
type HeadersConfig struct {
	AuthUser        string `yaml:"auth-user"`
	BearerTokenPath string `yaml:"impersonator-token"`
	BearerToken     string
	ToForward       []string `yaml:"to-forward"`
	toForwardOnce   sync.Once
	toForwardLookup map[string]bool // for performance reasons
}

// ShouldForward returns true if the provided header should be kept/forwarded
// this method is memoized over ToForward slices
func (hc *HeadersConfig) ShouldForward(header string) bool {
	hc.toForwardOnce.Do(func() {
		hc.toForwardLookup = make(map[string]bool)
		for _, v := range hc.ToForward {
			hc.toForwardLookup[v] = true
		}
	})
	return hc.toForwardLookup[header]
}

// ServerConfig type to be loaded in the configuration
// (must be exported for yaml unmarshalling)
type ServerConfig struct {
	Host    string        `yaml:"host" envconfig:"SERVER_HOST"`
	Port    int           `yaml:"port" envconfig:"SERVER_PORT"`
	Headers HeadersConfig `yaml:"headers"`
}

// Config type to be loaded in the configuration
// (must be exported for yaml unmarshalling)
type Config struct {
	Server   ServerConfig `yaml:"server"`
	Bindings []Binding    `yaml:"bindings"`
}

var opts = &flags{}
var configs = &Config{}

var cfg Config

func main() {

	configFile := flag.String("config", "config.yaml", "configuration file path")
	flag.Parse()

	readConfig(*configFile, &cfg)
	readEnv(&cfg)

	log.Printf("cfg : %+v\n", &cfg)
	serve(&cfg)
}

func readConfig(path string, cfg *Config) {
	file, _ := os.Open(path)
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err := decoder.Decode(&cfg)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	readToken(cfg.Server.Headers.BearerTokenPath, cfg)
}

func readToken(path string, cfg *Config) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	cfg.Server.Headers.BearerToken = string(content)
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func serve(cfg *Config) {

	serveAddress := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
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
	hc := &cfg.Server.Headers

	if hc.ShouldForward("Host") {
		hs["Host"] = []string{r.Host}
	}
	for k, v := range r.Header {
		if hc.ShouldForward(k) {
			hs[k] = v
		}
	}
	users, hasHeader := r.Header[hc.AuthUser]

	if !hasHeader {
		// TODO: check if this is the desired behaviour (or if should be configured)
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

	// check if any (should be just one, at most) user has a binding, and if so impersonate it
	for _, user := range users {
		for _, b := range cfg.Bindings {
			if b.User != "" && b.User == user {
				c.impersonate(b, w)
				return
			}
		}
	}

	// if no account to impersonate, accept auth but do not enrich with bearer token (i.e. skip)
	c.skip(w)
}

func (c *controller) impersonate(b Binding, w http.ResponseWriter) {
	c.logger.Printf(`msg="impersonating", binding="%+v"`, b)
	user := b.ServiceAccount

	if b.UserAccount != "" {
		user = b.UserAccount
	}

	if user == "" {
		c.logger.Print(`msg="cannot impersonate user without account defined, please review configs"`)
		c.skip(w)
	} else {
		w.Header().Add("Authorization", "Bearer "+cfg.Server.Headers.BearerToken)
		w.Header().Add("Impersonate-User", user)
	}
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
