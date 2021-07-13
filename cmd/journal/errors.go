package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

type errorResponse struct {
	Error interface{} `json:"error"`
}

func (a *App) writeError(w http.ResponseWriter, status int, message interface{}) {
	payload := errorResponse{Error: message}
	err := a.writeJSON(w, status, payload)
	if err != nil {
		w.WriteHeader(500)
	}
}

func (a *App) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := "the server encountered a problem and could not process your request"
	a.errorLog.Println(err)
	a.errorLog.Println(string(debug.Stack()))
	a.writeError(w, http.StatusInternalServerError, message)
}

func (a *App) notFoundResponse(w http.ResponseWriter) {
	message := "the requested resource could not be found"
	a.writeError(w, http.StatusNotFound, message)
}

func (a *App) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	a.writeError(w, http.StatusMethodNotAllowed, message)
}

func (a *App) badRequestResponse(w http.ResponseWriter, err error) {
	a.writeError(w, http.StatusBadRequest, err.Error())
}

func (a *App) unauthorizedResponse(w http.ResponseWriter) {
	message := "authorization is required for this method or the provided authorization was invalid"
	a.writeError(w, http.StatusUnauthorized, message)
}
