// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package session

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func bearerTokenFromRequest(r *http.Request) (string, bool) {
	parts := strings.Split(r.Header.Get("Authorization"), " ")

	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1], true
	}

	return "", false
}

const (
	CloudFlareLatitudeHeader = "Cf-Iplatitude"
	CloudFlareLongitude      = "Cf-Iplongitude"
)

type coordinates struct {
	long float64
	lat  float64
}

// getCoordinatesFromRequest fetches
// a longitude (from a header cf-iplongitude) and a latitude (from a header cf-iplatitude).
// to follow the convention, it's written with capital letters.
// source: https://developers.cloudflare.com/rules/transform/managed-transforms/reference/#add-visitor-location-headers
func getCoordinatesFromRequest(r *http.Request) (coordinates, error) {
	long, err := strconv.ParseFloat(r.Header.Get(CloudFlareLongitude), 10)
	if err != nil {
		return coordinates{}, fmt.Errorf("parsing %q header: %v", CloudFlareLongitude, err)
	}

	lat, err := strconv.ParseFloat(r.Header.Get(CloudFlareLatitudeHeader), 10)
	if err != nil {
		return coordinates{}, fmt.Errorf("parsing %q header: %v", CloudFlareLatitudeHeader, err)
	}

	return coordinates{
		long: long,
		lat:  lat,
	}, nil
}

func mapDeviceInfoToCoordinates(devices []Device) []CoordinatesWithLoginTime {
	out := make([]CoordinatesWithLoginTime, 0, len(devices))

	for _, device := range devices {
		// After realising the feature, we might have situations, that not all sessions will have coordinates.
		// To prevent a "bad user experience", we would like to skip "old" saved info about sessions.
		if device.Coordinates.Longitude == 0 && device.Coordinates.Latitude == 0 {
			continue
		}

		out = append(out, CoordinatesWithLoginTime{
			Longitude: device.Coordinates.Longitude,
			Latitude:  device.Coordinates.Latitude,
			LoginTime: device.CreatedAt,
		})
	}

	return out
}
