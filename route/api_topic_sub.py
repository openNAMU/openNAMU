from .tool.func import *

def api_topic_sub_2(conn, name, sub, time):
    curs = conn.cursor()

    if flask.request.args.get('time', None):
        curs.execute("select id, data, ip from topic where title = ? and sub = ? and date >= ? order by id + 0 asc", [name, sub, flask.request.args.get('time', None)])
    else:
        curs.execute("select id, data, ip from topic where title = ? and sub = ? order by id + 0 asc", [name, sub])
    data = curs.fetchall()
    if data:
        json_data = {}
        for i in data:
            json_data[i[0]] =   {
                "data" : i[1],
                "id" : i[2]
            }

        return flask.jsonify(json_data)
    else:
        return flask.jsonify({})