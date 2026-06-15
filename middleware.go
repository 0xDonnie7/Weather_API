package main

import (
	"net/http"

	"github.com/redis/go-redis/v9"
)

func (app *application) fetchWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		app.writeJSON(w, r, http.StatusBadRequest, envelope{
			"error": "empty city parameter, city required",
		})
	}

	val, err := app.rdb.Get(app.ctx, city).Result()
	if err == redis.Nil {
		//cache miss, fetch from visualcrossing

	} else if err != nil {
		// something went wrong with Redis
		app.writeJSON(w, r, http.StatusInternalServerError, envelope{
			"error": err,
		})
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(val))
		return
	}

}
