package app

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/divan/num2words"
	"github.com/pkg/errors"
)

//ParseTemplate parsing template function
func ParseTemplate(templateFileName string, data interface{}) (*bytes.Buffer, error) {
	funcMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"inc": func(i int) int {
			return i + 1
		},
		"div": func(i *float32) float32 {
			return *i / 2
		},
		"inWords": func(i float64) string {
			return strings.ToUpper(num2words.Convert(int(i)))
		},
	}
	fName := templateFileName[strings.LastIndex(templateFileName, "/")+1:]
	t, err := template.New(fName).Funcs(funcMap).ParseFiles(templateFileName)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf, nil
}

//GeneratePDF generate pdf function
func GeneratePDF(body *bytes.Buffer) (*bytes.Buffer, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create PDFGenerator instance")
	}
	pdfg.AddPage(wkhtmltopdf.NewPageReader(body))
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Dpi.Set(300)
	err = pdfg.Create()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pdf document")
	}
	buf := pdfg.Buffer()
	return buf, nil
}
