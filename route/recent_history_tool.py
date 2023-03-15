from .tool.func import *

def recent_history_tool(name = 'Test', rev = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        num = str(rev)

        data = '' + \
            '<h2>' + load_lang('tool') + '</h2>' + \
            '<ul class="opennamu_ul">' + \
                '<li><a href="/raw_rev/' + num + '/' + url_pas(name) + '">' + load_lang('raw') + '</a></li>' + \
        ''

        data += '<li><a href="/revert/' + num + '/' + url_pas(name) + '">' + load_lang('revert') + ' (r' + num + ')</a></li>'
        if rev - 1 > 0:
            data += '<li><a href="/revert/' + str(rev - 1) + '/' + url_pas(name) + '">' + load_lang('revert') + ' (r' + str(rev - 1) + ')</a></li>'

        if rev - 1 > 0:
            data += '<li><a href="/diff/' + str(rev - 1) + '/' + num + '/' + url_pas(name) + '">' + load_lang('compare') + '</a></li>'

        data += '<li><a href="/history/' + url_pas(name) + '">' + load_lang('history') + '</a></li>'
        data += '</ul>'

        if admin_check(6) == 1:
            data += '<h3>' + load_lang('admin') + '</h3>'
            data += '<ul class="opennamu_ul">'
            curs.execute(db_change('' + \
                'select title from history ' + \
                'where title = ? and id = ? and hide = "O"' + \
            ''), [name, num])
            data += '<li><a href="/history_hidden/' + num + '/' + url_pas(name) + '">'
            if curs.fetchall():
                data += load_lang('hide_release') 
            else:
                data += load_lang('hide')

            data += '</a></li>'
            data += '</ul>'

        if admin_check() == 1:
            data += '<h3>' + load_lang('owner') + '</h3>'
            data += '<ul class="opennamu_ul">'
            data += '<li><a href="/history_delete/' + num + '/' + url_pas(name) + '">' + load_lang('history_delete') + '</li>'
            data += '<li><a href="/history_send/' + num + '/' + url_pas(name) + '">' + load_lang('send_edit') + '</li>'
            data += '</ul>'

        return easy_minify(flask.render_template(skin_check(),
            imp = [name, wiki_set(), wiki_custom(), wiki_css(['(r' + num + ')', 0])],
            data = data,
            menu = [['history/' + url_pas(name), load_lang('return')]]
        ))