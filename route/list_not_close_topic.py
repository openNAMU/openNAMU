from .tool.func import *
import pymysql

def list_not_close_topic_2(conn):
    curs = conn.cursor()

    div = '<ul>'
    
    curs.execute('select title, sub from rd where stop != "O" order by date desc')
    n_list = curs.fetchall()
    for data in n_list:
        div += '<li><a href="/topic/' + url_pas(data[0]) + '/sub/' + url_pas(data[1]) + '">' + html.escape(data[0]) + ' (' + data[1] + ')</a></li>'
            
    div += '</ul>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('open_discussion_list'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['manager', load_lang('return')]]
    ))