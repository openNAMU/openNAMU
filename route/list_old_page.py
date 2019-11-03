from .tool.func import *

def list_old_page_2(conn):
    curs = conn.cursor()

    num = int(number_check(flask.request.args.get('num', '1')))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0

    div = '<ul>'
    
    curs.execute('' + \
        'select title, date from history ' + \
        "where title not like 'user:%' and title not like 'category:%' and title not like 'file:%' " + \
        'order by id desc, date asc ' + \
        'limit ?, "50"' + \
    '', [str(sql_num)])
    n_list = curs.fetchall()
    for data in n_list:
        div += '<li><a href="/w/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a> (' + re.sub(' .*$', '', data[1]) + ')</li>'
            
    div += '</ul>' + next_fix('/old_page?num=', num, n_list)

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('old_page'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['other', load_lang('return')]]
    ))