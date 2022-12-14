package builder

import (
	"fmt"
	"html/template"
	"killedthis/shared"
	"log"
	"os"
	"sort"
	"strings"
)

type Renderer struct {
	ServiceProvider string
	Shows           []shared.KilledShow
	OtherServices   []string
	template        *template.Template
}

type templateFields struct {
	Title  string
	Sites  []string
	Years  []string
	Months []string
	Shows  []shared.KilledShow
}

func NewRenderer(sp string, otherServices []string, shows []shared.KilledShow) *Renderer {
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
		"gotemplates/custom_styles.html",
	)
	if err != nil {
		log.Panic(err)
	}

	m.template = siteTemplate
}

func (m *Renderer) RenderHtml(cfg *shared.ContentGeneratorConfig) {
	filename := fmt.Sprintf(cfg.FormatSiteFile, strings.ToLower(m.ServiceProvider))
	file, err := os.OpenFile(cfg.OutputDirectory+"/"+filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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

	years := make([]string, 0)
	months := make([]string, 0)

	// only show years, months we have data for
	for _, show := range m.Shows {
		year := fmt.Sprintf("%d", show.Year())
		if !contains(years, year) {
			years = append(years, year)
		}
		if !contains(months, show.Month()) {
			months = append(months, show.Month())
		}
	}

	// sort them
	sort.Strings(years)
	sort.Strings(months)

	err = m.template.Execute(file, templateFields{
		Title:  fmt.Sprintf("%s killed the following shows", m.ServiceProvider),
		Sites:  m.OtherServices,
		Shows:  m.Shows,
		Years:  years,
		Months: months,
	})

	if err != nil {
		log.Panic(err)
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
