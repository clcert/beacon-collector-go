package main

import (
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"time"
)

func getExternalEvents(timestamp time.Time) []string {
	db := connectDB()
	defer db.Close()

	var externalEvents []string

	getEventsStatement := `SELECT digest FROM events_collected WHERE pulse_timestamp = $1`
	rows, err := db.Query(getEventsStatement, timestamp)
	if err != nil {
		log.WithFields(log.Fields{
			"pulseTimestamp": timestamp,
		}).Panic("Failed to get events collected")
	}
	defer rows.Close()
	for rows.Next() {
		var externalEvent string
		err = rows.Scan(&externalEvent)
		if err != nil {
			log.WithFields(log.Fields{
				"pulseTimestamp": timestamp,
			}).Panic("No events collected for this pulse")
		}
		externalEvents = append(externalEvents, externalEvent)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return externalEvents

}

func generateExternalEventField(externalEvents []string, timestamp time.Time) {
	db := connectDB()
	defer db.Close()

	hashedEvents := hashEvents(externalEvents)
	externalEvent := vdf(hashedEvents)
	addEventStatement := `INSERT INTO external_events (value, pulse_timestamp) VALUES ($1, $2)`

	_, err := db.Exec(addEventStatement, hex.EncodeToString(externalEvent[:]), timestamp)
	if err != nil {
		log.WithFields(log.Fields{
			"pulseTimestamp": timestamp,
		}).Panic("Failed to add External Events to database")
	}

}

// H(e1 || e2 || ... || en)
func hashEvents(events []string) [64]byte {
	var concatenatedEvents string
	for _, l := range events {
		concatenatedEvents = concatenatedEvents + l
	}
	byteEvents := []byte(concatenatedEvents)
	return sha3.Sum512(byteEvents)
}

func vdf(events [64]byte) [64]byte {
	return sha3.Sum512(events[:])
}