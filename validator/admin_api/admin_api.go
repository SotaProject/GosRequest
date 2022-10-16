package admin_api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Config struct {
	AdminAPIURL   string
	AdminAPIToken string
}

type Subnet struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Ranges []string `json:"ranges"`
	Tags   []string `json:"tags"`
}

type AdminAPIData struct {
	Subnets     []Subnet `json:"subnets"`
	LastUpdated int64    `json:"last_updated,omitempty"`
}

type FetchNotificationResponse struct {
	TrackerName string   `json:"tracker_name"`
	ChatIDs     []string `json:"chat_ids"`
	LastUpdated int64    `json:"last_updated,omitempty"`
}

type Request struct {
	TrackerUUID string `json:"tracker_uuid"`
	URL         string `json:"url"`
	IP          string `json:"ip"`
	UserAgent   string `json:"user_agent"`
	SubnetUUID  string `json:"subnet_uuid"`
}

var config Config

var subnetsCache *AdminAPIData

func SetConfig(token string, url string) {
	config = Config{AdminAPIToken: token, AdminAPIURL: url}
}

func GetSubnetsData() (AdminAPIData, error) {
	if subnetsCache != nil && time.Now().Sub(time.Unix(subnetsCache.LastUpdated, 0)) < 120 {
		return *subnetsCache, nil
	}

	req, err := http.NewRequest("GET", config.AdminAPIURL+"/subnets_data", nil)
	if err != nil {
		return AdminAPIData{}, err
	}

	req.Header.Set("x-api-key", config.AdminAPIToken)
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		return AdminAPIData{}, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return AdminAPIData{}, err
	}

	data := AdminAPIData{}
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		return AdminAPIData{}, err
	}

	subnetsCache = &data
	return data, nil
}

func GetNotifications(trackerUUID string) (FetchNotificationResponse, error) {
	req, err := http.NewRequest("GET", config.AdminAPIURL+"/fetch_notifications", nil)
	if err != nil {
		return FetchNotificationResponse{}, err
	}

	req.Header.Set("x-api-key", config.AdminAPIToken)

	q := req.URL.Query()
	q.Add("tracker_uuid", trackerUUID)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		return FetchNotificationResponse{}, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return FetchNotificationResponse{}, err
	}

	data := FetchNotificationResponse{}
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		return FetchNotificationResponse{}, err
	}

	return data, nil
}

func AddRequest(data Request) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", config.AdminAPIURL+"/new_request", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.AdminAPIToken)
	_, err = http.DefaultClient.Do(req)
	return err
}
