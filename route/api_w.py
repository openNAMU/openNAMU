from .tool.func import *

def api_w_2(conn, name):
    curs = conn.cursor()

    if flask.request.method == 'POST':
        json_data = { "title" : name, "data" : render_set(title = name, data = flask.request.form.get('data', '')) }
        
        return flask.jsonify(json_data)
    else:
        curs.execute("select data from data where title = ?", [name])
        data = curs.fetchall()
        if data:
            json_data = { "title" : name, "data" : render_set(title = name, data = data[0][0]) }
        
            return flask.jsonify(json_data)
        else:
            return flask.jsonify({})