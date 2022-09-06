package pdfGenerator

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

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
	
	t := time.Now().Unix()
	// write whole the body
	file:=envs["TEMPDIR"]+"/" + strconv.FormatInt(int64(t), 10) + ".html";
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
	pdfg.AddPage(page)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Dpi.Set(70)
	pdfg.ImageDpi.Set(100)
	pdfg.ImageQuality.Set(50)
	pdfg.LowQuality.Set(true)
	//pdfg.MarginLeft.Set(100)
	// pdfg.MarginRight.Set(100)
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

