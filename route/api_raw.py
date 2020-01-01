from .tool.func import *

def api_raw_2(conn, name):
    curs = conn.cursor()

    if acl_check(name, 'render') != 1:
        curs.execute(db_change("select data from data where title = ?"), [name])
        data = curs.fetchall()
        if data:
            json_data = { "title" : name, "data" : render_set(title = name, data = data[0][0], s_data = 1) }

            return flask.jsonify(json_data)
        else:
            return flask.jsonify({})
    else:
        return flask.jsonify({})