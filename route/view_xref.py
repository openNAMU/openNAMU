from .tool.func import *
import pymysql

def view_xref_2(conn, name):
    curs = conn.cursor()

    num = int(number_check(flask.request.args.get('num', '1')))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0
        
    div = '<ul>'
    
    curs.execute("select link, type from back where title = %s and not type = 'cat' and not type = 'no' order by link asc limit %s, '50'", [name, str(sql_num)])
    data_list = curs.fetchall()
    for data in data_list:
        div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a>'
        
        if data[1]:                
            div += ' (' + data[1] + ')'
        
        curs.execute("select title from back where title = %s and type = 'include'", [data[0]])
        db_data = curs.fetchall()
        if db_data:
            div += ' <a id="inside" href="/xref/' + url_pas(data[0]) + '">(' + load_lang('backlink') + ')</a>'

        div += '</li>'
      
    div += '</ul>' + next_fix('/xref/' + url_pas(name) + '?num=', num, data_list)
    
    return easy_minify(flask.render_template(skin_check(), 
        imp = [name, wiki_set(), custom(), other2([' (' + load_lang('backlink') + ')', 0])],
        data = div,
        menu = [['w/' + url_pas(name), load_lang('return')]]
    ))