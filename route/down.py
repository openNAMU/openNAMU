from .tool.func import *

def down_2(conn, name):
    curs = conn.cursor()

    div = '<ul>'

    curs.execute("select title from data where title like ?", ['%' + name + '/%'])
    for data in curs.fetchall():
        div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'
        
    div += '</ul>'
    
    return easy_minify(flask.render_template(skin_check(), 
        imp = [name, wiki_set(), custom(), other2([' (' + load_lang('sub') + ')', 0])],
        data = div,
        menu = [['w/' + url_pas(name), load_lang('return')]]
    ))