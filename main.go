package main

import (
	u "go-pdf/pdfGenerator"
	"os"
	"strconv"
	"sync"
	"time"
	"github.com/xuri/excelize/v2"
)

const (
	DDMMYYYYhhmmss  = "2006-01-02 15:04:05"
	DDMMYYYYhhmmss2 = "20060102150405"
)
var loc *time.Location
var runtime string
var wg sync.WaitGroup
var guard = make(chan struct{}, 50)
func main() {
	lo, _ := time.LoadLocation("Asia/Bangkok")
	loc = lo
	runtime = time.Now().In(loc).Format(DDMMYYYYhhmmss2)
	CallClear()
	InitLog()
	InitFlag()
	r := u.NewRequestPdf("", l, envs)
	//html template path
	genLog("Aplication name:" , envs["APPNAME"])
	genLog("VERSION:", envs["VERSION"])
	mx, _ := strconv.Atoi(envs["MAXGOROUTINES"])
	guard = make(chan struct{}, mx)
	genLog("MAX GOROUTINES PER RUNTIME = ", envs["MAXGOROUTINES"])
	genLog("Generate struct template complete !")
	templatePath := envs["TEMPDIR"]+"/tpl/" + envs["TEMPFILE"]
	//#####################################################
	//path for download pdf
	if _, err := os.Stat(envs["STORAGE"]); err != nil {
		if err := os.MkdirAll(envs["STORAGE"], os.ModePerm); err != nil {
			l.Fatal(err)
		}
	}
	if _, err := os.Stat(envs["TEMPDIR"]+ "/qrcode/" + runtime); err != nil {
		if err := os.MkdirAll(envs["TEMPDIR"]+"/qrcode", os.ModePerm); err != nil {
			l.Fatal(err)
		}
	}
	outputPath := envs["STORAGE"]
	//######################################################
	var xlsFile = envs["EXCELFILE"]
	genLog("Open ",xlsFile)
	f, err := excelize.OpenFile(xlsFile)
	if err != nil {
		genLog(err)
		return
	}
	genLog("Open file complete!")
	// Get value from cell by given worksheet name and axis.
	rows, err := f.GetRows("Sheet1")

	if err != nil {
		genLog(err)
		return
	}
	if err := os.MkdirAll(outputPath+"/"+runtime, 0777); err != nil {
		genLog(err)
	}
	if EndRow > (len(rows) - 1) {
		// fmt.Println("end > ", len(rows)-1)
		genLog("end > ",strconv.Itoa(len(rows)-1))
		return
	}
	useRows := rows[StartRow : EndRow+1]
	genLog("Start generate from",useRows[0][0],"(",strconv.Itoa(StartRow),") to",useRows[len(useRows)-1][0],"(",strconv.Itoa(EndRow),")")
	genLog("Totals",strconv.Itoa(len(useRows)),"records")
	bar := InitBar(len(useRows))
	for i := 0; i < len(useRows); i++ {
		guard <- struct{}{}
		wg.Add(1)
		go generate(i, useRows[i], templatePath, outputPath, r, bar)
	}
	close(guard)
	wg.Wait()

	defer func() {
		// Close the spreadsheet.
		
		
		
		if err := f.Close(); err != nil {
			genLog(err)
		}
		
	}()
}

