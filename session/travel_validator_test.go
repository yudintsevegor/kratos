package session

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_travelValidator_IsValid(t *testing.T) {
	type fields struct {
		maxAllowedTravelSpeed float64
	}

	type args struct {
		logins []CoordinatesWithLoginTime
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "it returns TRUE for one login only",
			fields: fields{
				maxAllowedTravelSpeed: 10,
			},
			args: args{
				logins: []CoordinatesWithLoginTime{
					{
						Longitude: 37.618423,
						Latitude:  55.751244,
						LoginTime: time.Now().UTC(),
					},
				},
			},
			want: true,
		},
		{
			name: "it returns TRUE for logins ~10km apart with a max speed 10 km/h and 2h time difference",
			fields: fields{
				maxAllowedTravelSpeed: 10,
			},
			args: args{
				logins: []CoordinatesWithLoginTime{
					{
						Longitude: 37.618423,
						Latitude:  55.751244,
						LoginTime: time.Date(2024, time.March, 1, 14, 00, 00, 00, time.UTC),
					},
					{
						Longitude: 37.45912124487777,
						Latitude:  55.77524387014726,
						LoginTime: time.Date(2024, time.March, 1, 16, 00, 00, 00, time.UTC),
					},
				},
			},
			want: true,
		},
		{
			name: "it returns TRUE for logins ~10km apart with a max speed 10 km/h and 2h time dif and reversed orderd logins",
			fields: fields{
				maxAllowedTravelSpeed: 10,
			},
			args: args{
				logins: []CoordinatesWithLoginTime{
					{
						Longitude: 37.618423,
						Latitude:  55.751244,
						LoginTime: time.Date(2024, time.March, 1, 16, 00, 00, 00, time.UTC),
					},
					{
						Longitude: 37.45912124487777,
						Latitude:  55.77524387014726,
						LoginTime: time.Date(2024, time.March, 1, 14, 00, 00, 00, time.UTC),
					},
				},
			},
			want: true,
		},
		{
			name: "it returns TRUE for logins 0km apart with a max speed 10 km/h and 0 min time difference",
			fields: fields{
				maxAllowedTravelSpeed: 10,
			},
			args: args{
				logins: []CoordinatesWithLoginTime{
					{
						Longitude: 37.618423,
						Latitude:  55.751244,
						LoginTime: time.Date(2024, time.March, 1, 14, 00, 00, 00, time.UTC),
					},
					{
						Longitude: 37.618423,
						Latitude:  55.751244,
						LoginTime: time.Date(2024, time.March, 1, 14, 00, 00, 00, time.UTC),
					},
				},
			},
			want: true,
		},
		{
			name: "it returns FALSE for logins 10km apart with a max speed 10 km/h and 0 min time difference",
			fields: fields{
				maxAllowedTravelSpeed: 10,
			},
			args: args{
				logins: []CoordinatesWithLoginTime{
					{
						Longitude: 37.618423,
						Latitude:  55.751244,
						LoginTime: time.Date(2024, time.March, 1, 14, 00, 00, 00, time.UTC),
					},
					{
						Longitude: 37.45912124487777,
						Latitude:  55.77524387014726,
						LoginTime: time.Date(2024, time.March, 1, 14, 00, 00, 00, time.UTC),
					},
				},
			},
			want: false,
		},
		{
			name: "it returns FALSE for logins ~10km apart with a max speed 10 km/h and 30 min time difference",
			fields: fields{
				maxAllowedTravelSpeed: 10,
			},
			args: args{
				logins: []CoordinatesWithLoginTime{
					{
						Longitude: 37.618423,
						Latitude:  55.751244,
						LoginTime: time.Date(2024, time.March, 1, 14, 00, 00, 00, time.UTC),
					},
					{
						Longitude: 37.45912124487777,
						Latitude:  55.77524387014726,
						LoginTime: time.Date(2024, time.March, 1, 14, 30, 00, 00, time.UTC),
					},
				},
			},
			want: false,
		},
		{
			name: "it returns FALSE for multiple valid logins and the one which is made from a 1960km apart in less than 2h",
			fields: fields{
				maxAllowedTravelSpeed: 10,
			},
			args: args{
				logins: []CoordinatesWithLoginTime{
					{
						Longitude: 37.618423,
						Latitude:  55.751244,
						LoginTime: time.Date(2024, time.March, 1, 14, 00, 00, 00, time.UTC),
					},
					{
						Longitude: 37.616376252227084,
						Latitude:  55.756094714362014,
						LoginTime: time.Date(2024, time.March, 1, 14, 10, 00, 00, time.UTC),
					},
					{
						Longitude: 37.605474880975905,
						Latitude:  55.764490414061235,
						LoginTime: time.Date(2024, time.March, 1, 14, 20, 00, 00, time.UTC),
					},
					{
						Longitude: 37.59620641906158,
						Latitude:  55.76984009358484,
						LoginTime: time.Date(2024, time.March, 1, 14, 30, 00, 00, time.UTC),
					},
					{
						Longitude: 37.584734757147196,
						Latitude:  55.77641564850109,
						LoginTime: time.Date(2024, time.March, 1, 14, 35, 00, 00, time.UTC),
					},
					{
						Longitude: 11.578439526689374,
						Latitude:  48.13609197557769,
						LoginTime: time.Date(2024, time.March, 1, 15, 22, 00, 00, time.UTC),
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("case=%s", tt.name), func(t *testing.T) {
			tv := NewTravelValidator(
				WithMaxAllowedSpeed(tt.fields.maxAllowedTravelSpeed),
			)

			assert.Equalf(t, tt.want, tv.IsValid(tt.args.logins), "IsValid(%v)", tt.args.logins)
		})
	}
}
