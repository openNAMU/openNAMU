from .tool.func import *
import pymysql

def api_raw_2(conn, name):
    curs = conn.cursor()

    curs.execute("select data from data where title = %s", [name])
    data = curs.fetchall()
    if data:
        json_data = { "title" : name, "data" : render_set(title = name, data = data[0][0], s_data = 1) }
    
        return flask.jsonify(json_data)
    else:
        return flask.jsonify({})