package pdfGenerator

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	"strings"
	"bufio"
	"github.com/google/uuid"
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
func (r *RequestPdf) GeneratePDF(pdfPath string,qrfile string,runtime string) (bool, error) {
	ck :=true
	var e error 
	// write whole the body
	u :=uuid.New()
	id := strings.Replace(u.String(), "-", "", -1)
	file := envs["TEMPDIR"]+ "/" + id +"-"+runtime+ ".html"
	// fmt.Println(r.body);
	fi, _ := os.Create(file)
	os.Chmod(file,os.FileMode(0644))
	// fi.Sync()
	w := bufio.NewWriter(fi)
	 n, _ := w.WriteString(r.body)
	w.Flush()
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("error=",err,n)
		e=err
		ck=false
	}
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		ck=false
		e=err
	}
	page :=wkhtmltopdf.NewPage(file)
	workingDir, err := os.Getwd()
	if err != nil {
		ck=false
		e=err
	}
	page.Allow.Set(workingDir)
	page.EnableLocalFileAccess.Set(true)
	js,_:=strconv.ParseUint(envs["JAVASCRIPTDELAY"],10,64)
	page.JavascriptDelay.Set(uint(js))
	pdfg.AddPage(page)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	dpi,_:=strconv.ParseUint(envs["DPI"],10,64)
	pdfg.Dpi.Set(uint(dpi))
	imgDpi,_:=strconv.ParseUint(envs["IMAGEDPI"],10,64)
	pdfg.ImageDpi.Set(uint(imgDpi))
	imageQuality,_:=strconv.ParseUint(envs["IMAGEQUALITY"],10,64)
	pdfg.ImageQuality.Set(uint(imageQuality))
	lowQuality,_ :=strconv.ParseBool(envs["LOWQUALITY"])
	pdfg.LowQuality.Set(lowQuality)
	left,_ :=strconv.ParseUint(envs["MARGINLEFT"],10,64)
	pdfg.MarginLeft.Set(uint(left))
	right,_ :=strconv.ParseUint(envs["MARGINRIGHT"],10,64)
	pdfg.MarginRight.Set(uint(right))
	top,_ :=strconv.ParseUint(envs["MARGINTOP"],10,64)
	pdfg.MarginTop.Set(uint(top))
	bottom,_ :=strconv.ParseUint(envs["MARGINBOTTOM"],10,64)
	pdfg.MarginBottom.Set(uint(bottom))
	err = pdfg.Create()
	if err != nil {
		ck=false
		e=err
	}
	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		ck=false
	}
	dir, err := os.Getwd()
	if err != nil {
		ck=false
		e=err
	}
	defer os.Remove(dir + "/"+file)
	defer os.Remove(dir + "/"+qrfile)
	defer fi.Close()
	if ck {
		e=nil
	}
	return ck, e
}