package daemon

import (
	"fmt"
	"net/http"
)

var api = []*Command{
	jsonCmd,
}

var (
	jsonCmd = &Command{
		Path: "/json",
		GET:  getJson,
	}
)

func getJson(c *Command, r *http.Request) Response {
	data, err := c.d.GetData()
	if err != nil {
		fmt.Println("error retrieving data: ", err)
		return &resp{
			Status: http.StatusInternalServerError,
			Result: nil,
		}
	}

	return SyncResponse(data)
}

func getData(d *Daemon) (interface{}, error) {
        URL, err := url.Parse(d.Config.DataURL)
        if err != nil {
                return nil, fmt.Errorf("could not parse url from config: ", err)
        }
        switch URL.Scheme {
        case "mongodb":
                return getMongoData(d.Config.DataURL, d.Config.CollectionName)
        case "http":
                return nil, fmt.Errorf("Not implemented yet")
        default:
                return nil, fmt.Errorf("Unrecognized scheme")
        }
}

func getMongoData(URL string, col string) (interface{}, error) {
        dialInfo, err := mgo.ParseURL(URL)
        if err != nil {
                return nil, fmt.Errorf("unable to parse db url: ", err)
        }

        dialInfo.Timeout = time.Duration(5) * time.Second
        session, err := mgo.DialWithInfo(dialInfo)

        if err != nil {
                return nil, fmt.Errorf("unable to connect to db: ", err)
        }
        c := session.DB(dialInfo.Database).C(col)
        results := []Data{}
        c.Find(nil).All(&results)
        // FIXME -- don't always return nil for err, probably:
        return results, nil
}

