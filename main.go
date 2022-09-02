package main

import (
	"fmt"
	u "go-pdf/pdfGenerator"
	"github.com/schollz/progressbar/v3"
	"time"
	"github.com/xuri/excelize/v2"
	"log"
    "os"
)
var (
    outfile, _ = os.Create("my.log") // update path for your needs
    l      = log.New(outfile, "", 0)
)
const (
    DDMMYYYYhhmmss = "2006-01-02 15:04:05"
)
func main() {

	r := u.NewRequestPdf("")

	//html template path
	templatePath := "templates/sample2.html"

	//path for download pdf
	outputPath := "storage/example.pdf"

	//html template data
	templateData := struct {
		Title       string
		Description string
		Company     string
		Contact     string
		Country     string
	}{
		Title:       "HTML to PDF generator",
		Description: "This is the simple HTML to PDF file.",
		Company:     "Jhon Lewis",
		Contact:     "Maria Anders",
		Country:     "Germany",
	}

	if err := r.ParseTemplate(templatePath, templateData); err == nil {
		// ok, _ := r.GeneratePDF(outputPath)
		// fmt.Println(ok, "pdf generated successfully")
	} else {
		fmt.Println(err)
	}

	/*
	bar := progressbar.NewOptions(
		100, 
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionShowCount(),
		progressbar.OptionSetDescription("กำลังประมวลผล..."),
	)
	
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(40 * time.Millisecond)
	}
	
	fmt.Println("")
	*/
	//######################################################
	var xlsFile ="test.xlsx"
	fmt.Printf("Open %s \n",xlsFile)
	f, err := excelize.OpenFile(xlsFile)
    if err != nil {
        fmt.Println(err)
        return
    }
	fmt.Println("Open file complete!")
    defer func() {
        // Close the spreadsheet.
        if err := f.Close(); err != nil {
            fmt.Println(err)
        }
    }()
    // Get value from cell by given worksheet name and axis.
	
	rows, err := f.GetRows("Sheet1")

    if err != nil {
        fmt.Println(err)
        return
    }
	fmt.Printf("Total %d rows \n",len(rows)-1)
	bar := progressbar.NewOptions(
		len(rows)-1, 
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionShowCount(),
		progressbar.OptionSetDescription("กำลังประมวลผล..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]▓[reset]",
        	SaucerHead:    "[green]▶[reset]",
			SaucerPadding: "",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	/*
    for _, row := range rows {
        for _, colCell := range row {
            fmt.Print(colCell, "\t")
        }
        fmt.Println()
    }
	*/
	
	for i := 1; i < len(rows); i++ {
		r.GeneratePDF(outputPath)
		l.Println(time.Now().UTC().Format(DDMMYYYYhhmmss),":","Number",i)
		bar.Add(1)

		// time.Sleep(1000 * time.Millisecond)
	}
	
	fmt.Println("")

}


