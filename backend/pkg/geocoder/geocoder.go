package geocoder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	_ "strconv"
	"strings"
)

type GeoResponse struct {
	Response struct {
		GeoObjectCollection struct {
			FeatureMember []struct {
				GeoObject struct {
					Point struct {
						Pos string `json:"pos"`
					} `json:"Point"`
				} `json:"GeoObject"`
			} `json:"featureMember"`
		} `json:"GeoObjectCollection"`
	} `json:"response"`
}

func GetCoordinates(fullAddress string) (string, string, error) {
	baseURL := "https://geocode-maps.yandex.ru/1.x/"
	params := url.Values{}
	params.Set("apikey", "e0b9b1d8-b8fa-4c07-ae3a-8b524663ebf5")
	params.Set("format", "json")
	params.Set("geocode", fullAddress)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var geoResp GeoResponse
	if err := json.Unmarshal(body, &geoResp); err != nil {
		return "", "", err
	}

	if len(geoResp.Response.GeoObjectCollection.FeatureMember) == 0 {
		return "", "", fmt.Errorf("coordinates not found")
	}

	pos := geoResp.Response.GeoObjectCollection.FeatureMember[0].GeoObject.Point.Pos
	parts := strings.Split(pos, " ")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid coordinates format")
	}
	return parts[1], parts[0], nil
}
