package main

// func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
// 	env := envelope{"error": message}

// 	err := app.writeJSON(w, r, status, env)
// 	if err != nil {
// 		w.WriteHeader(500)
// 		return
// 	}
// }

// func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
// 	message := "the requested resource could not be found"
// 	app.errorResponse(w, r, http.StatusNotFound, message)
// }
