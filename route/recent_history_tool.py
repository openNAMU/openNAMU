from .tool.func import *

async def recent_history_tool(name = 'Test', rev = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        num = str(rev)

        data = '' + \
            '<h2>' + get_lang(conn, 'tool') + '</h2>' + \
            '<ul>' + \
                '<li><a href="/raw_rev/' + num + '/' + url_pas(name) + '">' + get_lang(conn, 'raw') + '</a></li>' + \
        ''

        data += '<li><a href="/revert/' + num + '/' + url_pas(name) + '">' + get_lang(conn, 'revert') + ' (r' + num + ')</a></li>'
        if rev - 1 > 0:
            data += '<li><a href="/revert/' + str(rev - 1) + '/' + url_pas(name) + '">' + get_lang(conn, 'revert') + ' (r' + str(rev - 1) + ')</a></li>'

        if rev - 1 > 0:
            data += '<li><a href="/diff/' + str(rev - 1) + '/' + num + '/' + url_pas(name) + '">' + get_lang(conn, 'compare') + '</a></li>'

        data += '<li><a href="/history/' + url_pas(name) + '">' + get_lang(conn, 'history') + '</a></li>'
        data += '</ul>'

        if await acl_check(tool = 'hidel_auth') != 1:
            data += '<h3>' + get_lang(conn, 'admin') + '</h3>'
            data += '<ul>'
            curs.execute(db_change('select title from history where title = ? and id = ? and hide = "O"'), [name, num])
            data += '<li><a href="/history_hidden/' + num + '/' + url_pas(name) + '">'
            if curs.fetchall():
                data += get_lang(conn, 'hide_release') 
            else:
                data += get_lang(conn, 'hide')

            data += '</a></li>'
            data += '</ul>'

        if await acl_check('', 'owner_auth', '', '') != 1:
            data += '<h3>' + get_lang(conn, 'owner') + '</h3>'
            data += '<ul>'
            data += '<li><a href="/history_delete/' + num + '/' + url_pas(name) + '">' + get_lang(conn, 'history_delete') + '</a></li>'
            data += '<li><a href="/history_send/' + num + '/' + url_pas(name) + '">' + get_lang(conn, 'send_edit') + '</a></li>'
            data += '</ul>'

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(r' + num + ')', 0])],
            data = data,
            menu = [['history/' + url_pas(name), get_lang(conn, 'return')]]
        ))