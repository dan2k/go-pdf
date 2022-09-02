package main

import (
	"fmt"
	u "go-pdf/pdfGenerator"
	"log"
	"os"
	"time"

	"path/filepath"

	"github.com/schollz/progressbar/v3"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/xuri/excelize/v2"
)

var (
	outfile, _ = os.Create("my.log") // update path for your needs
	l          = log.New(outfile, "", 0)
)

const (
	DDMMYYYYhhmmss = "2006-01-02 15:04:05"
)

func main() {
	
	err := qrcode.WriteFile("https://example.org", qrcode.Medium, 256, "templates/qr.png")
	if err != nil {
		fmt.Println(err)
	}
	r := u.NewRequestPdf("")

	//html template path
	templatePath := "templates/sample2.html"

	//path for download pdf
	outputPath := "storage/example.pdf"

	//html template data
	type templateData struct {
		Title       string
		Description string
		Company     string
		Contact     string
		Country     string
		AddPath     string
	}
	tmp := templateData{
		Title:       "HTML to PDF generator",
		Description: "This is the simple HTML to PDF file.",
		Company:     "Jhon Lewis",
		Contact:     "Maria Anders",
		Country:     "Germany",
		AddPath:     AddPath("qr.png"),
	}
	if err := r.ParseTemplate(templatePath, tmp); err == nil {
		// ok, _ := r.GeneratePDF(outputPath)
		// fmt.Println(ok, "pdf generated successfully")
	} else {
		fmt.Println(err)
	}

	//######################################################
	var xlsFile = "test.xlsx"
	fmt.Printf("Open %s \n", xlsFile)
	l.Printf("%s Open %s \n", time.Now().UTC().Format(DDMMYYYYhhmmss), xlsFile)
	f, err := excelize.OpenFile(xlsFile)
	if err != nil {
		fmt.Println(err)
		l.Println(time.Now().UTC().Format(DDMMYYYYhhmmss), err)
		return
	}
	fmt.Println("Open file complete!")
	l.Println(time.Now().UTC().Format(DDMMYYYYhhmmss), "Open file complete!")
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
			l.Println(time.Now().UTC().Format(DDMMYYYYhhmmss), err)
		}
	}()
	// Get value from cell by given worksheet name and axis.

	rows, err := f.GetRows("Sheet1")

	if err != nil {
		fmt.Println(err)
		l.Println(time.Now().UTC().Format(DDMMYYYYhhmmss), err)
		return
	}
	fmt.Printf("Total %d rows \n", len(rows)-1)
	l.Printf("%s Total %d rows \n", time.Now().UTC().Format(DDMMYYYYhhmmss), len(rows)-1)
	bar := progressbar.NewOptions(
		len(rows)-1,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionShowCount(),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("กำลังประมวลผล..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]▓[reset]",
			SaucerHead:    "[green]▶[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	for i := 1; i < len(rows); i++ {
		r.GeneratePDF(outputPath)
		l.Println(time.Now().UTC().Format(DDMMYYYYhhmmss), "Number", i)
		bar.Add(1)

		// time.Sleep(1000 * time.Millisecond)
	}

	fmt.Println("complete !")
	l.Println(time.Now().UTC().Format(DDMMYYYYhhmmss), "complete !")

}
func AddPath(f string) string {
	return fmt.Sprintf("file://%s/%s", filepath.Dir(os.Args[0]), f)
}
