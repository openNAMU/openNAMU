from .tool.func import *

def api_recent_discuss(num = 10, get_type = 'normal'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        num = 50 if num > 50 else num
        data_list = []

        if get_type == 'stop':
            curs.execute(db_change("select title, sub, date, code, stop from rd where stop = 'O' order by date desc limit ?"), [num])
        elif get_type == 'all':
            curs.execute(db_change("select title, sub, date, code, stop from rd order by date desc limit ?"), [num])
        else:
            curs.execute(db_change("select title, sub, date, code, stop from rd where not stop = 'O' order by date desc limit ?"), [num])

        for i in curs.fetchall():
            data_list += [i]

        return flask.jsonify(data_list if data_list else {}) 