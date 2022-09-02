package pdfGenerator

import (
	"bytes"
	// "fmt"
	// "fmt"
	// "path/filepath"
	// "strings"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

//pdf requestpdf struct
type RequestPdf struct {
	body string
}

//new request to pdf function
func NewRequestPdf(body string) *RequestPdf {
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
	file:="templates/" + strconv.FormatInt(int64(t), 10) + ".html";
	
	err1 := ioutil.WriteFile(file, []byte(r.body), 0644)
	if err1 != nil {
		panic(err1)
	}
	
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Fatal(err)
	}
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}
	// page := wkhtmltopdf.NewPageReader(f)
	
	page :=wkhtmltopdf.NewPage(file)
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	page.Allow.Set(workingDir)
	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		log.Fatal(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	defer os.Remove(dir + "/"+file)
	return true, nil
}

