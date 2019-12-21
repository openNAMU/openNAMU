from .tool.func import *

def list_not_close_topic_2(conn):
    curs = conn.cursor()

    div = '<ul>'

    curs.execute(db_change('select title, sub, date from rd where stop != "O" order by date desc'))
    n_list = curs.fetchall()
    for data in n_list:
        curs.execute(db_change("select code from topic where id = '1' and title = ? and sub = ?"), [data[0], data[1]])
        div += '<li><a href="/thread/' + url_pas(curs.fetchall()[0][0]) + '">' + html.escape(data[0]) + '</a> (' + data[1] + ') | ' + data[2] + '</li>'

    div += '</ul>'

    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('open_discussion_list'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['manager', load_lang('return')]]
    ))