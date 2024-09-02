package session

import (
	"sort"
	"time"

	geo "github.com/kellydunn/golang-geo"
)

type TravelValidator interface {
	SessionTravelValidator() Validator
}

type Validator interface {
	IsValid(logins []CoordinatesWithLoginTime) bool
}

const defaultAllowedSpeed = 10

type CoordinatesWithLoginTime struct {
	Longitude float64
	Latitude  float64

	LoginTime time.Time
}

type travelValidator struct {
	maxAllowedTravelSpeed float64 // km/h
}

func NewTravelValidator(opts ...ValidatorOptions) *travelValidator {
	cfg := &validatorOptions{
		maxAllowedTravelSpeed: defaultAllowedSpeed,
	}

	for _, o := range opts {
		o(cfg)
	}

	return &travelValidator{
		maxAllowedTravelSpeed: cfg.maxAllowedTravelSpeed,
	}
}

// IsValid validates incoming logins (coords with login time) on "impossible travels".
// Every login is compared to each other to detect "unusual travels" which are not possible
// to achieve for a given speed `maxAllowedTravelSpeed` (it's configured).
// If "impossible travel" is detected, the function returns "false", otherwise "true".
//
// NOTE: I might assume it can be simplified:
//	we could compare everything in one iteration, but it's just an assumption, so that's why I
//	leave it as it is.
func (tv *travelValidator) IsValid(logins []CoordinatesWithLoginTime) bool {
	sort.Slice(logins, func(i, j int) bool {
		return logins[i].LoginTime.After(logins[j].LoginTime)
	})

	for currentInd, currentLogin := range logins {
		for targetInd := currentInd + 1; targetInd <= len(logins)-1; targetInd++ {
			targetLogin := logins[targetInd]

			currentGeoPoint := geo.NewPoint(currentLogin.Latitude, currentLogin.Longitude)
			targetGeoPoint := geo.NewPoint(targetLogin.Latitude, targetLogin.Longitude)

			distanceKM := currentGeoPoint.GreatCircleDistance(targetGeoPoint)

			timeDiff := currentLogin.LoginTime.Sub(targetLogin.LoginTime).Seconds()

			if timeDiff == 0 {
				if distanceKM == 0 {
					continue
				} else {
					return false
				}
			}

			if distanceKM/(timeDiff/60/60) > tv.maxAllowedTravelSpeed {
				return false
			}
		}
	}

	return true
}

type ValidatorOptions func(*validatorOptions)

type validatorOptions struct {
	maxAllowedTravelSpeed float64 // km/h
}

// WithMaxAllowedSpeed passes along the max allowed speed in km/h for the distance sessions.
func WithMaxAllowedSpeed(maxAllowedTravelSpeed float64) ValidatorOptions {
	return func(opts *validatorOptions) {
		if maxAllowedTravelSpeed == 0 {
			return
		}

		opts.maxAllowedTravelSpeed = maxAllowedTravelSpeed
	}
}
