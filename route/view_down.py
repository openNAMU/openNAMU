from .tool.func import *

async def view_down(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        div = '<ul>'

        curs.execute(db_change("select title from data where title like ?"), [name + '/%'])
        for data in curs.fetchall():
            div += '<li><a href="/w/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a></li>'

        div += '</ul>'

        return await render_template(
            name,
            div,
            '(' + await get_lang('sub') + ')',
            [['w/' + url_pas(name), await get_lang('return')]]
        )
