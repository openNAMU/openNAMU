package tool

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var Err_sitemap_domain_require = errors.New("sitemap domain require")
var Err_sitemap_file_write = errors.New("file write error")
var Err_sitemap_file_remove = errors.New("file remove error")

// Domain_setting_is_empty reports whether the full domain string is unusable
// because the domain setting is missing (empty or scheme only).
func Domain_setting_is_empty(domain string) bool {
	return domain == "" || domain == "http://" || domain == "https://"
}

func Make_sitemap(db *sql.DB) (int, int, error) {
	domain := ""

	sitemap_auto_exclude_domain := Get_setting_value(db, "sitemap_auto_exclude_domain", "", "")
	if sitemap_auto_exclude_domain == "" {
		domain = Get_domain(db, true)
		if Domain_setting_is_empty(domain) {
			return 0, 0, Err_sitemap_domain_require
		}
	}

	sql_add := ""

	sitemap_auto_exclude_user_page := Get_setting_value(db, "sitemap_auto_exclude_user_page", "", "")
	if sitemap_auto_exclude_user_page != "" {
		sql_add += " title not like 'user:%'"
	}

	sitemap_auto_exclude_file_page := Get_setting_value(db, "sitemap_auto_exclude_file_page", "", "")
	if sitemap_auto_exclude_file_page != "" {
		if sql_add != "" {
			sql_add += " and"
		}

		sql_add += " title not like 'file:%'"
	}

	sitemap_auto_exclude_category_page := Get_setting_value(db, "sitemap_auto_exclude_category_page", "", "")
	if sitemap_auto_exclude_category_page != "" {
		if sql_add != "" {
			sql_add += " and"
		}

		sql_add += " title not like 'category:%'"
	}

	if sql_add != "" {
		sql_add = " where" + sql_add
	}

	titles, err := Query_titles(db, "select title from data"+sql_add)
	if err != nil {
		return 0, 0, err
	}

	const sitemap_max_count = 30000

	count := (len(titles) + sitemap_max_count - 1) / sitemap_max_count
	if count == 0 {
		count = 1
	}

	data := strings.Builder{}
	data.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	data.WriteString("\n")
	data.WriteString(`<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	data.WriteString("\n")

	for i := 0; i < count; i++ {
		data.WriteString("<sitemap><loc>")
		data.WriteString(domain)
		data.WriteString("/sitemap_")
		data.WriteString(fmt.Sprint(i))
		data.WriteString(".xml</loc></sitemap>\n")
	}

	data.WriteString("</sitemapindex>")

	if err := Write_file("sitemap.xml", data.String(), 0644); err != nil {
		return 0, 0, fmt.Errorf("%w : sitemap.xml", Err_sitemap_file_write)
	}

	for i := 0; i < count; i++ {
		data.Reset()

		data.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
		data.WriteString("\n")
		data.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
		data.WriteString("\n")

		start := sitemap_max_count * i
		end := start + sitemap_max_count

		if end > len(titles) {
			end = len(titles)
		}

		for _, title := range titles[start:end] {
			data.WriteString("<url><loc>")
			data.WriteString(domain)
			data.WriteString("/w/")
			data.WriteString(Url_parser(title))
			data.WriteString("</loc></url>\n")
		}

		data.WriteString("</urlset>")

		file_name := "sitemap_" + fmt.Sprint(i) + ".xml"
		if err := Write_file(file_name, data.String(), 0644); err != nil {
			return 0, 0, fmt.Errorf("%w : %s", Err_sitemap_file_write, file_name)
		}
	}

	// The index and every current part are written, so parts left over from a
	// previous run with more documents can be removed safely now.
	if err := Remove_sitemap_stale_files(count); err != nil {
		return 0, 0, err
	}

	return count, len(titles), nil
}

// Remove_sitemap_stale_files removes sitemap_N.xml files whose index is not
// referenced by the current sitemap index (index >= count).
func Remove_sitemap_stale_files(count int) error {
	files, err := filepath.Glob("sitemap_*.xml")
	if err != nil {
		return fmt.Errorf("%w : %s", Err_sitemap_file_remove, "sitemap_*.xml")
	}

	for _, file_name := range files {
		middle := strings.TrimSuffix(strings.TrimPrefix(file_name, "sitemap_"), ".xml")
		index, err := strconv.Atoi(middle)
		if err != nil {
			continue
		}
		if index < count {
			continue
		}

		if err := os.Remove(file_name); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%w : %s", Err_sitemap_file_remove, file_name)
		}
	}

	return nil
}

func Query_titles(db *sql.DB, query string) ([]string, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	titles := []string{}

	for rows.Next() {
		title := ""

		if err := rows.Scan(&title); err != nil {
			return nil, err
		}

		titles = append(titles, title)
	}

	return titles, nil
}

func Write_file(path string, content string, mode os.FileMode) error {
	return os.WriteFile(path, []byte(content), mode)
}
