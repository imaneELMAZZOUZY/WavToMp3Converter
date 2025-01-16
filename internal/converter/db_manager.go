package converter

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/models"
	_ "modernc.org/sqlite"
)

// DB Goroutine that listens to dbChan and inserts conversion records into SQLite
func ManageDb(dbChan <-chan models.ConversionRecord) {

	db , err:= connectDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	// Create table if it doesn't exist
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS conversion_records (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        input_file TEXT,
        output_file TEXT,
        codec TEXT,
        bitrate TEXT,
        sample_rate TEXT,
        channels TEXT,
        conversion_status TEXT,
        start_time TEXT,
        end_time TEXT
    );
    `
	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	// Listen for records and insert into DB
	for record := range dbChan {
		insertSQL := `
        INSERT INTO conversion_records (
            input_file,
            output_file,
            codec,
            bitrate,
            sample_rate,
            channels,
            conversion_status,
            start_time,
            end_time
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        `
		_, err := db.Exec(insertSQL, record.InputFile, record.OutputFile, record.Codec,
			record.Bitrate, record.SampleRate, record.Channels, record.ConversionStatus,
			record.StartTime, record.EndTime)
		if err != nil {
			fmt.Println("Error inserting record into database:", err)
		} else {
			fmt.Println("Record inserted successfully:", record)
		}
	}
}


func connectDB() (*sql.DB , error) {
	// Define the path to sqlite.exe in the assets folder
	sqlitePath := filepath.Join("bin", "sqlite3.exe")

	// Check if sqlite.exe exists
	if _, err := os.Stat(sqlitePath); os.IsNotExist(err) {
		fmt.Printf("sqlite3.exe not found in %s\n", sqlitePath)
		return nil, err
	}

	// Open or create the SQLite database file
	db, err := sql.Open("sqlite", "conversions.db")
	if err != nil {
		fmt.Println("Error opening SQLite database:", err)
		return nil, err
	}
	return db, nil
}

// GetRecords retrieves conversion records from the database based on the provided status.
// If status is an empty string, it retrieves all records.
func GetRecords(status string) ([]models.ConversionRecord, error) {

	db, err := connectDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil, err
	}
	defer db.Close()

	// Query the database for records based on the status
	var rows *sql.Rows
	if status == "" {
		rows, err = db.Query("SELECT * FROM conversion_records")
	} else {
		rows, err = db.Query("SELECT * FROM conversion_records WHERE conversion_status = ?", status)
	}

	if err != nil {
		fmt.Println("Error querying database:", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the result set and scan each row into a ConversionRecord
	var records []models.ConversionRecord
	for rows.Next() {
		var record models.ConversionRecord
		err := rows.Scan(&record.Id, &record.InputFile, &record.OutputFile, &record.Codec,
			&record.Bitrate, &record.SampleRate, &record.Channels, &record.ConversionStatus,
			&record.StartTime, &record.EndTime)
		if err != nil {
			fmt.Println("Error scanning record:", err)
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}






