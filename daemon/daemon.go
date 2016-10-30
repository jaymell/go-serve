package daemon

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
)

type Response interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type ResponseFunc func(*Command, *http.Request) Response

type Resp struct {
	Status int         `json:"Status"`
	Result interface{} `json:"Result"`
}

type Command struct {
	Path   string
	GET    ResponseFunc
	PUT    ResponseFunc
	POST   ResponseFunc
	DELETE ResponseFunc

	API API
}

type DaemonConfig struct {
	ListenAddress  string `json:"ListenAddress"`
}

type Daemon struct {
	Listener net.Listener
	Config   *DaemonConfig
	api	API
}

type API interface {
	Routes() []*Command
}

// user auth can be validated here, if implemented -- see snapd code:
func (c *Command) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var rspf ResponseFunc

	switch r.Method {
	case "GET":
		rspf = c.GET
	case "PUT":
		rspf = c.PUT
	case "POST":
		rspf = c.POST
	case "DELETE":
		rspf = c.DELETE
	}

	if rspf != nil {
		rsp := rspf(c, r)
		rsp.ServeHTTP(w, r)
	} else {
		var badResp = &Resp{
			Status: http.StatusMethodNotAllowed,
			Result: nil,
		}
		badResp.ServeHTTP(w, r)
	}
}

func (r *Resp) MarshalJSON() ([]byte, error) {
	return json.Marshal(Resp{
		Status: r.Status,
		Result: r.Result,
	})
}

// called by Command.ServeHTTP -- the Response header/body actually gets written:
func (r *Resp) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	bs, err := r.MarshalJSON()
	if err != nil {
		r.Status = http.StatusInternalServerError
		bs = nil
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	w.Write(bs)
}

func (d *Daemon) loadConfig(f *os.File) error {
	decoder := json.NewDecoder(f)
	config := DaemonConfig{}

	err := decoder.Decode(&config)
	if err != nil {
		return fmt.Errorf("unable to decode json: ", err)
	}

	d.Config = &config

	return nil
}

// register listener, initialize routes
func (d *Daemon) Init(api API, f *os.File) error {
	err := d.loadConfig(f)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", d.Config.ListenAddress)
	if err != nil {
		return fmt.Errorf("failed to initialize listener")
	}

	d.Listener = listener
	d.addRoutes(api)

	return nil
}

func (d *Daemon) addRoutes(api API) {

	routes := api.Routes()
	for _, c := range routes {
		http.Handle(c.Path, c)
	}

	// everything else assumed to be static content
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
}

// start the daemon
func (d *Daemon) Start() {
	http.Serve(d.Listener, nil)
}

func SyncResponse(result interface{}) Response {
	if _, ok := result.(error); ok {
		return &Resp{
			Status: http.StatusInternalServerError,
			Result: nil,
		}
	}
	return &Resp{
		Status: http.StatusOK,
		Result: result,
	}
}
