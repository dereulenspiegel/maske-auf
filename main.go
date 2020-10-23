package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
	"github.com/sirupsen/logrus"
)

var (
	updateURL      = flag.String("url", "", "Specify the URL to update the GeoJSON from")
	updateInterval = flag.Duration("interval", time.Minute*5, "The duration after the GeoJSON should be updated")
	httpAddr       = flag.String("addr", ":8080", "The address to listen on")
)

var (
	Version = "undefined"
	Commit  = "undefined"
)

func main() {
	flag.Parse()
	logger := logrus.WithFields(logrus.Fields{
		"version": Version,
		"commit":  Commit,
	})

	logger.Info("Starting")

	killSignal := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(killSignal, os.Interrupt, os.Kill)

	go func() {
		cache := newGeoCache(*updateURL, *updateInterval)

		r := mux.NewRouter()

		r.Path("/maske").Methods(http.MethodGet).HandlerFunc(handleGet(cache))
		r.PathPrefix("/").Methods(http.MethodGet).Handler(frontendHandler())

		server := &http.Server{
			Handler:      r,
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
			Addr:         *httpAddr,
		}

		if err := server.ListenAndServe(); err != nil {
			logrus.WithError(err).Fatal("HTTP server failed")
			errChan <- err
		}
	}()
	select {
	case <-killSignal:
		os.Exit(0)
	case err := <-errChan:
		logger.WithError(err).Fatal("Something failed")
	}
}

func maskNeeded(cache *geoCache, lat, lon float64) (bool, string, error) {
	point := orb.Point{lon, lat}
	for _, feature := range cache.getCollection().Features {
		multiPoly, isMulti := feature.Geometry.(orb.MultiPolygon)
		if isMulti {
			if planar.MultiPolygonContains(multiPoly, point) {
				return true, feature.Properties.MustString("Name", "unbekannt"), nil
			}
		} else {
			polygon, isPoly := feature.Geometry.(orb.Polygon)
			if isPoly {
				if planar.PolygonContains(polygon, point) {
					return true, feature.Properties.MustString("Name", "unbekannt"), nil
				}
			} else {
				return true, "", errors.New("Unhandled geometry type")
			}
		}
	}
	return false, "", nil
}
