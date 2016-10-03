package daemon

import (
	"net/http"
	"fmt"
)

var api = []*Command{
	jsonCmd,
}

var (
	jsonCmd = &Command{
		Path: "/",
		GET: getJson,
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

