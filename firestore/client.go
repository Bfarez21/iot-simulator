package firestore

import (
    "context"
    "fmt"
    "log"

    "cloud.google.com/go/firestore"
    "google.golang.org/api/option"
    "iot-simulator/models" 
)

var client *firestore.Client

func InitFirestore() {
    ctx := context.Background()
    sa := option.WithCredentialsFile("clave.json") // archivo descargado desde GCP

    var err error
    client, err = firestore.NewClient(ctx, "iot-simulator-467016", sa)
    if err != nil {
        log.Fatalf("Error al inicializar Firestore: %v", err)
    }

    fmt.Println("🚀 Conectado a Firestore")
}

func SaveSensorData(data models.SensorData) error {
    ctx := context.Background()
    _, _, err := client.Collection("sensores").Add(ctx, data)
    return err
}


func GetLastReadings(limit int) ([]models.SensorData, error) {
    ctx := context.Background()
    docs, err := client.Collection("sensores").OrderBy("timestamp", firestore.Desc).Limit(limit).Documents(ctx).GetAll()
    if err != nil {
        return nil, err
    }

    var results []models.SensorData
    for _, doc := range docs {
        var d models.SensorData
        doc.DataTo(&d)
        results = append(results, d)
    }

    return results, nil
}
