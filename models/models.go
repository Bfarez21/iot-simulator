package models

import "time"

type SensorData struct {
    ID          string    `firestore:"id"`
    Temperatura float64   `firestore:"temperatura"`
    Timestamp   time.Time `firestore:"timestamp"`
}
