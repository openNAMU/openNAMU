from .tool.func import *

async def list_admin():
    with get_db_connect() as conn:
        curs = conn.cursor()

        div = '<ul>'

        curs.execute(db_change(
            "select id, data from user_set where name = 'acl' and not data = 'user'"
        ))
        for data in curs.fetchall():
            name = '' + \
                await ip_pas(data[0]) + ' ' + \
                '<a href="/auth/list/add/' + url_pas(data[1]) + '">(' + data[1] + ')</a>' + \
            ''

            div += '<li>' + name + '</li>'

        div += '</ul>'

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'admin_list'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
            data = div,
            menu = [['other', get_lang(conn, 'return')]]
        ))