package database

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Data struct {
	TrackerUUID string
	TrackerName string
	SubnetCIDR  string
	SubnetUUID  string
	SubnetName  string
	SubnetTag   string
	ChatID      []int64
}
type Request struct {
	TrackerUUID string `json:"trackerUUID" validate:"required,uuid4"`
	URL         string `json:"url" validate:"required,url"`
	IP          string `validate:"required,ipv4"`
	UserAgent   string `validate:"required"`
	SubnetUUID  string
}

func Connect(url string, maxConnection int, maxIdleConnections int) *sql.DB {
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(maxConnection)
	db.SetMaxIdleConns(maxIdleConnections)
	return db
}

func Migrate(url string) {
	m, err := migrate.New("file://./migrations", url)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil && err.Error() != "no change" {
		panic(err)
	}
}

func (request Request) Insert(db *sql.DB) error {
	query := "INSERT INTO requests (tracker_uuid, url, ip, subnet_uuid, user_agent) VALUES ($1, $2, $3, $4, $5)"
	_, err := db.Query(query, request.TrackerUUID, request.URL, request.IP, request.SubnetUUID, request.UserAgent)
	return err
}

func GetData(db *sql.DB, ip string, trackerUUID string) (error, Data) {
	query := "SELECT trackers.uuid, trackers.name, subnet_ranges.cidr, subnets.uuid, subnets.name, subnets_tags.name, " +
		"array(select chat_id from notifications where tracker_uuid = trackers.uuid) as chat_ids " +
		"FROM subnet_ranges " +
		"LEFT JOIN subnets ON subnets.uuid = subnet_ranges.subnet_uuid " +
		"LEFT JOIN subnets_tags ON subnets.tag_uuid = subnets_tags.uuid " +
		"LEFT JOIN trackers ON trackers.uuid = $2 " +
		"WHERE subnet_ranges.cidr >>= $1 AND trackers.uuid IS NOT NULL"
	row := db.QueryRow(query, ip, trackerUUID)
	data := Data{}
	err := row.Scan(&data.TrackerUUID, &data.TrackerName, &data.SubnetCIDR, &data.SubnetUUID, &data.SubnetName, &data.SubnetTag, (*pq.Int64Array)(&data.ChatID))
	return err, data
}
