from .tool.func import *

def list_image_file_2(conn):
    curs = conn.cursor()

    num = int(number_check(flask.request.args.get('num', '1')))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0

    list_data = '<ul class="inside_ul">'
    back = ''

    curs.execute(db_change("select title from data where title like 'file:%' limit ?, 50"), [sql_num])
    data_list = curs.fetchall()
    for data in data_list:
        list_data += '<li><a href="/w/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a></li>'

    list_data += next_fix('/image_file_list?num=', num, data_list)

    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('image_file_list'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
        data = list_data,
        menu = [['other', load_lang('return')]]
    ))    
