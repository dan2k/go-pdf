package pdfGenerator

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"time"

	"strconv"
	"strings"

	// "time"
	"github.com/google/uuid"
	// "fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

//pdf requestpdf struct
type RequestPdf struct {
	body string
}
var l *log.Logger
var envs map[string]string 
//new request to pdf function
func NewRequestPdf(body string,ll *log.Logger,e map[string]string) *RequestPdf {
	l=ll
	envs=e 
	return &RequestPdf{
		body: body,
	}
}

//parsing template function
func (r *RequestPdf) ParseTemplate(templateFileName string, data interface{}) error {
	ck:=true
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		ck=false
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, data);
	if  err != nil {
		ck=false
	}
	if !ck {
		l.Println(templateFileName)
		// recover()
	}
	r.body = buf.String()
	var e error 
	if !ck {
		e=err
	} else {
		e=nil
	}
	return e  
}

//generate pdf function
func (r *RequestPdf) GeneratePDF(pdfPath string,qrfile string) (bool, error) {
	ck :=true
	var e error 
	t := time.Now().Unix()
	// write whole the body
	u :=uuid.New()
	id := strings.Replace(u.String(), "-", "", -1)
	// file := envs["TEMPDIR"] + "/" + strconv.FormatInt(int64(t), 10) + ".html"
	file := envs["TEMPDIR"] + "/" + id +strconv.FormatInt(int64(t), 10)+ ".html"
	// fmt.Println(r.body);
	err1 := ioutil.WriteFile(file, []byte(r.body), 0644)
	if err1 != nil {
		// panic(err1)
		fmt.Println("error=",err1)
		// l.Fatalln("error=",err1)
		//l.Println(time.Now(),"error1:", err1)
		e=err1
		ck=false
		
	}
	
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("error=",err)
		// l.Fatalln("error=",err)
		//l.Println(time.Now(),"error2:", err1)
		e=err
		ck=false
	}
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		//l.Fatal(err)
		//l.Println(time.Now(),"error3:", err)
		ck=false
		e=err
	}
	page :=wkhtmltopdf.NewPage(file)
	workingDir, err := os.Getwd()
	if err != nil {
		//l.Println(time.Now(),"error4:", err)
		ck=false
		e=err
	}
	page.Allow.Set(workingDir)
	page.EnableLocalFileAccess.Set(true)
	
	// page.JavascriptDelay.Set(700)
	js,_:=strconv.ParseUint(envs["JAVASCRIPTDELAY"],10,64)
	page.JavascriptDelay.Set(uint(js))
	// page.LoadErrorHandling.Set("ignore")
	// page.LoadMediaErrorHandling.Set("ignore")
	
	pdfg.AddPage(page)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	// pdfg.Dpi.Set(100)
	dpi,_:=strconv.ParseUint(envs["DPI"],10,64)
	pdfg.Dpi.Set(uint(dpi))
	
	// pdfg.ImageDpi.Set(150)
	imgDpi,_:=strconv.ParseUint(envs["IMAGEDPI"],10,64)
	pdfg.ImageDpi.Set(uint(imgDpi))

	// pdfg.ImageQuality.Set(150)
	imageQuality,_:=strconv.ParseUint(envs["IMAGEQUALITY"],10,64)
	pdfg.ImageQuality.Set(uint(imageQuality))
	
	// pdfg.LowQuality.Set(true)
	lowQuality,_ :=strconv.ParseBool(envs["LOWQUALITY"])
	pdfg.LowQuality.Set(lowQuality)

	// pdfg.MarginLeft.Set(10)
	left,_ :=strconv.ParseUint(envs["MARGINLEFT"],10,64)
	pdfg.MarginLeft.Set(uint(left))

	// pdfg.MarginRight.Set(10)
	right,_ :=strconv.ParseUint(envs["MARGINRIGHT"],10,64)
	pdfg.MarginRight.Set(uint(right))

	// pdfg.MarginTop.Set(10)
	top,_ :=strconv.ParseUint(envs["MARGINTOP"],10,64)
	pdfg.MarginTop.Set(uint(top))

	// pdfg.MarginBottom.Set(10)
	bottom,_ :=strconv.ParseUint(envs["MARGINBOTTOM"],10,64)
	pdfg.MarginBottom.Set(uint(bottom))
	
	err = pdfg.Create()
	if err != nil {
		
		//l.Println(time.Now(),"error5:", err)
		ck=false
		e=err
	}
	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		//l.Println(time.Now(),"error6:", err)
		ck=false
	}
	dir, err := os.Getwd()
	if err != nil {
		//l.Println(time.Now(),"error7:", err)
		ck=false
		e=err
	}
	if !ck {
		dir, _ := os.Getwd()
		os.Remove(dir + "/"+file)
		os.Remove(dir + "/"+qrfile)
		
		// ck=false
		// l.Println(pdfPath)
		// recover()
	}
	defer os.Remove(dir + "/"+file)
	defer os.Remove(dir + "/"+qrfile)
	if ck {
		e=nil
	}
	return ck, e
}

