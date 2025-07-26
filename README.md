# 🌡️ IoT Sensor Simulator con Go + Firestore + Mailjet

Este proyecto simula sensores de temperatura utilizando **Golang**, guarda las lecturas en **Google Firestore** y muestra un **frontend web en tiempo real**. Además, envía **alertas por correo electrónico** cuando la temperatura supera los 30°C.


## 🚀 Características

-  Simulación de múltiples sensores de temperatura
-  Almacenamiento en **Google Cloud Firestore**
-  Visualización de datos en frontend embebido con **Chart.js**
-  Alerta visual y por **correo electrónico** con Mailjet
-  Actualización automática de datos cada 10 segundos
-  Variables de entorno para mantener las credenciales seguras
-  API REST para acceso a datos
-  Interfaz web responsive

## 🧱 Estructura del proyecto

```
iot-simulator/
├── main.go                 # Servidor principal, simulación y frontend
├── firestore/
│   └── client.go           # Inicialización y funciones para Firestore
├── models/
│   └── models.go           # Modelo de datos de sensores
├── .env                    # Claves Mailjet API (no incluir en git)
├── go.mod                  # Dependencias del proyecto
├── go.sum                  # Checksums de dependencias
├── clave.json              # credenciales por GCP 
└── README.md               # Este archivo
```

## ⚙️ Requisitos

- [Go 1.19+](https://go.dev/dl/)
- Cuenta en [Google Cloud](https://console.cloud.google.com/) con Firestore habilitado
- Cuenta gratuita en [Mailjet](https://www.mailjet.com/)
- Git (para clonar el repositorio)

## 📦 Instalación

1. **Clona este repositorio:**
   ```bash
   git clone https://github.com/tuusuario/iot-simulator.git
   cd iot-simulator
   ```

2. **Instala las dependencias:**
   ```bash
   go mod tidy
   ```

## ⚙️ Configuración

### 1. Configuración de Google Cloud Firestore

1. Ve a [Google Cloud Console](https://console.cloud.google.com/)
2. Crea un nuevo proyecto o selecciona uno existente
3. Habilita la API de Firestore
4. Crea una cuenta de servicio y descarga el archivo de credenciales JSON
5. Agrégalo a la carpeta raíz del proyecto
  

### 2. Configuración de Mailjet

1. Regístrate en [Mailjet](https://www.mailjet.com/)
2. Obtén tu API Key y API Secret desde el dashboard
3. Crea el archivo `.env` en la carpeta raíz:
4. Edita el archivo `.env` con tus credenciales:
   ```env
   MAILJET_API_KEY=tu_api_key_aqui
   MAILJET_API_SECRET=tu_api_secret_aqui
   MAILJET_FROM_EMAIL=tuemailverificado@gmail.com
   MAILJET_FROM_NAME=IoT Alerta
   MAILJET_TO_EMAIL=destinatario@ejemplo.com
   ```

## 🚀 Uso

1. **Ejecuta el proyecto:**
   ```bash
   go run main.go
   ```

2. **El servidor iniciará en:**
   ```
   http://localhost:9090
   ```

3. **Los sensores comenzarán a generar datos automáticamente**

## 🌐 Vista Web

Una vez que el servidor esté ejecutándose, abre tu navegador y visita:
** http://localhost:9090/**

La interfaz web incluye:

- **Dashboard principal** con métricas en tiempo real
- **Tabla de lecturas** con las últimas mediciones de todos los sensores
- **Gráfico interactivo** con Chart.js que se actualiza cada 10 segundos
- **Sistema de alertas** visual cuando un sensor supera los 30°C
- **Indicadores de estado** para cada sensor

## 📧 Notificaciones por Correo

El sistema envía automáticamente correos electrónicos cuando:

- Un sensor excede los 30°C de temperatura

Las notificaciones se envían de forma asíncrona usando Mailjet:

```go
go enviarAlertaCorreo(sensorID, temperatura)
```

## 🔌 API Endpoints

El proyecto expone los siguientes endpoints REST:

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/` | Dashboard principal con interfaz web |
| `GET` | `/api` | Obtiene las últimas 10 lecturas de sensores |


### Ejemplo de respuesta:

```json
[
  {
    "ID": "sensor-1",
    "Temperatura": 25.4,
    "Timestamp": "2025-07-25T10:30:00Z"
  },
  {
    "ID": "sensor-2", 
    "Temperatura": 28.7,
    "Timestamp": "2025-07-25T10:30:15Z"
  }
]
```

## 📈 Mejoras Futuras

-  Conexión con dispositivos reales vía MQTT
-  WebSockets para actualización instantánea
-  Aplicación móvil con Flutter
-  Análisis histórico y generación de reportes PDF
-  Despliegue con Docker
-  Deploy en Google Cloud App Engine
-  Notificaciones push
-  Dashboard de administración avanzado
-  Soporte para múltiples zonas horarias
-  Sistema de autenticación y autorización

## 🔐 Seguridad

-  Las variables sensibles están almacenadas en `.env`
-  El archivo `.env` está incluido en `.gitignore`
-  Las credenciales de Google Cloud se manejan mediante variables de entorno
-  No se exponen claves API en el código fuente

### ⚠️ Importante:
- **Nunca** subas archivos `.env` o `credentials.json` a GitHub
- Usa siempre `.env.example` como plantilla
- Rota tus claves API regularmente

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Consulta el archivo `LICENSE` para más detalles.

## ✨ Autor

**Hecho por Bryan Fárez**

- GitHub: https://github.com/bfarez21

---

