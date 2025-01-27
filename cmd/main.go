package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/app"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/helpers"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/models"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/services"

	_ "modernc.org/sqlite"
)



func main() {

	port := flag.String("port", ":5000", "HTTP network address")
	
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := connectDB()
	if err != nil {
		logger.Error("Error connecting to database:", "error", err)
		return
	}
	defer db.Close()

	conversionRecords := &models.ConversionRecordModel{Db: db}
    appDep := app.NewAppDep(logger, conversionRecords)

	go services.Watch(appDep)

	go services.Process(appDep,services.NewJobConverter())

	go func() {

		appDep.ConversionRecords.CreateTable()

		for record := range appDep.DbChan {
			err := appDep.ConversionRecords.Insert(record.InputFile,
				record.OutputFile, record.Codec, record.Bitrate,
				record.SampleRate, record.Channels, record.ConversionStatus,
				record.StartTime, record.EndTime)

			if err != nil {
				logger.Error("Error inserting record into database:", "error", err)
			} else {
				
				logger.Info("Insertion in db successfully done!", "input_file", record.InputFile, "status", record.ConversionStatus)
			}
		}
	}()

	log.Printf("Starting server on %s", *port)
	err = http.ListenAndServe(*port, appDep.Routes())
	log.Fatal(err)

}

func connectDB() (*sql.DB, error) {
	// Define the path to sqlite.exe in the assets folder
	sqlitePath := filepath.Join("bin", "sqlite3.exe")

	if !helpers.IsFileExist(sqlitePath) {
		return nil, fmt.Errorf("%w in %s", app.ErrSqliteNotFound, sqlitePath)
	}
	// Open or create the SQLite database file
	db, err := sql.Open("sqlite", "conversions.db")
	if err != nil {
		fmt.Println("Error opening SQLite database:", err)
		return nil, err
	}
	return db, nil
}
