from .tool.func import *

def list_admin():
    with get_db_connect() as conn:
        curs = conn.cursor()

        div = '<ul class="opennamu_ul">'

        curs.execute(db_change(
            "select id, data from user_set where name = 'acl' and not data = 'user'"
        ))
        for data in curs.fetchall():
            name = '' + \
                ip_pas(data[0]) + ' ' + \
                '<a href="/auth/list/add/' + url_pas(data[1]) + '">(' + data[1] + ')</a>' + \
            ''

            div += '<li>' + name + '</li>'

        div += '</ul>'

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('admin_list'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = div,
            menu = [['other', load_lang('return')]]
        ))