from .tool.func import *

def list_please(arg_num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        sql_num = (arg_num * 50 - 50) if arg_num * 50 > 0 else 0

        div = '<ul class="opennamu_ul">'

        curs.execute(db_change("select distinct title, link from back where type = 'no' limit ?, 50"), [sql_num])
        data_list = curs.fetchall()
        for data in data_list:
            div += '' + \
                '<li>' + \
                    '<a class="opennamu_not_exist_link" href="/w/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a> ' + \
                    '<a href="/w/' + url_pas(data[1]) + '">(' + html.escape(data[1]) + ')</a>' + \
                '</li>' + \
            ''

        div += '</ul>' + next_fix('/list/document/need/', arg_num, data_list)

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('need_document'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = div,
            menu = [['other', load_lang('return')]]
        ))