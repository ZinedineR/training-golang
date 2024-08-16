package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	basedir  = "./data"
	fileName = "data.xlsx" // Change this to your actual file name
)

func main() {
	if err := os.MkdirAll(basedir, 0777); err != nil {
		panic(err)
	}

	file, err := os.Open(basedir + "/" + fileName)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	// Detect the file extension and process accordingly
	extension, err := GetFileExtension(file)
	if err != nil {
		log.Fatalf("file extension error: %s", err)
	}

	// Reset the file pointer after checking the extension
	file.Seek(0, io.SeekStart)

	reportFile, err := os.Create(basedir + "/report.log")
	if err != nil {
		log.Fatalf("failed creating report file: %s", err)
	}
	defer reportFile.Close()
	reportFileWriter := bufio.NewWriter(reportFile)

	if extension == ".csv" {
		err := processCSV(file, reportFileWriter)
		if err != nil {
			log.Fatalf("failed processing csv file: %s", err)
		}
	} else if extension == ".xlsx" {
		err := processXLSX(file, reportFileWriter)
		if err != nil {
			log.Fatalf("failed processing xlsx file: %s", err)
		}
	} else {
		log.Fatalf("unsupported file extension: %s", extension)
	}
}

func processCSV(file io.Reader, reportFileWriter *bufio.Writer) error {
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return err
	}

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		name, address := line[0], line[1]
		report := fmt.Sprintf("Row %d: name: %s address: %s \n", i, name, address)
		_, _ = fmt.Fprintf(reportFileWriter, report)
		_ = reportFileWriter.Flush()
	}

	return nil
}

func processXLSX(file io.Reader, reportFileWriter *bufio.Writer) error {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return err
	}

	f, err := excelize.OpenReader(buf)
	if err != nil {
		return err
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}

	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}
		name := row[0]
		address := row[1]
		report := fmt.Sprintf("Row %d: name: %s address: %s \n", i, name, address)
		_, _ = fmt.Fprintf(reportFileWriter, report)
		_ = reportFileWriter.Flush()
	}

	return nil
}

func GetFileExtension(file io.Reader) (string, error) {
	byteFile, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	mimeType := mimetype.Detect(byteFile)

	// Reset the reader after reading it
	file = bytes.NewReader(byteFile)

	extension := filepath.Ext(fileName)

	// Extension or mimetype checker, if not .csv/.xlsx return error
	if extension != ".csv" && mimeType.String() != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		return mimeType.String(), errors.New("extension invalid")
	}

	return extension, nil
}
