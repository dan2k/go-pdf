package main

import (
	"fmt"
	"flag"
	log "go-pdf/log"
	u "go-pdf/pdfGenerator"
	"time"
	"github.com/schollz/progressbar/v3"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/xuri/excelize/v2"
)

const (
	DDMMYYYYhhmmss = "2006-01-02 15:04:05"
)

var envs map[string]string

func main() {
	var startRow,endRow int 
	// endRow :=flag.Int("end",0,"an int ")
	// startRow := flag.Int("start", 0, "an int")
	// fmt.Println(*startRow,*endRow);
	flag.IntVar(&startRow, "start", 0, "a string var")
	flag.IntVar(&endRow, "end", 0, "a string var")
	flag.Parse()

	fmt.Println(startRow,endRow);
	log.New()
	envs :=log.Envs
	l :=log.L
	err := qrcode.WriteFile("https://example.org", qrcode.Medium, 256, envs["QRCODE"])
	if err != nil {
		fmt.Println(err)
	}
	r := u.NewRequestPdf("",l,envs)

	//html template path
	templatePath := envs["TEMPDIR"] + "/" + envs["TEMPFILE"]
	//path for download pdf
	outputPath := envs["STORAGE"] + "/" + envs["PDFFILE"]

	//html template data
	type templateData struct {
		Title       string
		Description string
		Company     string
		Contact     string
		Country     string
		Labels      []string
		Data        []int
	}
	tmp := templateData{
		Title:       "HTML to PDF generator",
		Description: "This is the simple HTML to PDF file.",
		Company:     "Jhon Lewis",
		Contact:     "Maria Anders",
		Country:     "Germany",
		Labels: 	 []string{"Red", "Blue", "Yellow", "Green", "Purple", "Orange"},
		Data:        []int{12, 19, 3, 5, 2, 3},
	}
	if err := r.ParseTemplate(templatePath, tmp); err == nil {
		// ok, _ := r.GeneratePDF(outputPath)
		// fmt.Println(ok, "pdf generated successfully")
	} else {
		fmt.Println(err)
		l.Println(time.Now().UTC().Format(DDMMYYYYhhmmss), err)
	}

	//######################################################
	var xlsFile = envs["EXCELFILE"]
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
		progressbar.OptionOnCompletion(func(){
			fmt.Println("")
			fmt.Println("complete !")
			l.Println(time.Now().UTC().Format(DDMMYYYYhhmmss), "complete !")
		}),
	)
	for i := 1; i < len(rows); i++ {
		r.GeneratePDF(outputPath)
		l.Println(time.Now().UTC().Format(DDMMYYYYhhmmss), "Number", i)
		bar.Add(1)

		// time.Sleep(1000 * time.Millisecond)
	}
	

	

}
