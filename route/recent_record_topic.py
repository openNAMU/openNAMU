from .tool.func import *

async def recent_record_topic(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        num = int(number_check(flask.request.args.get('num', '1')))
        sql_num = (num * 50 - 50) if num * 50 > 0 else 0

        div = '''
            <table id="main_table_set">
                <tr id="main_table_top_tr">
                    <td id="main_table_width">''' + get_lang(conn, 'discussion_name') + '''</td>
                    <td id="main_table_width">''' + get_lang(conn, 'writer') + '''</td>
                    <td id="main_table_width">''' + get_lang(conn, 'time') + '''</td>
                </tr>
        '''
        sub = '(' + html.escape(name) + ')'
        pas_name = await ip_pas(name)

        curs.execute(db_change("select code, id, date from topic where ip = ? order by date desc limit ?, 50"), [name, sql_num])
        data_list = curs.fetchall()
        for data in data_list:
            title = html.escape(data[0])

            curs.execute(db_change("select title, sub from rd where code = ?"), [data[0]])
            other_data = curs.fetchall()

            div += '' + \
                '<tr>' + \
                    '<td>' + \
                        '<a href="/thread/' + data[0] + '#' + data[1] + '">' + other_data[0][1] + '#' + data[1] + '</a> (' + other_data[0][0] + ')' + \
                    '</td>' + \
                    '<td>' + pas_name + '</td>' + \
                    '<td>' + data[2] + '</td>' + \
                '</tr>' + \
            ''

        div += '</table>'
        div += get_next_page_bottom(conn, '/record/topic/' + url_pas(name) + '?num={}', num, data_list)

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'discussion_record'), await wiki_set(), await wiki_custom(conn), wiki_css([sub, 0])],
            data = div,
            menu = [['other', get_lang(conn, 'other')], ['user/' + url_pas(name), get_lang(conn, 'user_tool')]]
        ))