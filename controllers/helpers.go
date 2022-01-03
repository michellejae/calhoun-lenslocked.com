package controllers

import (
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

func parseForm(r *http.Request, dst interface{}) error {
	// parseform causing r.PostForm field to be filled with data from sign up form
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.PostForm, dst)
}

func parseURLParams(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return nil
	}
	return parseValues(r.Form, dst)
}
func parseValues(values url.Values, dst interface{}) error {
	// create a decoder with schema package to decode form and save it into our signupform type
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	if err := dec.Decode(dst, values); err != nil {
		return err
	}
	return nil
}
