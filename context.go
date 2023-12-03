package goboot

import (
	"encoding/json"
	"net/http"
)

const (
	ContentType               = "Content-Type"
	TextPlain                 = "text/plain"
	ApplicationJSON           = "application/json"
	ApplicationXML            = "application/xml"
	ApplicationFormURLEncoded = "application/x-www-form-urlencoded"
	MultipartFormData         = "multipart/form-data"
	ApplicationOctetStream    = "application/octet-stream"
	ApplicationPDF            = "application/pdf"
)

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter
}

func (c *Context) String(status int, data string) {
	c.Writer.WriteHeader(status)
	c.Writer.Header().Add(ContentType, TextPlain)
	c.Writer.Write([]byte(data))
}

func (c *Context) JSON(status int, data any) {
	c.Writer.WriteHeader(status)
	c.Writer.Header().Add(ContentType, ApplicationJSON)
	if err := json.NewEncoder(c.Writer).Encode(&data); err != nil {
		return
	}
}

func (c *Context) BindJSON(data any) {
	defer c.Request.Body.Close()
	if err := json.NewDecoder(c.Request.Body).Decode(data); err != nil {
		return
	}
}
