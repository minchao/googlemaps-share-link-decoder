package googlemaps_share_link_decoder

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Service interface {
	Decode(*Request) (*Response, error)
}

type ShareLinkService struct{}

// Decode
func (ShareLinkService) Decode(req *Request) (*Response, error) {
	url := req.URL
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var latLngArr []float64
	var address, phone, name interface{}

	// Determine the url is location or place
	reg := regexp.MustCompile(`maps\/(@|search\/)(\-?\d+(\.\d+)?),(\-?\d+(\.\d+)?)`)
	tmp := reg.FindAllStringSubmatch(r.Request.URL.String(), -1)
	if len(tmp) == 1 {
		// location
		lat, _ := strconv.ParseFloat(tmp[0][2], 64)
		lng, _ := strconv.ParseFloat(tmp[0][4], 64)

		latLngArr = []float64{lat, lng}
	} else {
		// place
		rawBody, _ := ioutil.ReadAll(r.Body)
		body := strings.Replace(string(rawBody), "\n", "", -1)
		reg := regexp.MustCompile(`cacheResponse\((.*?)\);`)
		tmp := reg.FindAllStringSubmatch(body, -1)
		if len(tmp) == 0 {
			return nil, errors.New("requested JSON data not found")
		}

		// Parse JSON data
		var data, place, locations []json.RawMessage
		err = json.Unmarshal([]byte(tmp[0][1]), &data)
		if err != nil || len(data) < 55 {
			return nil, errors.New("wrong JSON data format")
		}
		if string(data[8]) == "null" {
			return nil, errors.New("place data not found")
		}
		json.Unmarshal(data[8], &place)
		if len(place) < 14 {
			return nil, errors.New("place data not found")
		}
		if string(place[0]) == "null" {
			return nil, errors.New("location data not found")
		}
		json.Unmarshal(place[0], &locations)
		err = json.Unmarshal(locations[2], &latLngArr)
		if err != nil {
			return nil, errors.New("incorrect coordinate format")
		}
		// address
		json.Unmarshal(place[13], &address)
		// phone
		json.Unmarshal(place[7], &phone)
		// name
		json.Unmarshal(place[1], &name)
	}

	result := Response{
		FormattedAddress:     address,
		FormattedPhoneNumber: phone,
		Geometry:             Geometry{Location{latLngArr[0], latLngArr[1]}},
		Name:                 name,
	}

	return &result, nil
}

type Request struct {
	URL string `json:"url"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Geometry struct {
	Location Location `json:"location"`
}

type Response struct {
	FormattedAddress     interface{} `json:"formatted_address"`
	FormattedPhoneNumber interface{} `json:"formatted_phone_number"`
	Geometry             Geometry    `json:"geometry"`
	Name                 interface{} `json:"name"`
}
