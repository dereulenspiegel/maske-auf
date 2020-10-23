package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaskeAuf(t *testing.T) {
	cache := newGeoCache("https://geoweb1.digistadtdo.de/doris_gdi/gisserver3/Fussgaengerzone/gisserver?SERVICE=WFS&VERSION=1.1.0&REQUEST=GetFeature&TYPENAME=Fussgaengerzone&OUTPUTFORMAT=GeoJSON", time.Hour)

	validPoint := []float64{51.514333, 7.461445}

	maske, name, err := maskNeeded(cache, validPoint[0], validPoint[1])
	require.NoError(t, err)
	assert.True(t, maske)
	assert.Equal(t, `MNS Westen-/ Ostenhellweg`, name)

	invalidPoint := []float64{51.490846, 7.7119190}

	maske, name, err = maskNeeded(cache, invalidPoint[0], invalidPoint[1])
	require.NoError(t, err)
	assert.False(t, maske)
	assert.Equal(t, "", name)
}
