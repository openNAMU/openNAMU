from .tool.func import *

def list_give_2(conn):
    curs = conn.cursor()

    list_data = '<ul>'
    back = ''

    curs.execute("select distinct name from alist order by name asc")
    for data in curs.fetchall():                      
        if back != data[0]:
            back = data[0]

        list_data += '<li><a href="/admin_plus/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'
    
    list_data += '</ul><hr class=\"main_hr\"><a href="/manager/8">(' + load_lang('add') + ')</a>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('admin_group_list'), wiki_set(), custom(), other2([0, 0])],
        data = list_data,
        menu = [['manager', load_lang('return')]]
    ))    