package daemon

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	//"time"

	"github.com/gorilla/mux"
)

type Response interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type ResponseFunc func(*Command, *http.Request) Response

type Command struct {
	Path string
	GET    ResponseFunc
	PUT    ResponseFunc
	POST   ResponseFunc
	DELETE ResponseFunc

	d *Daemon
}

type DaemonConfig struct {
	DataURL url.URL `json: "DataURL"`
	ListenAddress string `json: "ListenAddress"`
}
type Daemon struct {
	Listener net.Listener
	router  *mux.Router
	Config	*DaemonConfig
}

// func logIt(handler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		t0 := time.Now()
// 		handler.ServeHTTP(w, r)
// 		t := time.Now().Sub(t0)
// 		url := r.URL.String()
// 		logger.Debugf("%s %s %s %s", r.RemoteAddr, r.Method, r.URL, t)
// 	})
// }

func (d *Daemon) loadConfig() error {
	f, err := os.Open("config.json")
	if err != nil {
		return fmt.Errorf("unable to open config file: ", err)
	}
	decoder := json.NewDecoder(file)
	config := DaemonConfig{}

	err := decoder.Decode(&config)
	if err != nil {
		return fmt.Errorf("unable to decode json: ", err)
	}

	d.Config = &config

	return nil
}

// register listener, initialize routes
func (d *Daemon) Init() error {
	err := d.loadConfig()
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", d.Config.ListenAddress)
    if err != nil {
		return fmt.Errorf("failed to initialize listener")
	}

	d.Listener = listener;
	d.addRoutes()

	return nil
}

func (d *Daemon) addRoutes() {
	d.Router = mux.NewRouter()

	for _, c := range api {
		c.d = d
		d.router.Handle(c.Path, c).Name(c.Path)
	}

}

// start the daemon
func (d *Daemon) Start() {
	http.Serve(d.Listener, d.router)
}

func (d *Daemon) GetData() error {
}

