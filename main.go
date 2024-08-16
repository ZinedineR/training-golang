package main

import (
	"bytes"
	"context"
	"csv-xlsx-read/entity"
	"csv-xlsx-read/externalapi"
	"csv-xlsx-read/httpclient"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/xuri/excelize/v2"
	"html/template"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	basedir      = "./data"
	filename     = "yugioh.xlsx"
	archetype    = "Melodious"
	htmlTemplate = "static/index.html"
	pdfOutput    = "output.pdf"
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
	processXLSX(ctx, *yugiohData)
	processPDF(ctx, *yugiohData)
	slog.InfoContext(ctx, fmt.Sprintf("Finished generating xlsx yugioh. Elapsed Time: %d ms", time.Since(startTime).Milliseconds()))
}

func processXLSX(ctx context.Context, yugiohData entity.YugiohAPIResponse) {
	// Step 1: Create a new Excel file
	f := excelize.NewFile()

	// Step 2: Create a sheet and set headers
	sheetName := "YugiohCards"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		panic(err)
	}
	f.SetActiveSheet(index)

	// Write header
	headers := []string{"ID", "Name", "Type", "Description", "Atk", "Def", "Race", "Attribute", "Archetype"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
	}

	// Step 3: Write data to the Excel sheet
	for i, card := range yugiohData.Data {
		row := i + 2 // Start from row 2 (row 1 is the header)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), card.Id)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), card.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), card.Type)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), card.Desc)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), card.Atk)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), card.Def)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), card.Race)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), card.Attribute)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), card.Archetype)
	}

	// Step 4: Save the Excel file
	excelFilePath := filepath.Join(basedir, filename)
	if err := f.SaveAs(excelFilePath); err != nil {
		slog.ErrorContext(ctx, "Error saving Excel file", slog.Any("error", err))
		os.Exit(1)
	}

	fmt.Println("XLSX generation complete.")
}

func processPDF(ctx context.Context, yugiohData entity.YugiohAPIResponse) {
	// Step 1: Create a new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	// Step 2: Iterate through each card and generate HTML content, then add it as a page to the PDF
	for _, card := range yugiohData.Data {
		tmpl, err := template.ParseFiles(basedir + "/" + htmlTemplate)
		if err != nil {
			log.Fatal(err)
		}

		// Create a buffer to hold the HTML content
		slog.InfoContext(ctx, "Writing card to pdf...", slog.Any("card", card.Name))
		data := entity.YugiohTemplate{
			Id:            card.Id,
			Name:          card.Name,
			Type:          card.Type,
			Desc:          card.Desc,
			Atk:           card.Atk,
			Def:           card.Def,
			Level:         card.Level,
			Race:          card.Race,
			Attribute:     card.Attribute,
			Archetype:     card.Archetype,
			ImageUrl:      card.CardImages[0].ImageUrl,
			YgoprodeckUrl: card.YgoprodeckUrl,
		}
		// Execute the template with the card data, writing the output to the buffer
		if err := tmpl.Execute(buf, data); err != nil {
			log.Println(err)
			continue
		}
		buf.WriteString(`<P style="page-break-before: always">`)
	}
	// Convert the buffer content to a string and create a new page
	page := wkhtmltopdf.NewPageReader(strings.NewReader(buf.String()))

	// Add the page to the PDF generator
	pdfg.AddPage(page)
	// Step 3: Generate the PDF document
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Save the PDF to a file
	outputPath := filepath.Join(basedir, pdfOutput)
	err = pdfg.WriteFile(outputPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PDF generation complete.")
}
