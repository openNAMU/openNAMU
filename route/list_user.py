from .tool.func import *

def list_user(arg_num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        sql_num = (arg_num * 50 - 50) if arg_num * 50 > 0 else 0

        list_data = '<ul>'

        curs.execute(db_change("select id, data from user_set where name = 'date' order by data desc limit ?, 50"), [sql_num])
        user_list = curs.fetchall()
        for data in user_list:
            list_data += '<li>'
            list_data += ip_pas(data[0])
            list_data += ' | ' + data[1] if data[1] != '' else ''
            list_data += '</li>'

        list_data += '</ul>' + next_fix(conn, '/list/user/', arg_num, user_list)

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'member_list'), wiki_set(conn), wiki_custom(conn), wiki_css([0, 0])],
            data = list_data,
            menu = [['other', get_lang(conn, 'return')]]
        ))