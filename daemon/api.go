package daemon

import (
	"net/http"
)

var api = []*Command{
	jsonCmd,
}

var (
	jsonCmd = &Command{
		Path: "/json",
		GET: getJson,
	}
)

func getJson(c *Command, r *http.Request) Response {
	data, err := c.d.GetData()
	if err != nil {
		return &resp{
			Status: http.StatusInternalServerError,
			Result: nil,
		}
	}

	return SyncResponse(data)
}

