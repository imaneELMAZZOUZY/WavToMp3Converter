package models

import "database/sql"

type ConversionRecordInt interface {

	Insert(InputFile string, OutputFile string, Codec string,
		Bitrate string, SampleRate string, Channels string,
		ConversionStatus string, StartTime string, EndTime string) error

	GetAll() ([]ConversionRecord, error)
	
	Get(status string) ([]ConversionRecord, error)

	CreateTable() error
}

type ConversionRecordModel struct {
	Db *sql.DB
}

type ConversionRecord struct {
	Id               int
	InputFile        string
	OutputFile       string
	Codec            string
	Bitrate          string
	SampleRate       string
	Channels         string
	ConversionStatus string
	StartTime        string
	EndTime          string
}

func (c *ConversionRecordModel) CreateTable() error {
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
	_, err := c.Db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConversionRecordModel) Insert(InputFile string, OutputFile string, Codec string,
	Bitrate string, SampleRate string, Channels string, ConversionStatus string,
	StartTime string, EndTime string) error {

	_, err := c.Db.Exec(`INSERT INTO conversion_records (
	input_file, 
	output_file, 
	codec, bitrate, 
	sample_rate, 
	channels, 
	conversion_status, 
	start_time, end_time) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		InputFile, OutputFile, Codec, Bitrate, SampleRate,
		Channels, ConversionStatus, StartTime, EndTime)

	if err != nil {
		return err
	}
	
	return nil
}



func (c *ConversionRecordModel) GetAll() ([]ConversionRecord, error) {
	
	rows, err := c.Db.Query("SELECT * FROM conversion_records")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversionRecords []ConversionRecord

	for rows.Next() {
		var cr ConversionRecord

		err = rows.Scan(&cr.Id, &cr.InputFile, &cr.OutputFile, &cr.Codec,&cr.Bitrate,
		 &cr.SampleRate, &cr.Channels, &cr.ConversionStatus, &cr.StartTime, &cr.EndTime)
		if err != nil {
			return nil, err
		}

		conversionRecords = append(conversionRecords, cr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return conversionRecords, nil

}


func (c *ConversionRecordModel) Get(status string) ([]ConversionRecord, error) {
	
	rows, err := c.Db.Query("SELECT * FROM conversion_records WHERE conversion_status = ?", status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversionRecords []ConversionRecord

	for rows.Next() {
		var cr ConversionRecord

		err = rows.Scan(&cr.Id, &cr.InputFile, &cr.OutputFile, &cr.Codec,&cr.Bitrate,
		 &cr.SampleRate, &cr.Channels, &cr.ConversionStatus, &cr.StartTime, &cr.EndTime)
		if err != nil {
			return nil, err
		}

		conversionRecords = append(conversionRecords, cr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return conversionRecords, nil

}