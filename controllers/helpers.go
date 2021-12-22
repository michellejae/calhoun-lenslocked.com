package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
)

func parseForm(r *http.Request, dst interface{}) error {
	// parseform causing r.PostForm field to be filled with data from sign up form
	if err := r.ParseForm(); err != nil {
		return err
	}
	// create a decoder with schema package to decode form and save it into our signupform type
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	if err := dec.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}
