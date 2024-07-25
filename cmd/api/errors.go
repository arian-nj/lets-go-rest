package main

import "net/http"

// func (app *application) errRespos(w http.ResponseWriter, r *http.Request, errors map[string]string)  {

// }

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	WriteJSON(w, http.StatusUnprocessableEntity, envelope{"error": errors})
}
