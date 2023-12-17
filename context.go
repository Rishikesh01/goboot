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
	HttpRoutePath string
	Request       *http.Request
	Writer        http.ResponseWriter
}

func (c *Context) String(status int, data string) {
	c.Writer.WriteHeader(status)
	c.Writer.Header().Add(ContentType, TextPlain)
	c.Writer.Write([]byte(data))
}

func (c *Context) JSON(status int, data any) error {
	c.Writer.WriteHeader(status)
	c.Writer.Header().Add(ContentType, ApplicationJSON)
	if err := json.NewEncoder(c.Writer).Encode(&data); err != nil {
		return err
	}

	return nil
}

func (c *Context) BindJSON(data any) error {
	defer c.Request.Body.Close()
	if err := json.NewDecoder(c.Request.Body).Decode(data); err != nil {
		return err
	}

	return nil
}
