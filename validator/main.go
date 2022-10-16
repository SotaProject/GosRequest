package main

import (
	"fmt"
	"github.com/SotaProject/GosRequest/validator/admin_api"
	"github.com/SotaProject/GosRequest/validator/notifications"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net"
	"sync"
)

type Config struct {
	NotificationsURL   string `required:"true" envconfig:"NOTIFICATIONS_URL"`
	NotificationsToken string `required:"true" envconfig:"NOTIFICATIONS_TOKEN"`
	AdminAPIURL        string `required:"true" envconfig:"ADMIN_API_URL"`
	AdminAPIToken      string `required:"true" envconfig:"ADMIN_API_TOKEN"`
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func LoadAppConfig() Config {
	var config Config
	envconfig.MustProcess("", &config)
	return config
}

func Handler(api events.APIGatewayV2HTTPRequest) (Response, error) {
	userAgent := api.Headers["User-Agent"]
	ip := net.ParseIP(api.Headers["X-Forwarded-For"])
	trackerUUID := api.QueryStringParameters["tracker_uuid"]
	url := api.QueryStringParameters["url"]

	data, err := admin_api.GetSubnetsData()
	if err != nil {
		err = fmt.Errorf("failed to fetch subnets data from admin_api: %w", err)
		log.Println(err)
		return Response{}, err
	}

	for _, sn := range data.Subnets {
		for _, sr := range sn.Ranges {
			_, ipNet, _ := net.ParseCIDR(sr)
			if ipNet.Contains(ip) {
				notifyData, err := admin_api.GetNotifications(trackerUUID)
				if err != nil {
					err = fmt.Errorf("failed to fetch notification data from admin_api: %w", err)
					log.Println(err)
					return Response{}, err
				}

				wg := sync.WaitGroup{}
				wg.Add(2)

				go func() {
					defer wg.Done()
					err = notifications.SendNotifications(notifications.SendData{
						TrackerUUID: trackerUUID,
						TrackerName: notifyData.TrackerName,
						SubnetCIDR:  sr,
						SubnetName:  sn.Name,
						SubnetTags:  sn.Tags,
						URL:         url,
						IP:          ip.String(),
						UserAgent:   userAgent,
						ChatID:      notifyData.ChatIDs,
					})
					if err != nil {
						err = fmt.Errorf("failed to send notification: %w", err)
						log.Println(err)
					}
				}()

				go func() {
					defer wg.Done()
					err = admin_api.AddRequest(admin_api.Request{
						TrackerUUID: trackerUUID,
						URL:         url,
						IP:          ip.String(),
						UserAgent:   userAgent,
						SubnetUUID:  sn.ID,
					})
					if err != nil {
						err = fmt.Errorf("failed to save request: %w", err)
						log.Println(err)
					}
				}()

				wg.Wait()

				return Response{StatusCode: 200, Body: "{\"gos\": true}"}, nil
			}
		}
	}

	return Response{StatusCode: 200, Body: "{\"gos\": false}"}, nil
}

func main() {
	config := LoadAppConfig()

	admin_api.SetConfig(config.AdminAPIToken, config.AdminAPIURL)
	notifications.SetConfig(config.NotificationsToken, config.NotificationsURL)

	lambda.Start(Handler)
}
