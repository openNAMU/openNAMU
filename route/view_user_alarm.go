package route

import "database/sql"

func challenge_is_complete(db *sql.DB, id string, name string) bool {
	return user_value(db, id, name) != ""
}

func challenge_design(image string, title string, info string, complete bool) string {
	border := "red"
	if complete {
		border = "green"
	}

	return `<table id="main_table_set" style="border: 2px solid ` + border + `">
		<tr>
			<td id="main_table_width_quarter" rowspan="2"><span style="font-size: 64px;">` + image + `</span></td>
			<td><span style="font-size: 32px;">` + title + `</span></td>
		</tr>
		<tr><td>` + info + `</td></tr>
	</table>
	<hr class="main_hr">`
}
