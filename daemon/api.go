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

	err := c.d.GetData()
	if err != nil {
		return http.StatusInternalServerError
	}
}

