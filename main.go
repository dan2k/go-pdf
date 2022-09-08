package main

import (
	"fmt"
	u "go-pdf/pdfGenerator"
	"os"
	"time"
	"github.com/xuri/excelize/v2"
	"sync"
)
const (
	DDMMYYYYhhmmss = "2006-01-02 15:04:05"
	DDMMYYYYhhmmss2 = "20060102150405"
)
// var runtime =time.Now().In(loc).Format(DDMMYYYYhhmmss2)
var loc *time.Location
var runtime string 
var wg sync.WaitGroup
func main() {
	
	lo, _ := time.LoadLocation("Asia/Bangkok")
	loc=lo
	runtime =time.Now().In(loc).Format(DDMMYYYYhhmmss2)
	CallClear() 
	InitLog()
	InitFlag()
	r := u.NewRequestPdf("",l,envs)
	//html template path
	templatePath := envs["TEMPDIR"] + "/" + envs["TEMPFILE"]
	//#####################################################
	//path for download pdf
	if _, err := os.Stat(envs["STORAGE"]); err !=nil {
		if err := os.MkdirAll(envs["STORAGE"], os.ModePerm); err != nil {
			l.Fatal(err)
		}
	}
	outputPath := envs["STORAGE"]
	//######################################################
	var xlsFile = envs["EXCELFILE"]
	fmt.Printf("Open %s \n", xlsFile)
	l.Printf("%s Open %s \n", time.Now().In(loc).Format(DDMMYYYYhhmmss), xlsFile)
	f, err := excelize.OpenFile(xlsFile)
	if err != nil {
		fmt.Println(err)
		l.Println(time.Now().In(loc).Format(DDMMYYYYhhmmss), err)
		return
	}
	fmt.Println("Open file complete!")
	l.Println(time.Now().In(loc).Format(DDMMYYYYhhmmss), "Open file complete!")
	// Get value from cell by given worksheet name and axis.
	rows, err := f.GetRows("Sheet1")
	
	if err != nil {
		fmt.Println(err)
		l.Println(time.Now().In(loc).Format(DDMMYYYYhhmmss), err)
		return
	}
	if err := os.MkdirAll(outputPath+"/"+runtime, 0777); err != nil {
		fmt.Println("Open file complete!2",runtime)
		l.Println(time.Now().In(loc).Format(DDMMYYYYhhmmss), err)
        l.Fatal(err)
    }
	if EndRow >(len(rows)-1) {
		fmt.Println("end > ",len(rows)-1)	
		return
	}
	useRows :=rows[StartRow:EndRow+1]
	fmt.Printf("Start generate from %s(%d) to %s(%d) \n",useRows[0][0],StartRow,useRows[len(useRows)-1][0],EndRow)
	l.Printf("%s Start generate from %s(%d) to %s(%d) \n",time.Now().In(loc).Format(DDMMYYYYhhmmss),useRows[0][0],StartRow,useRows[len(useRows)-1][0],EndRow)
	fmt.Printf("Totals %d records \n",len(useRows))
	l.Printf("%s Totals %d records \n",time.Now().In(loc).Format(DDMMYYYYhhmmss),len(useRows))
	bar:=InitBar(len(useRows))
	qrfile :=envs["QRCODE"]+"/qr-"+runtime+".png"
	for i := 0; i < len(useRows); i++ {
		//fmt.Println("pid=",useRows[i][0])
		go generate(i,useRows[i],qrfile,templatePath,outputPath,r,bar)
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
			l.Println(time.Now().In(loc).Format(DDMMYYYYhhmmss), err)
		}
		if _, err := os.Stat(qrfile); err ==nil {
			dir, err := os.Getwd()
			if err == nil {
				os.Remove(dir+"/"+qrfile)
			}
		}

	}()
}
