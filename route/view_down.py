from .tool.func import *

def view_down(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        div = '<ul class="opennamu_ul">'

        curs.execute(db_change("select title from data where title like ?"), [name + '/%'])
        for data in curs.fetchall():
            div += '<li><a href="/w/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a></li>'

        div += '</ul>'

        return easy_minify(flask.render_template(skin_check(),
            imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('sub') + ')', 0])],
            data = div,
            menu = [['w/' + url_pas(name), load_lang('return')]]
        ))