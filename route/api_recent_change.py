from .tool.func import *

def api_recent_change_2(conn):
    curs = conn.cursor()

    num = int(number_check(flask.request.args.get('num', '10')))
    num = 50 if (1 if not num > 0 else num) > 50 else num

    curs.execute(db_change('' + \
        'select id, title, date, ip, send, leng from history ' + \
        "where not title like 'user:%' " + \
        'order by date desc ' + \
        'limit ?' + \
    ''), [num])
    all_list = curs.fetchall()
    if all_list:
        return flask.jsonify(all_list)
    else:
        return flask.jsonify({})