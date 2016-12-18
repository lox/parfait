package args

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

// TemplateSourceValue is a Kingpin type for either a local file or a remote url
type TemplateSourceValue bytes.Buffer

func (h *TemplateSourceValue) Set(value string) error {
	log.Printf("Read from %s", value)
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return readURL(value, (*bytes.Buffer)(h))
	}
	f, err := os.Open(value)
	if err == nil {
		defer f.Close()
		io.Copy((*bytes.Buffer)(h), f)
	}
	return err
}

func readURL(u string, w io.Writer) error {
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Response status was %s", resp.Status)
	}

	_, err = io.Copy(w, resp.Body)
	return err
}

func (h TemplateSourceValue) String() string {
	return "[bytes]"
}

func TemplateSource(s kingpin.Settings) (target *bytes.Buffer) {
	target = &bytes.Buffer{}
	s.SetValue((*TemplateSourceValue)(target))
	return
}
