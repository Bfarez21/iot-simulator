package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
	"github.com/joho/godotenv"  //para usar variable de entorno
	"os"
	"log"

	"iot-simulator/firestore"
	"iot-simulator/models"
	mailjet "github.com/mailjet/mailjet-apiv3-go/v4"
)

func simulateSensor(id string, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 20; i++ { // Más iteraciones para simulación más larga
		data := models.SensorData{
			ID:          id,
			Temperatura: 20 + rand.Float64()*15, // entre 20°C y 35°C
			Timestamp:   time.Now(),
		}

		//  Alerta si la temperatura supera los 30°C	
		// guardar en firestore
		err := firestore.SaveSensorData(data)
		if data.Temperatura > 30 {
			fmt.Printf("  ALERTA: Sensor %s reportó alta temperatura: %.2f°C\n", id, data.Temperatura)
			go enviarAlertaCorreo(id, data.Temperatura) //  Llamada en goroutine

		}

		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Sensor %s: %.2f°C\n", id, data.Temperatura)
		}

		time.Sleep(30 * time.Second) // envía cada 30 segundos para mayor frecuencia
	}
}

func main() {
	// Cargar .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando el archivo .env")
	}

	rand.Seed(time.Now().UnixNano())
	firestore.InitFirestore()
	go startSimulation() // Simulación en background

	// Definir handlers
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/api", apiDataHandler) // debe coincide con el fetch de JavaScript

	fmt.Println("Servidor web en http://localhost:9090/")
	http.ListenAndServe(":9090", nil)
}

// genera o simula los datos de los sensores
func startSimulation() {
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go simulateSensor(fmt.Sprintf("sensor-%d", i), &wg)
	}
	wg.Wait()
}

