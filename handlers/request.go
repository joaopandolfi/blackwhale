package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/joaopandolfi/blackwhale/configurations"
	"github.com/joaopandolfi/blackwhale/utils"
)

// --- Responses ---

// header -
func header(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", configurations.Configuration.CORS)
	w.Header().Add("Content-Type", "application/json")
}

// responseError - Private function to make response
func responseError(w http.ResponseWriter, message string) {
	var response map[string]interface{}
	response = make(map[string]interface{})
	response["message"] = message
	b, _ := json.Marshal(response)

	w.WriteHeader(500)
	w.Write(b)
}

// restResponseError - Private function to response in mode RES error
func restResponseError(w http.ResponseWriter, message string) {
	var response map[string]interface{}
	response = make(map[string]interface{})
	response["success"] = true
	response["message"] = message
	b, _ := json.Marshal(response)

	w.Write(b)
}

// RESTResponse - Make default REST API response
func RESTResponse(w http.ResponseWriter, resp interface{}) {
	var response map[string]interface{}
	response = make(map[string]interface{})
	response["success"] = true
	response["data"] = resp

	Response(w, response)
}

// Response - Make default generic response
func Response(w http.ResponseWriter, resp interface{}) {
	// set Header
	header(w)
	b, err := json.Marshal(resp)

	if err == nil {
		// Responde
		w.Write(b)
	} else {
		utils.Error("Error on convert response to JSON", err)
		ResponseError(w, "Error on convert response to JSON")
	}
}

// ResponseError - Make default generic response
func ResponseError(w http.ResponseWriter, resp interface{}) {
	// set Header
	header(w)
	b, _ := json.Marshal(resp)
	responseError(w, string(b))
}

// RESTResponseError - Make REST API default response
func RESTResponseError(w http.ResponseWriter, resp interface{}) {
	// set Header
	header(w)
	b, _ := json.Marshal(resp)
	restResponseError(w, string(b))
}

// Redirect - Redirect page
func Redirect(r *http.Request, w http.ResponseWriter, url string) {
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// --- Parameters ---

// GetVars - Return url vars
// @example /api/{key}/send
// @vars = {"key":data}
func GetVars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

// GetHeader - Return Header value stored on passed key
func GetHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
}

// GetQueryes - Return queryes values
// @example /api?key=data
func GetQueryes(r *http.Request) url.Values {
	return r.URL.Query()
}

// GetBody - Return byte body data
func GetBody(r *http.Request) ([]byte, error) {
	return ioutil.ReadAll(r.Body)
}

// GetForm - Return parsed form data
func GetForm(r *http.Request) (form url.Values, err error) {
	err = r.ParseForm()
	form = r.Form
	return
}

// DecodeForm - Decoded parsed form data on interface
func DecodeForm(dst interface{}, src map[string][]string) error {
	decoder := schema.NewDecoder()
	return decoder.Decode(dst, src)
}

// GetSession returns stored Session
// @global
// Login session keys: `logged`, `username`, `institution`, `level`, `token`
func GetSession(r *http.Request) (*sessions.Session, error) {
	return configurations.Configuration.Session.Store.Get(r, configurations.Configuration.Session.Name)
}

// GetNamedSession - Return data sored on specific session
func GetNamedSession(r *http.Request, name string) (*sessions.Session, error) {
	return configurations.Configuration.Session.Store.Get(r, name)
}
