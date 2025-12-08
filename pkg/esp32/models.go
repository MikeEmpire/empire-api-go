package esp32

import "time"

type SensorReading struct {
	ID           int       `json:"id"`
	DeviceID     string    `json:"device_id"`
	TemperatureC float64   `json:"temperature_c"`
	TemperatureF float64   `json:"temperature_f"`
	Voltage      float64   `json:"voltage"`
	RawADC       int       `json:"raw_adc"`
	Timestamp    time.Time `json:"timestamp"`
}

type SensorRequest struct {
	DeviceID     string  `json:"device_id" validate:"required"`
	TemperatureC float64 `json:"temperature_c" validate:"required"`
	TemperatureF float64 `json:"temperature_f" validate:"required"`
	Voltage      float64 `json:"voltage" validate:"required"`
	RawADC       int     `json:"raw_adc" validate:"required"`
}

type StatsResponse struct {
	Current float64 `json:"current"`
	Average float64 `json:"average"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
}
