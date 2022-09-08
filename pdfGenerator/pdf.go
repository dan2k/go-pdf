package pdfGenerator

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"

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

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

//generate pdf function
func (r *RequestPdf) GeneratePDF(pdfPath string) (bool, error) {
	
	// t := time.Now().Unix()
	// write whole the body
	u :=uuid.New()
	id := strings.Replace(u.String(), "-", "", -1)
	// file := envs["TEMPDIR"] + "/" + strconv.FormatInt(int64(t), 10) + ".html"
	file := envs["TEMPDIR"] + "/" + id + ".html"
	// fmt.Println(r.body);
	err1 := ioutil.WriteFile(file, []byte(r.body), 0644)
	if err1 != nil {
		panic(err1)
	}
	
	if _, err := os.Stat(file); os.IsNotExist(err) {
		l.Fatal(err)
	}
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		l.Fatal(err)
	}
	page :=wkhtmltopdf.NewPage(file)
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	page.Allow.Set(workingDir)
	page.EnableLocalFileAccess.Set(true)
	// page.JavascriptDelay.Set(700)
	js,_:=strconv.ParseUint(envs["JAVASCRIPTDELAY"],10,64)
	page.JavascriptDelay.Set(uint(js))
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
		l.Fatal(err)
	}
	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		l.Fatal(err)
	}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	defer os.Remove(dir + "/"+file)
	return true, nil
}

