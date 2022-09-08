package main

import (
	"fmt"
	"flag"
	"time"
	"log"
	"os"
	"os/exec"
	"strings"
	"strconv"
	rt "runtime"
	"github.com/schollz/progressbar/v3"
	qrcode "github.com/skip2/go-qrcode"
	lo "go-pdf/log"
	u "go-pdf/pdfGenerator"
)
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
	MX			int
	MY			int
	Media       string
}
var StartRow,EndRow int
var l *log.Logger
var envs map[string]string
var clear map[string]func() //create a map for storing clear funcs
func InitBar(totals int) *progressbar.ProgressBar{
	return progressbar.NewOptions(
		totals,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionShowCount(),
		progressbar.OptionFullWidth(),
		// progressbar.OptionSetVisibility(true),
		// progressbar.OptionSetWidth(500),
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
			l.Println(time.Now().In(loc).Format(DDMMYYYYhhmmss), "complete !")
		}),
	)
}
func GenQr(data string,file string) error{
	err := qrcode.WriteFile(data, qrcode.Medium, 256, file)
	return err
}
func InitFlag(){
	// EndRow :=flag.Int("end",0,"an int ")
	// StartRow := flag.Int("start", 0, "an int")
	flag.IntVar(&StartRow, "start", 1, "a string var")
	flag.IntVar(&EndRow, "end", 1, "a string var")
	flag.Usage = func() {                                                  // [4]
		fmt.Fprintf(os.Stderr, "Options:\n-start int   number of records for start\n-end int number of records for end \nExample:\n./go-pdf -start=1 -end=10 \n")
	}
	
	flag.Parse()
	if flag.NFlag() !=2 { 
		flag.Usage()
		os.Exit(0)
	}
	
	if EndRow <StartRow {
		fmt.Println("start < end ")	
		os.Exit(0)
	}
	if StartRow <1 {
		fmt.Println("start < 0")	
		os.Exit(0)
	}
}
func InitLog(){
	lo.New(runtime)
	envs =lo.Envs
	l =lo.L
}
func init() {
    clear = make(map[string]func()) //Initialize it
    clear["linux"] = func() { 
        cmd := exec.Command("clear") //Linux example, its tested
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
    clear["windows"] = func() {
        cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested 
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
}

func CallClear() {
    value, ok := clear[rt.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
    if ok { //if we defined a clear func for that platform:
        value()  //we execute it
    } else { //unsupported platform
        panic("Your platform is unsupported! I can't clear terminal screen :(")
    }
}

func generate(number int,row []string,qrfile string,templatePath string,outputPath string,r *u.RequestPdf,bar *progressbar.ProgressBar) error{
	pid:=row[0]
	
	t:=strings.Split(qrfile, "/")
	t2:=strings.Join(t[1:int(len(t))],"/")
	if err :=GenQr(pid,qrfile);err != nil {
		fmt.Println(err)
	}
	mx,_ :=strconv.Atoi(envs["MX"])
	my,_ :=strconv.Atoi(envs["MY"])
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
			MX:	         mx,
			MY:	         my,
			Media:       envs["MEDIA"],
	}
	if err := r.ParseTemplate(templatePath, tmp); err == nil {
		r.GeneratePDF(outputPath+"/"+runtime+"/"+pid+".pdf")
		l.Println(time.Now().In(loc).Format(DDMMYYYYhhmmss), "PID", pid)
	} else {
			fmt.Println(err)
			l.Println(time.Now().In(loc).Format(DDMMYYYYhhmmss), err)
	}
	bar.Add(1)
	return nil
}