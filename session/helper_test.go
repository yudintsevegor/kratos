// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package session

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBearerTokenFromRequest(t *testing.T) {
	for k, tc := range []struct {
		h http.Header
		t string
		f bool
	}{
		{
			h: http.Header{"Authorization": {"Bearer token"}},
			t: "token", f: true,
		},
		{
			h: http.Header{"Authorization": {"bearer token"}},
			t: "token", f: true,
		},
		{
			h: http.Header{"Authorization": {"beaRer token"}},
			t: "token", f: true,
		},
		{
			h: http.Header{"Authorization": {"BEARER token"}},
			t: "token", f: true,
		},
		{
			h: http.Header{"Authorization": {"notbearer token"}},
		},
		{
			h: http.Header{"Authorization": {"token"}},
		},
		{
			h: http.Header{"Authorization": {}},
		},
		{
			h: http.Header{},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			token, found := bearerTokenFromRequest(&http.Request{Header: tc.h})
			assert.Equal(t, tc.f, found)
			assert.Equal(t, tc.t, token)
		})
	}
}

// I would call it Test_getCoordinatesFromRequest, but it's needed to follow the code style
func TestGetCoordinatesFromRequest(t *testing.T) {
	for _, tc := range []struct {
		name        string
		h           http.Header
		want        coordinates
		errExpected bool
	}{
		{
			name:        "happy path",
			h:           http.Header{CloudFlareLongitude: {"55.751244"}, CloudFlareLatitudeHeader: {"37.618423"}},
			want:        coordinates{long: 55.751244, lat: 37.618423},
			errExpected: false,
		},
		{
			name:        "should have the error: no cf-iplongitude and no cf-iplatitude",
			h:           http.Header{},
			want:        coordinates{long: 0, lat: 0},
			errExpected: true,
		},
		{
			name:        "should have the error: no cf-iplongitude",
			h:           http.Header{CloudFlareLatitudeHeader: {"37.618423"}},
			want:        coordinates{long: 0, lat: 0},
			errExpected: true,
		},
		{
			name:        "should have the error: no cf-iplatitude",
			h:           http.Header{CloudFlareLongitude: {"55.751244"}},
			want:        coordinates{long: 0, lat: 0},
			errExpected: true,
		},
		{
			name:        "should have the error: empty value for cf-iplongitude",
			h:           http.Header{CloudFlareLongitude: {""}, CloudFlareLatitudeHeader: {"37.618423"}},
			want:        coordinates{long: 0, lat: 0},
			errExpected: true,
		},
		{
			name:        "should have the error: empty value for cf-iplatitude",
			h:           http.Header{CloudFlareLongitude: {"55.751244"}, CloudFlareLatitudeHeader: {""}},
			want:        coordinates{long: 0, lat: 0},
			errExpected: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := getCoordinatesFromRequest(&http.Request{Header: tc.h})
			assert.Equal(t, tc.want, result)
			if tc.errExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
