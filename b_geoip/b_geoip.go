package b_geoip

import (
	"github.com/djghostghost/go-basic/geoip"
)

func init() {
	geoip.Init()
}

func GetCountryRecord(ip string) (*GeoRecorder, error) {

	result, err := geoip.GetCountryRecord(ip)
	if err != nil {
		return nil, err
	}
	return &GeoRecorder{
		IP:      ip,
		Country: result.Country,
		ISOCode: result.ISOCode,
	}, nil
}

func GetCityRecord(ip string) (*GeoRecorder, error) {
	result, err := geoip.GetCityRecord(ip)
	if err != nil {
		return nil, err
	}
	return &GeoRecorder{
		IP:       ip,
		City:     result.City,
		CityCN:   result.CityCN,
		Country:  result.Country,
		ISOCode:  result.ISOCode,
		TimeZone: result.TimeZone,
	}, nil
}
