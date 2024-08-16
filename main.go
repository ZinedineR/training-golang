package main

import (
	"bufio"
	"context"
	"csv-xlsx-read/externalapi"
	"csv-xlsx-read/httpclient"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

const (
	basedir   = "./data"
	filename  = "yugioh.csv"
	archetype = "Lyrilusc"
)

var (
	httpClient httpclient.Client
	logger     *slog.Logger
)

func main() {
	if err := os.MkdirAll(basedir, 0777); err != nil {
		panic(err)
	}
	httpClientFactory := httpclient.New()
	httpClient = httpClientFactory.CreateClient()

	yugiohExternal := externalapi.NewYugiohExternalImpl(httpClient)

	ctx := context.Background()
	startTime := time.Now()
	logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)
	slog.InfoContext(ctx, "Start retrieving yugioh data", slog.Any("coba-key", "coba-value"))
	yugiohData, statusCode, err := yugiohExternal.Get(archetype)
	if err != nil {
		slog.ErrorContext(ctx, "Error retrieving data", slog.Any("error", err))
		os.Exit(1)
	}
	if statusCode != 200 {
		slog.ErrorContext(ctx, "Error retrieving data", slog.Any("error", err))
		os.Exit(1)
	}
	// Handle report file creation and writing
	reportFilePath := filepath.Join(basedir, filename)
	reportFile, err := os.Create(reportFilePath)
	if err != nil {
		slog.ErrorContext(ctx, "Error creating report file", slog.Any("error", err))
		os.Exit(1)
	}
	defer reportFile.Close()

	reportFileWriter := bufio.NewWriter(reportFile)
	defer func() {
		if err := reportFileWriter.Flush(); err != nil {
			slog.ErrorContext(ctx, "Error flushing report file", slog.Any("error", err))
		}
	}()

	// Write header
	_, err = fmt.Fprintf(reportFileWriter, "ID,Name,Type,Description,Atk,Def,Race,Attribute,Archetype\n")
	if err != nil {
		slog.ErrorContext(ctx, "Error writing to report file", slog.Any("error", err))
		os.Exit(1)
	}

	// Write data
	for _, card := range yugiohData.Data {
		_, err = fmt.Fprintf(reportFileWriter, "%d,%s,%s,%s,%d,%d,%s,%s,%s\n",
			card.Id, card.Name, card.Type, card.Desc, card.Atk, card.Def, card.Race, card.Attribute, card.Archetype)
		if err != nil {
			slog.ErrorContext(ctx, "Error writing to report file", slog.Any("error", err))
			os.Exit(1)
		}
	}

	slog.InfoContext(ctx, fmt.Sprintf("Finished generating csv yugioh. Elapsed Time: %d ms", time.Since(startTime).Milliseconds()))
}
