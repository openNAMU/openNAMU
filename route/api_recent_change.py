from .tool.func import *

def api_recent_change(num = 10):
    with get_db_connect() as conn:
        curs = conn.cursor()

        num = 50 if (1 if not num > 0 else num) > 50 else num
        admin = admin_check(6)

        data_list = []

        curs.execute(db_change('select id, title from rc where type = "normal" order by date desc limit ?'), [num])
        for for_a in curs.fetchall():
            curs.execute(db_change('select id, title, date, ip, send, leng, hide from history where id = ? and title = ?'), for_a)
            db_data = curs.fetchall()
            if db_data:
                db_data = list(db_data[0])
                if db_data[6] == '' or admin == 1:
                    if admin == 1:
                        data_list += [db_data]
                    else:
                        db_data[3] = ip_pas(db_data[3], 1)
                        data_list += [db_data]
                else:
                    data_list += [['', '', '', '', '', '', db_data[6]]]
            else:
                data_list += [['', '', '', '', '', '', '']]

        return flask.jsonify(data_list if data_list else {}) 