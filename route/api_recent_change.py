from .tool.func import *

def api_recent_change(num = 10):
    with get_db_connect() as conn:
        curs = conn.cursor()

        num = 50 if (1 if not num > 0 else num) > 50 else num
        repeat_ok = flask.request.args.get('repeat', '1')
        data_list = []
        admin = admin_check(6)
        get_title = ''

        curs.execute(db_change('select id, title from rc where type = "" order by date desc'))
        for i in curs.fetchall():
            if repeat_ok == '1' or i[1] != get_title:
                get_title = i[1]
                curs.execute(db_change('select id, title, date, ip, send, leng, hide from history where id = ? and title = ?'), i)
                get_data = curs.fetchall()
                if get_data:
                    if get_data[0][6] == '' or admin == 1:
                        data_list += get_data
                    else:
                        data_list += [['', '', '', '', '', '', get_data[0][6]]]
                else:
                    data_list += [['', '', '', '', '', '', '']]

        return flask.jsonify(data_list if data_list else {}) 