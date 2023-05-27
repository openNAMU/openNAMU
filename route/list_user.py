from .tool.func import *

def list_user(arg_num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        sql_num = (arg_num * 50 - 50) if arg_num * 50 > 0 else 0

        list_data = '<ul class="opennamu_ul">'

        curs.execute(db_change("select id, data from user_set where name = 'date' order by data desc limit ?, 50"), [sql_num])
        user_list = curs.fetchall()
        for data in user_list:
            list_data += '' + \
                '<li>' + \
                    ip_pas(data[0]) + (' (' + data[1] + ')' if data[1] != '' else '') + \
                '</li>' + \
            ''

        list_data += '</ul>' + next_fix('/list/user/', arg_num, user_list)

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('member_list'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = list_data,
            menu = [['other', load_lang('return')]]
        ))