package app

import (
	"flag"
	"log/slog"
	"sync"

	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/models"
)

type AppDep struct {
    Logger            *slog.Logger
    SharedMap         *sync.Map
    CurrentJobs       *sync.Map
    DbChan            chan models.ConversionRecord
    ConversionRecords models.ConversionRecordInt
}

var DirectoryToWatch  = flag.String("i", "samples/Watched_folder", "Directory to watch for changes")
var OutputDirectory  = flag.String("o", "samples/Output_folder", "Directory to store the mp3 files")

func NewAppDep(logger *slog.Logger, conversionRecords models.ConversionRecordInt) *AppDep {
    return &AppDep{
        Logger:            logger,
        SharedMap:         &sync.Map{},
        CurrentJobs:       &sync.Map{},
        DbChan:            make(chan models.ConversionRecord, 10),
        ConversionRecords: conversionRecords,
    }
}