// html para visualización de datos
func serveHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, `
	<!DOCTYPE html>
	<html lang="es">
	<head>
		<meta charset="UTF-8">
		<title>Panel de Sensores IoT</title>
		<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
		<style>
			body {
				font-family: 'Segoe UI', sans-serif;
				background-color: #f0f2f5;
				color: #333;
				padding: 20px;
			}
			h2 {
				text-align: center;
				color: #007BFF;
			}
			#chart-container {
				width: 80%;
				margin: 20px auto;
			}
			table {
				width: 80%;
				margin: 30px auto;
				border-collapse: collapse;
				box-shadow: 0 2px 8px rgba(0,0,0,0.1);
				background-color: #fff;
			}
			th, td {
				padding: 10px;
				text-align: center;
				border-bottom: 1px solid #ddd;
			}
			th {
				background-color: #007BFF;
				color: white;
			}
			tr:hover {
				background-color: #f9f9f9;
			}
			.status-indicator {
				padding: 4px 8px;
				border-radius: 4px;
				font-size: 0.9em;
			}
			.status-normal { background-color: rgb(212, 237, 218); color: rgb(21, 87, 36); }
			.status-warning { background-color: rgb(255, 243, 205); color: rgb(133, 100, 4); }
			.status-critical { background-color: rgb(248, 215, 218); color: rgb(114, 28, 36); }
			
			#chart-container {
				background: white;
				border-radius: 8px;
				padding: 20px;
				box-shadow: 0 2px 8px rgba(0,0,0,0.1);
			}
		</style>
	</head>
	<body>
		<h2>Lecturas de temperatura - IoT en tiempo real</h2>
		<div id="chart-container">
			<canvas id="tempChart"></canvas>
		</div>

		<table id="dataTable">
			<tr><th>Sensor</th><th>Temperatura (°C)</th><th>Fecha</th></tr>
		</table>

		<script>
			let lastUpdateTime = null;
			let sensorDatasets = new Map();

			async function fetchData() {
				try {
					const res = await fetch('/api');
					if (!res.ok) {
						throw new Error('Network response was not ok');
					}
					const data = await res.json();
					return data;
				} catch (error) {
					console.error('Error fetching data:', error);
					return [];
				}
			}

			// Configuración de colores para cada sensor
			const sensorColors = {
				'sensor-1': { border: 'rgba(255, 99, 132, 1)', background: 'rgba(255, 99, 132, 0.2)' },
				'sensor-2': { border: 'rgba(54, 162, 235, 1)', background: 'rgba(54, 162, 235, 0.2)' },
				'sensor-3': { border: 'rgba(255, 206, 86, 1)', background: 'rgba(255, 206, 86, 0.2)' }
			};

			const ctx = document.getElementById('tempChart').getContext('2d');
			const chart = new Chart(ctx, {
				type: 'line',
				data: {
					labels: [],
					datasets: []
				},
				options: {
					responsive: true,
					interaction: {
						mode: 'index',
						intersect: false,
					},
					scales: {
						x: {
							display: true,
							title: {
								display: true,
								text: 'Tiempo'
							}
						},
						y: {
							display: true,
							title: {
								display: true,
								text: 'Temperatura (°C)'
							},
							beginAtZero: false,
							suggestedMin: 15,
							suggestedMax: 40
						}
					},
					plugins: {
						legend: {
							display: true,
							position: 'top'
						},
						title: {
							display: true,
							text: 'Temperaturas en Tiempo Real'
						}
					}
				}
			});

			function updateChart(data) {
				// Agrupar datos por sensor
				const sensorGroups = {};
				data.forEach(entry => {
					if (!sensorGroups[entry.ID]) {
						sensorGroups[entry.ID] = [];
					}
					sensorGroups[entry.ID].push(entry);
				});

				// Limpiar datasets existentes
				chart.data.datasets = [];
				const allLabels = new Set();

				// Crear dataset para cada sensor
				Object.keys(sensorGroups).forEach(sensorId => {
					const sensorData = sensorGroups[sensorId];
					const color = sensorColors[sensorId] || { 
						border: 'rgba(75, 192, 192, 1)', 
						background: 'rgba(75, 192, 192, 0.2)' 
					};

					const dataset = {
						label: sensorId,
						data: [],
						borderColor: color.border,
						backgroundColor: color.background,
						fill: false,
						tension: 0.4,
						pointRadius: 4,
						pointHoverRadius: 6
					};

					sensorData.forEach(entry => {
						const timeLabel = new Date(entry.Timestamp).toLocaleTimeString();
						dataset.data.push({
							x: timeLabel,
							y: entry.Temperatura
						});
						allLabels.add(timeLabel);
					});

					chart.data.datasets.push(dataset);
				});

				// Ordenar labels por tiempo
				chart.data.labels = Array.from(allLabels).sort();
				chart.update('none'); // Actualización sin animación para mejor rendimiento
			}

			function updateTable(data) {
				const table = document.getElementById('dataTable');
				
				// Crear filas de la tabla
				let tableHTML = "<tr><th>Sensor</th><th>Temperatura (°C)</th><th>Estado</th><th>Fecha</th></tr>";
				
				data.forEach(entry => {
					const temp = entry.Temperatura;
					let estado = "Normal";
					let estadoClass = "";
					
					if (temp > 32) {
						estado = "CRITICA";
						estadoClass = "style='color: red; font-weight: bold;'";
					} else if (temp > 28) {
						estado = "ALTA";
						estadoClass = "style='color: orange; font-weight: bold;'";
					} else {
						estado = "Normal";
						estadoClass = "style='color: green;'";
					}

tableHTML += `+"`"+`
	<tr>
		<td>${entry.ID}</td>
		<td>${temp.toFixed(2)}°C</td>
		<td ${estadoClass}>${estado}</td>
		<td>${new Date(entry.Timestamp).toLocaleString()}</td>
	</tr>
`+"`"+`;
				});
				
				
				table.innerHTML = tableHTML;
			}

			async function updateUI() {
				const data = await fetchData();
				
				if (data && data.length > 0) {
					// Verificar si hay nuevos datos
					const latestTimestamp = Math.max(...data.map(d => new Date(d.Timestamp).getTime()));
					
					if (!lastUpdateTime || latestTimestamp > lastUpdateTime) {
						console.log('Actualizando datos...', new Date().toLocaleTimeString());
						lastUpdateTime = latestTimestamp;
						
						updateTable(data);
						updateChart(data);
						
						// Mostrar notificación si hay temperaturas críticas
						const criticalTemp = data.find(d => d.Temperatura > 32);
						if (criticalTemp) {
							showNotification('ALERTA: ' + criticalTemp.ID + ' reporta ' + criticalTemp.Temperatura.toFixed(2) + '°C');
						}
					}
				}
			}

			function showNotification(message) {
				// Crear notificación visual
				const notification = document.createElement('div');
				notification.style.position = 'fixed';
				notification.style.top = '20px';
				notification.style.right = '20px';
				notification.style.background = 'rgb(255, 68, 68)';
				notification.style.color = 'white';
				notification.style.padding = '15px';
				notification.style.borderRadius = '5px';
				notification.style.boxShadow = '0 4px 8px rgba(0,0,0,0.3)';
				notification.style.zIndex = '1000';
				notification.style.fontWeight = 'bold';
				
				notification.textContent = message;
				document.body.appendChild(notification);
				
				// Remover después de 5 segundos
				setTimeout(() => {
					document.body.removeChild(notification);
				}, 5000);
			}

			// Inicial
			updateUI();
			
			// Actualiza cada 5 segundos para mayor responsividad
			setInterval(updateUI, 5000);
			
			// Título dinámico con estado
			setInterval(() => {
				const now = new Date().toLocaleTimeString();
				document.title = 'Panel IoT - ' + now;
			}, 1000);
		</script>
	</body>
	</html>
	`)
}

