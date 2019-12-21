from .tool.func import *

def alarm_2(conn):
    curs = conn.cursor()

    num = int(number_check(flask.request.args.get('num', '1')))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0

    data = '<ul>'

    curs.execute(db_change("select data, date from alarm where name = ? order by date desc limit ?, 50"), [ip_check(), sql_num])
    data_list = curs.fetchall()
    if data_list:
        data = '<a href="/del_alarm">(' + load_lang('delete') + ')</a><hr class=\"main_hr\">' + data

        for data_one in data_list:
            data += '<li>' + data_one[0] + ' (' + data_one[1] + ')</li>'

    data += '</ul>' + next_fix('/alarm?num=', num, data_list)

    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('notice'), wiki_set(), custom(), other2([0, 0])],
        data = data,
        menu = [['user', load_lang('return')]]
    ))