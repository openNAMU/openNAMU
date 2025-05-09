from .tool.func import *

async def view_down(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        div = '<ul>'

        curs.execute(db_change("select title from data where title like ?"), [name + '/%'])
        for data in curs.fetchall():
            div += '<li><a href="/w/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a></li>'

        div += '</ul>'

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'sub') + ')', 0])],
            data = div,
            menu = [['w/' + url_pas(name), get_lang(conn, 'return')]]
        ))