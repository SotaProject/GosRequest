package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"log"
	"validator/database"
	"validator/notifications"
)

var db *sql.DB

func SetupDB(dbInstance *sql.DB) {
	db = dbInstance
}

func ProcessRequest(c *fiber.Ctx) error {
	c.Accepts("application/json")

	body := new(database.Request)
	if err := c.BodyParser(body); err != nil {
		return errors.New("invalid request")
	}

	body.IP = c.IP()
	body.UserAgent = c.GetReqHeaders()["User-Agent"]
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return errors.New("invalid request")
	}

	err, data := database.GetData(db, body.IP, body.TrackerUUID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(fmt.Sprintf("Unable to find subnet: %s", err.Error()))
		return c.SendStatus(500)
	}

	if data.SubnetUUID != "" && len(data.ChatID) > 0 {
		body.SubnetUUID = data.SubnetUUID
		if err := body.Insert(db); err != nil {
			log.Println(fmt.Sprintf("Unable to insert into requests: %s", err.Error()))
			return c.SendStatus(500)
		}

		err := notifications.SendNotifications(notifications.SendData{
			TrackerUUID: data.TrackerUUID,
			TrackerName: data.TrackerName,
			SubnetCIDR:  data.SubnetCIDR,
			SubnetName:  data.SubnetName,
			SubnetTag:   data.SubnetTag,
			URL:         body.URL,
			IP:          body.IP,
			UserAgent:   body.UserAgent,
			ChatID:      data.ChatID,
		})
		if err != nil {
			log.Println(err.Error())
		}

		log.Println(fmt.Sprintf("Got request from %s(%s), tracker: %s, notifications: %d", body.IP, data.SubnetName, data.TrackerUUID, len(data.ChatID)))
	}

	return c.SendStatus(204)
}
