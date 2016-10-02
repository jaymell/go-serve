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

// this may differ across applications, obviously
type Data struct {
	XForwardedFor string `bson: "x-forwarded-for" json: "x-forwarded-for"`
}

type Response interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type ResponseFunc func(*Command, *http.Request) Response

type resp struct {
	Status int `json: "Status"`
	Result interface{} `json: "Result"`
}

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
	CollectionName string `json: "CollectionName"`
}

type Daemon struct {
	Listener net.Listener
	router  *mux.Router
	Config	*DaemonConfig
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
		rsp = rspf(c, r, user)
	}

	rsp.ServeHTTP(w, r)  
}

func (r *resp) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp{
			Status: r.Status,
			Result: r.Result,
		})
}

// called by Command.ServeHTTP -- the response header/body actually gets written:
func (r *resp) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	bs, err := r.MarshalJSON()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	w.Write(bs)
}

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

func (d *Daemon) GetData() (interface{}, error) {
	switch d.Config.DataURL.scheme {
	case "mongo":
		return getMongoData(d.Config.DataURL, d.Config.CollectionName)
	case "http":
		return nil, fmt.Errorf("Not implemented yet")
	default:
		return nil, fmt.Errorf("Unrecognized scheme")
	}	
}

func getMongoData(url URL, col string) (interface{}, error) {
	dialInfo, err := mgo.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("unable to parse db url: ", err)
	}
	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: ", err)
	}
	c := session.DB(dialinfo.Database).C(col)
	results := []Data{}
	c.Find(nil).All(&results)
	// FIXME -- don't always return nil for err, probably:
	return results, nil	
}

func SyncResponse(result interface{}) Response {
	if err, ok := result.(error); ok {
		return &resp{
			Status: http.StatusInternalServerError,
			Result: nil,
		}
	}
	return &resp{
		Status: http.StatusOK,
		Result: result,
	}	
}