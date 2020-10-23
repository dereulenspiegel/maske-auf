package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

var (
	invalidGetParam = &apiError{
		Status: 400,
		Type:   "invalid_parameter",
		Title:  "Invalid query parameter",
		Detail: "A query parameter is either missing or invalid",
	}
)

func handleGet(cache *geoCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithFields(logrus.Fields{
			"user_agent":  r.UserAgent(),
			"remote_addr": r.RemoteAddr,
		})
		latStr := r.URL.Query().Get("lat")
		lonStr := r.URL.Query().Get("lon")

		logger.WithFields(logrus.Fields{
			"lat": latStr,
			"lon": lonStr,
		})

		if latStr == "" || lonStr == "" {
			invalidGetParam.Write(w)
			return
		}

		lat, err := strconv.ParseFloat(latStr, 64)
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			logger.WithError(err).Error("Failed to parse latitude and longitude from query parameters")
			invalidGetParam.Write(w)
			return
		}
		mask, zone, err := maskNeeded(cache, lat, lon)

		json.NewEncoder(w).Encode(&response{
			MaskNeeded: mask,
			ZoneName:   zone,
		})
	}
}

type response struct {
	MaskNeeded bool   `json:"mask_needed"`
	ZoneName   string `json:"zone_name,omitempty"`
}

type apiError struct {
	Status int    `json:"-"`
	Type   string `json:"type"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (a *apiError) Write(w http.ResponseWriter) {
	w.WriteHeader(a.Status)
	json.NewEncoder(w).Encode(a)
}
