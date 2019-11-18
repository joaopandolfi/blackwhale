package utils

import (
	"crypto/md5"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/joaopandolfi/blackwhale/configurations"
)

func Round(value float64) float64 {
	return math.Round(value*100) / 100
}

func OnlyNumbers(value string) string {
	re := regexp.MustCompile("[0-9]+")
	numbers := re.FindAllString(value, -1)
	return strings.Join(numbers, "")
}

func GetStaticPageDir(page string) string {
	return configurations.Configuration.StaticPagesDir + page
}

func GetHbsPage(page string) *pongo2.Template {
	return pongo2.Must(pongo2.FromFile(GetStaticPageDir(page)))
}

func GenerateName() string {
	crutime := time.Now().Unix()
	h := md5.New()
	return fmt.Sprintf("%x-%d", h.Sum(nil), crutime)
}
