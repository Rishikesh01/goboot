package goboot

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter
}

func (c *Context) JSON(status int, data any) {
	json, err := json.Marshal(&data)
	if err != nil {
		return
	}

	c.Writer.Write(json)
}

func (c *Context) BindJSON(data any) {
	defer c.Request.Body.Close()
	if err := json.NewDecoder(c.Request.Body).Decode(data); err != nil {
		return
	}
}
