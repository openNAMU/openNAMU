from .tool.func import *

def list_image_file(arg_num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        sql_num = (arg_num * 50 - 50) if arg_num * 50 > 0 else 0

        list_data = '<ul class="opennamu_ul">'

        curs.execute(db_change("select title from data where title like 'file:%' limit ?, 50"), [sql_num])
        data_list = curs.fetchall()
        for data in data_list:
            list_data += '<li><a href="/w/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a></li>'

        list_data += next_fix('/list/file/', arg_num, data_list)

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('image_file_list'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = list_data,
            menu = [['other', load_lang('return')]]
        ))