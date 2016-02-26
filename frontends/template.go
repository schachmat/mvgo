package frontends

import (
	"bufio"
	"flag"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/schachmat/mvgo/iface"
)

type tmplConfig struct {
	tFileName string
	oFileName string
	message   string
}

func (c *tmplConfig) Setup() {
	flag.StringVar(&c.tFileName, "tmpl-template", "frontends/template.html", "template frontend: `FILENAME` of the template")
	flag.StringVar(&c.oFileName, "tmpl-output", "index.html", "template frontend: `FILENAME` of the output file")
	flag.StringVar(&c.message, "tmpl-message", "", "template frontend: `MESSAGE` to be displayed")
}

func (c *tmplConfig) RenderDepartures(station string, deps []iface.Departure) {
	t, err := template.ParseFiles(c.tFileName)
	if err != nil {
		log.Fatalf("template frontend: Unable to parse template %s:\n%v", c.tFileName, err)
	}

	fout, err := os.Create(c.oFileName)
	if err != nil {
		log.Fatalf("template frontend: Could not open %s for writing: %v", c.oFileName, err)
	}
	defer fout.Close()

	writer := bufio.NewWriter(fout)
	defer writer.Flush()

	if err := t.Execute(writer, struct {
		Station string
		Data    []iface.Departure
		Clock   string
		Message string
	}{station, deps, time.Now().Format("15:04"), c.message}); err != nil {
		log.Fatalln("template frontend: Error executing template: %v", err)
	}
}

func init() {
	iface.AllFrontends["template"] = &tmplConfig{}
}