func apiDataHandler(w http.ResponseWriter, r *http.Request) {
	// Configurar headers CORS si es necesario
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data, err := firestore.GetLastReadings(10)
	if err != nil {
		fmt.Printf("Error al obtener datos: %v\n", err)
		http.Error(w, `{"error": "Error al obtener datos"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Printf("Error al codificar JSON: %v\n", err)
		http.Error(w, `{"error": "Error al codificar datos"}`, http.StatusInternalServerError)
		return
	}
}

///funcion para enviar correos
func enviarAlertaCorreo(sensorID string, temperatura float64) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ Error cargando .env:", err)
		return
	}

	apiKey := os.Getenv("MAILJET_API_KEY")
	apiSecret := os.Getenv("MAILJET_API_SECRET")
	fromEmail := os.Getenv("MAILJET_FROM_EMAIL")
	fromName := os.Getenv("MAILJET_FROM_NAME")
	toEmail := os.Getenv("MAILJET_TO_EMAIL")

	mj := mailjet.NewMailjetClient(apiKey, apiSecret)

	messageInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: fromEmail,
				Name:  fromName,
			},
			To: &mailjet.RecipientsV31{
				{
					Email: toEmail,
					Name:  "Usuario de Alerta",
				},
			},
			Subject:  "🚨 Alerta de temperatura",
			TextPart: fmt.Sprintf("⚠️ El sensor %s reportó %.2f°C. Supera el umbral permitido.", sensorID, temperatura),
			HTMLPart: fmt.Sprintf("<h3>⚠️ Alerta desde IoT</h3><p>El <b>sensor %s</b> reportó <b>%.2f°C</b>.</p>", sensorID, temperatura),
		},
	}

	messages := mailjet.MessagesV31{Info: messageInfo}

	_, err = mj.SendMailV31(&messages)
	if err != nil {
		fmt.Println("❌ Error al enviar correo:", err)
	} else {
		fmt.Println("✅ Alerta enviada correctamente")
	}
}
