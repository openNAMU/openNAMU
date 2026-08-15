package route

import (
	"fmt"
	"os"
	"strings"

	"opennamu/route/tool"
)

func Api_setting_sitemap_post(config tool.Config) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	domain := ""

	sitemap_auto_exclude_domain := ""
	tool.QueryRow_DB(
		db,
		"select data from other where name = 'sitemap_auto_exclude_domain'",
		[]any{&sitemap_auto_exclude_domain},
	)
	if sitemap_auto_exclude_domain == "" {
		domain = tool.Get_domain(db, true)
	}

	sql_add := ""

	sitemap_auto_exclude_user_page := ""
	tool.QueryRow_DB(
		db,
		"select data from other where name = 'sitemap_auto_exclude_user_page'",
		[]any{&sitemap_auto_exclude_user_page},
	)
	if sitemap_auto_exclude_user_page != "" {
		sql_add += " title not like 'user:%'"
	}

	sitemap_auto_exclude_file_page := ""
	tool.QueryRow_DB(
		db,
		"select data from other where name = 'sitemap_auto_exclude_file_page'",
		[]any{&sitemap_auto_exclude_file_page},
	)
	if sitemap_auto_exclude_file_page != "" {
		if sql_add != "" {
			sql_add += " and"
		}

		sql_add += " title not like 'file:%'"
	}

	sitemap_auto_exclude_category_page := ""
	tool.QueryRow_DB(
		db,
		"select data from other where name = 'sitemap_auto_exclude_category_page'",
		[]any{&sitemap_auto_exclude_category_page},
	)
	if sitemap_auto_exclude_category_page != "" {
		if sql_add != "" {
			sql_add += " and"
		}

		sql_add += " title not like 'category:%'"
	}

	if sql_add != "" {
		sql_add = " where" + sql_add
	}

	rows := tool.Query_DB(db, "select title from data"+sql_add)
	defer rows.Close()

	all_data := []string{}

	for rows.Next() {
		title := ""

		err := rows.Scan(&title)
		if err != nil {
			panic(err)
		}

		all_data = append(all_data, title)
	}

	const sitemap_max_count = 30000

	len_all_data := len(all_data)
	count := (len_all_data + sitemap_max_count - 1) / sitemap_max_count
	if count == 0 {
		count = 1
	}
	other_count := len_all_data % sitemap_max_count

	// 현재 로직에서는 직접 사용하지 않음
	_ = other_count

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

	err := os.WriteFile(
		"sitemap.xml",
		[]byte(data.String()),
		0644,
	)
	if err != nil {
		return_data["response"] = "error"
		return_data["error"] = "file write error : sitemap.xml"

		return return_data
	}

	for i := 0; i < count; i++ {
		data.Reset()

		data.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
		data.WriteString("\n")
		data.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
		data.WriteString("\n")

		start := sitemap_max_count * i
		end := start + sitemap_max_count

		if end > len_all_data {
			end = len_all_data
		}

		for _, x := range all_data[start:end] {
			data.WriteString("<url><loc>")
			data.WriteString(domain)
			data.WriteString("/w/")
			data.WriteString(tool.Url_parser(x))
			data.WriteString("</loc></url>\n")
		}

		data.WriteString("</urlset>")

		err = os.WriteFile(
			"sitemap_"+fmt.Sprint(i)+".xml",
			[]byte(data.String()),
			0644,
		)
		if err != nil {
			return_data["response"] = "error"
			return_data["error"] = "file write error : sitemap_" + fmt.Sprint(i) + ".xml"

			return return_data
		}
	}

	return_data["response"] = "ok"
	return_data["sitemap_count"] = count
	return_data["document_count"] = len_all_data

	return return_data
}
