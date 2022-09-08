package main

import (
	"fmt"
	"flag"
	"time"
	"log"
	"os"
	"os/exec"
	rt "runtime"
	"github.com/schollz/progressbar/v3"
	qrcode "github.com/skip2/go-qrcode"
	lo "go-pdf/log"
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
