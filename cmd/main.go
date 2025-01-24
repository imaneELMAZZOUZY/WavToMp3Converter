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
	"sync"

	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/models"
	_ "modernc.org/sqlite"
)

type appDep struct {
	logger            *slog.Logger
	sharedMap         *sync.Map
	CurrentJobs       *sync.Map
	dbChan            chan models.ConversionRecord
	conversionRecords models.ConversionRecordInt
}

var DirectoryToWatch  = flag.String("d", "samples/Watched_folder", "Directory to watch for changes")

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

	appDep := &appDep{
		logger: logger,
		sharedMap: &sync.Map{},
		CurrentJobs: &sync.Map{},
		dbChan:            make(chan models.ConversionRecord, 10),
		conversionRecords: &models.ConversionRecordModel{Db: db},

	}

	go appDep.Watch()

	go appDep.Process(NewJobConverter())

	go func() {

		appDep.conversionRecords.CreateTable()

		for record := range appDep.dbChan {
			err := appDep.conversionRecords.Insert(record.InputFile,
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
	err = http.ListenAndServe(*port, appDep.routes())
	log.Fatal(err)

}

func connectDB() (*sql.DB, error) {
	// Define the path to sqlite.exe in the assets folder
	sqlitePath := filepath.Join("bin", "sqlite3.exe")

	if !isFileExist(sqlitePath) {
		return nil, fmt.Errorf("%w in %s", ErrSqliteNotFound, sqlitePath)
	}
	// Open or create the SQLite database file
	db, err := sql.Open("sqlite", "conversions.db")
	if err != nil {
		fmt.Println("Error opening SQLite database:", err)
		return nil, err
	}
	return db, nil
}
