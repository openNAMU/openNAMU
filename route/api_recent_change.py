from .tool.func import *

def api_recent_change_2(conn):
    curs = conn.cursor()

    num = int(number_check(flask.request.args.get('num', '10')))
    num = 50 if (1 if not num > 0 else num) > 50 else num

    data_list = []
    curs.execute(db_change('select id, title from rc order by date desc'))
    for i in curs.fetchall():
        curs.execute(db_change('select id, title, date, ip, send, leng, hide from history where id = ? and title = ?'), i)
        data_list += curs.fetchall()
        
    if data_list:
        return flask.jsonify(data_list)
    else:
        return flask.jsonify({})