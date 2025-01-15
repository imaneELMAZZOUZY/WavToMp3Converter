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

	// Define the path to sqlite.exe in the assets folder
	sqlitePath := filepath.Join("bin", "sqlite3.exe")

	// Check if sqlite.exe exists
	if _, err := os.Stat(sqlitePath); os.IsNotExist(err) {
		fmt.Printf("sqlite3.exe not found in %s\n", sqlitePath)
		return
	}

	// Open or create the SQLite database file
	db, err := sql.Open("sqlite", "conversions.db")
	if err != nil {
		fmt.Println("Error opening SQLite database:", err)
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
        duration TEXT,
        timestamp TEXT
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
            duration,
            timestamp
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        `
		_, err := db.Exec(insertSQL, record.InputFile, record.OutputFile, record.Codec,
			record.Bitrate, record.SampleRate, record.Channels, record.ConversionStatus,
			record.Duration, record.Timestamp)
		if err != nil {
			fmt.Println("Error inserting record into database:", err)
		} else {
			fmt.Println("Record inserted successfully:", record)
		}
	}
}
