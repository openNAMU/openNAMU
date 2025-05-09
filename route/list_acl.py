from .tool.func import *

async def list_acl(arg_num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        sql_num = (arg_num * 50 - 50) if arg_num * 50 > 0 else 0

        div = '<ul>'

        curs.execute(db_change(
            "select distinct title, data, type from acl where data != '' and not title like 'user:%' order by title desc limit ?, 50"
        ), [sql_num])
        list_data = curs.fetchall()
        for data in list_data:
            curs.execute(db_change("select time from re_admin where what like ? order by time desc limit 1"), ['acl (' + data[0] + ')%'])
            time_data = curs.fetchall()
            time_data = (time_data[0][0] + ' | ') if time_data else ''

            curs.execute(db_change("select data from acl where title = ? and type = 'why'"), [data[0]])
            why_data = curs.fetchall()
            why_data = (' | ' + why_data[0][0]) if why_data and why_data[0][0] != '' else ''

            div += '' + \
                '<li>' + \
                    time_data + \
                    '<a href="/acl/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a>' + \
                    why_data + \
                '</li>' + \
            ''

        div += '</ul>'
        div += get_next_page_bottom(conn, '/list/document/acl/{}', arg_num, list_data)

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'acl_document_list'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
            data = div,
            menu = [['other', get_lang(conn, 'return')]]
        ))