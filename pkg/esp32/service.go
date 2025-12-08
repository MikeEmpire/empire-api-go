package esp32

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Service struct {
	db *sql.DB
}

func NewService(dbPath string) (*Service, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS sensor_readings (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            device_id TEXT NOT NULL,
            temperature_c REAL NOT NULL,
            temperature_f REAL NOT NULL,
            voltage REAL NOT NULL,
            raw_adc INTEGER NOT NULL,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
        );
        CREATE INDEX IF NOT EXISTS idx_device_timestamp 
            ON sensor_readings(device_id, timestamp DESC);
    `)

	if err != nil {
		return nil, err
	}

	return &Service{db: db}, nil
}

func (s *Service) SaveReading(req SensorRequest) error {
	_, err := s.db.Exec(`
        INSERT INTO sensor_readings 
        (device_id, temperature_c, temperature_f, voltage, raw_adc)
        VALUES (?, ?, ?, ?, ?)
    `, req.DeviceID, req.TemperatureC, req.TemperatureF, req.Voltage, req.RawADC)

	return err
}

func (s *Service) GetRecentReadings(deviceID string, hours int) ([]SensorReading, error) {
	rows, err := s.db.Query(`
        SELECT id, device_id, temperature_c, temperature_f, voltage, raw_adc, timestamp
        FROM sensor_readings
        WHERE device_id = ? 
        AND timestamp >= datetime('now', '-' || ? || ' hours')
        ORDER BY timestamp DESC
    `, deviceID, hours)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []SensorReading
	for rows.Next() {
		var r SensorReading
		err := rows.Scan(&r.ID, &r.DeviceID, &r.TemperatureC,
			&r.TemperatureF, &r.Voltage, &r.RawADC, &r.Timestamp)
		if err != nil {
			continue
		}
		readings = append(readings, r)
	}

	return readings, nil
}

func (s *Service) GetStats(deviceID string, hours int) (*StatsResponse, error) {
	var stats StatsResponse

	// Get most recent reading for "current"
	err := s.db.QueryRow(`
        SELECT temperature_f 
        FROM sensor_readings 
        WHERE device_id = ? 
        ORDER BY timestamp DESC 
        LIMIT 1
    `, deviceID).Scan(&stats.Current)

	if err != nil {
		return nil, err
	}

	// Get aggregates
	err = s.db.QueryRow(`
        SELECT 
            AVG(temperature_f) as avg,
            MIN(temperature_f) as min,
            MAX(temperature_f) as max
        FROM sensor_readings
        WHERE device_id = ? 
        AND timestamp >= datetime('now', '-' || ? || ' hours')
    `, deviceID, hours).Scan(&stats.Average, &stats.Min, &stats.Max)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (s *Service) GetTimeOfDayStats(deviceID string, hours int) (map[string]float64, error) {
	rows, err := s.db.Query(`
        SELECT 
            CASE 
                WHEN CAST(strftime('%H', timestamp) AS INTEGER) BETWEEN 6 AND 11 THEN 'morning'
                WHEN CAST(strftime('%H', timestamp) AS INTEGER) BETWEEN 12 AND 17 THEN 'afternoon'
                WHEN CAST(strftime('%H', timestamp) AS INTEGER) BETWEEN 18 AND 23 THEN 'evening'
                ELSE 'night'
            END as time_period,
            AVG(temperature_f) as avg_temp
        FROM sensor_readings
        WHERE device_id = ? 
        AND timestamp >= datetime('now', '-' || ? || ' hours')
        GROUP BY time_period
    `, deviceID, hours)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]float64)
	for rows.Next() {
		var period string
		var avgTemp float64
		if err := rows.Scan(&period, &avgTemp); err != nil {
			continue
		}
		result[period] = avgTemp
	}

	return result, nil
}
