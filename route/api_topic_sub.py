from .tool.func import *

def api_topic_sub_2(conn, name, sub, time):
    curs = conn.cursor()

    if flask.request.args.get('num', None):
        if flask.request.args.get('over', '0') == '0':
            curs.execute("select id, data, date, ip from topic where title = ? and sub = ? and id + 0 = ? + 0 order by id + 0 asc", [name, sub, flask.request.args.get('num', '')])
        elif flask.request.args.get('over', '0') == '1':
            curs.execute("select id, data, date, ip from topic where title = ? and sub = ? and id + 0 >= ? + 0 order by id + 0 asc", [name, sub, flask.request.args.get('num', '')])
        else:
            curs.execute("select id, data, date, ip from topic where title = ? and sub = ? and id + 0 <= ? + 0 order by id + 0 asc", [name, sub, flask.request.args.get('num', '')])
    elif flask.request.args.get('time', None):
        if flask.request.args.get('over', '0') == '0':
            curs.execute("select id, data, date, ip from topic where title = ? and sub = ? and date = ? order by id + 0 asc", [name, sub, flask.request.args.get('time', '')])
        if flask.request.args.get('over', '0') == '1':
            curs.execute("select id, data, date, ip from topic where title = ? and sub = ? and date >= ? order by id + 0 asc", [name, sub, flask.request.args.get('time', '')])
        else:
            curs.execute("select id, data, date, ip from topic where title = ? and sub = ? and date <= ? order by id + 0 asc", [name, sub, flask.request.args.get('time', '')])
    else:
        curs.execute("select id, data, date, ip from topic where title = ? and sub = ? order by id + 0 asc", [name, sub])

    data = curs.fetchall()
    if data:
        json_data = {}
        for i in data:
            json_data[i[0]] =   {
                "data" : i[1],
                "date" : i[2],
                "id" : i[3]
            }

        return flask.jsonify(json_data)
    else:
        return flask.jsonify({}), 404