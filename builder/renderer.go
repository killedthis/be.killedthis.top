package builder

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
)

type Renderer struct {
	ServiceProvider string
	Shows           []KilledShow
	OtherServices   []string
	template        *template.Template
}

type templateFields struct {
	Title string
	Sites []string
}

func NewRenderer(sp string, otherServices []string, shows []KilledShow) *Renderer {
	r := Renderer{
		ServiceProvider: sp,
		OtherServices:   otherServices,
		Shows:           shows,
	}

	r.init()

	return &r
}

func (m *Renderer) init() {
	siteTemplate, err := template.ParseFiles(
		"gotemplates/site.html",
		"gotemplates/menu.html",
	)
	if err != nil {
		log.Panic(err)
	}

	m.template = siteTemplate
}

func (m *Renderer) RenderHtml() {
	outputFolder := os.Getenv("OUTPUT")

	if outputFolder == "" {
		log.Panic("unknown output folder, specify ENV 'OUTPUT', should probably go into a config file later")
		return
	}

	file, err := os.OpenFile(outputFolder+"/"+strings.ToLower(m.ServiceProvider)+".html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Panic("failed to open output file: ", err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println("failed to close file: ", err)
		}
	}()

	err = m.template.Execute(file, templateFields{
		Title: fmt.Sprintf("%s killed this", m.ServiceProvider),
		Sites: m.OtherServices,
	})

	if err != nil {
		log.Panic(err)
	}
}
