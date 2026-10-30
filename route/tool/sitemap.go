package tool

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func Make_sitemap(db *sql.DB) (int, int, error) {
	domain := ""

	sitemap_auto_exclude_domain := Get_setting_value(db, "sitemap_auto_exclude_domain", "", "")
	if sitemap_auto_exclude_domain == "" {
		domain = Get_domain(db, true)
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
		return 0, 0, fmt.Errorf("file write error : sitemap.xml")
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
			return 0, 0, fmt.Errorf("file write error : %s", file_name)
		}
	}

	return count, len(titles), nil
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
