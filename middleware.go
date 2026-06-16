package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func (app *application) fetchWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		app.writeJSON(w, r, http.StatusBadRequest, envelope{
			"error": "empty city parameter, city required",
		})
		return
	}

	val, err := app.rdb.Get(app.ctx, city).Result()

	var apiResp APIResponse

	if err == redis.Nil {
		//cache miss, fetch from visualcrossing
		url := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?include=current&unitGroup=metric&contentType=json&key=%s",
			city, app.config.APIKey)

		resp, err := http.Get(url)
		if err != nil {
			app.writeJSON(w, r, http.StatusInternalServerError, envelope{"error": err})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			app.writeJSON(w, r, http.StatusInternalServerError, envelope{"error": fmt.Sprintf("weather API returned status: %d", resp.StatusCode)})
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			app.writeJSON(w, r, http.StatusInternalServerError, envelope{"error": err})
			return
		}

		if err = app.readJSON(body, &apiResp); err != nil {
			app.writeJSON(w, r, http.StatusInternalServerError, envelope{"error": "failed to parse weather response"})
		}

		if err := app.rdb.Set(app.ctx, city, body, 5*time.Minute).Err(); err != nil {
			log.Printf("redis set error: %s", err)
		}

		w.Header().Set("Content-Type", "application/json")
		app.writeJSON(w, r, http.StatusOK, envelope{
			"city":    city,
			"weather": apiResp.CurrentConditions,
			"cached":  false,
		})

	} else if err != nil {
		// something went wrong with Redis
		app.writeJSON(w, r, http.StatusInternalServerError, envelope{
			"error": err.Error(),
		})
		return
	} else {
		//cache hit
		if err = app.readJSON([]byte(val), &apiResp); err != nil {
			app.writeJSON(w, r, http.StatusInternalServerError, envelope{
				"error": "failed to parse cached data",
			})
		}

		w.Header().Set("Content-Type", "application/json")
		app.writeJSON(w, r, http.StatusOK, envelope{
			"city":    city,
			"weather": apiResp.CurrentConditions,
			"cached":  true,
		})
		return

	}

}
