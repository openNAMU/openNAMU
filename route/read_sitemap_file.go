package route

import (
	"os"
	"strconv"
	"strings"
)

func Read_sitemap_file(name string) (string, bool) {
	valid := name == "sitemap.xml"
	if strings.HasPrefix(name, "sitemap_") && strings.HasSuffix(name, ".xml") {
		middle := strings.TrimSuffix(strings.TrimPrefix(name, "sitemap_"), ".xml")
		if _, err := strconv.Atoi(middle); err == nil {
			valid = true
		}
	}
	if !valid {
		return "", false
	}

	data, err := os.ReadFile(name)
	if err != nil {
		return "", false
	}

	return string(data), true
}
