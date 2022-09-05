package main

import (
	"flag"
	"fmt"
	log "go-pdf/log"
	u "go-pdf/pdfGenerator"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/xuri/excelize/v2"
)

const (
	DDMMYYYYhhmmss = "2006-01-02 15:04:05"
	DDMMYYYYhhmmss2 = "20060102-15:04:05"
)

var envs map[string]string
var runtime =time.Now().UTC().Format(DDMMYYYYhhmmss2)
func main() {
	
	var startRow,endRow int 

	// endRow :=flag.Int("end",0,"an int ")
	// startRow := flag.Int("start", 0, "an int")
	// fmt.Println(*startRow,*endRow);
	flag.IntVar(&startRow, "start", 0, "a string var")
	flag.IntVar(&endRow, "end", 0, "a string var")
	flag.Parse()
	if endRow <startRow {
		fmt.Println("start < end ")	
		return
	}
	if startRow <0 {
		fmt.Println("start < 0 ")	
		return
	}
	//fmt.Println(startRow,endRow);
	log.New(runtime)
	envs :=log.Envs
	l :=log.L
	
	r := u.NewRequestPdf("",l,envs)

	//html template path
	templatePath := envs["TEMPDIR"] + "/" + envs["TEMPFILE"]
	//path for download pdf
	outputPath := envs["STORAGE"] 

	//html template data
	type templateData struct {
		Title       string
		Description string
		Company     string
		Contact     string
		Country     string
		Labels      []string
		Data        []int
		Qrcode      string 
		Pid         string
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
	
	// Get value from cell by given worksheet name and axis.

	rows, err := f.GetRows("Sheet1")

	if err != nil {
		fmt.Println(err)
		l.Println(time.Now().UTC().Format(DDMMYYYYhhmmss), err)
		return
	}
	fmt.Printf("Total %d rows \n", len(rows)-1)
	l.Printf("%s Total %d rows \n", time.Now().UTC().Format(DDMMYYYYhhmmss), len(rows)-1)
	
	if err := os.Mkdir(outputPath+"/"+runtime, os.ModePerm); err != nil {
		l.Println(time.Now().UTC().Format(DDMMYYYYhhmmss), err)
        l.Fatal(err)
    }
	if endRow >(len(rows)-1) {
		fmt.Println("end > ",len(rows)-1)	
		return
	}
	useRows :=rows[startRow:endRow+1]
	bar := progressbar.NewOptions(
		len(useRows),
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
	qrfile :=envs["QRCODE"]+"/qr-"+runtime+".png"
	for i := 0; i < len(useRows); i++ {
		pid:=useRows[i][0]
		t:=strings.Split(qrfile, "/")
		t2:=strings.Join(t[1:int(len(t))],"/")
		err := qrcode.WriteFile(pid, qrcode.Medium, 256, qrfile)
		if err != nil {
			fmt.Println(err)
		}
		tmp := templateData{
			Title:       "HTML to PDF generator",
			Description: "This is the simple HTML to PDF file.",
			Company:     "Jhon Lewis",
			Contact:     "Maria Anders",
			Country:     "Germany",
			Labels: 	 []string{"Red", "Blue", "Yellow", "Green", "Purple", "Orange"},
			Data:        []int{12, 19, 3, 5, 2, 3},
			Qrcode:      t2,
			Pid:		 pid,
		}
		if err := r.ParseTemplate(templatePath, tmp); err == nil {
			r.GeneratePDF(outputPath+"/"+runtime+"/"+pid+".pdf")
			l.Println(time.Now().UTC().Format(DDMMYYYYhhmmss), "PID", pid)
		} else {
			fmt.Println(err)
			l.Println(time.Now().UTC().Format(DDMMYYYYhhmmss), err)
		}
		bar.Add(1)
		// time.Sleep(1000 * time.Millisecond)
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
			l.Println(time.Now().UTC().Format(DDMMYYYYhhmmss), err)
		}
		if _, err := os.Stat(qrfile); err ==nil {
			dir, err := os.Getwd()
			if err == nil {
				os.Remove(dir+"/"+qrfile)
			}
		}

	}()

	

}
