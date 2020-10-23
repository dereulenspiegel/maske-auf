package main

import (
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	geojson "github.com/paulmach/orb/geojson"
	"github.com/sirupsen/logrus"
)

type geoCache struct {
	currCollection *geojson.FeatureCollection

	lastUpdate     time.Time
	updateInterval time.Duration
	updateURL      string
	lock           *sync.Mutex
}

func newGeoCache(updateURL string, updateInterval time.Duration) *geoCache {
	return &geoCache{
		updateURL:      updateURL,
		updateInterval: updateInterval,
		lock:           &sync.Mutex{},
	}
}

func (g *geoCache) getCollection() *geojson.FeatureCollection {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.currCollection == nil {
		g.lock.Unlock()
		g.update()
		g.lock.Lock()
	}
	if time.Now().After(g.lastUpdate.Add(g.updateInterval)) {
		go g.update()
	}

	return g.currCollection
}

func (g *geoCache) update() {
	logger := logrus.WithField("update_url", g.updateInterval)
	req, err := http.NewRequest(http.MethodGet, g.updateURL, nil)
	if err != nil {
		logger.WithError(err).Error("Failed to create update request")
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.WithError(err).Error("Failed to get update from remote")
		return
	}
	logger.WithField("response_code", resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	newCollection, err := geojson.UnmarshalFeatureCollection(body)
	if err != nil {
		logger.WithError(err).Error("Failed to unmarshal feature collection")
		return
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	g.currCollection = newCollection
}
