from .tool.func import *

def api_w_2(conn, name):
    curs = conn.cursor()

    if flask.request.args.get('exist', None):
        curs.execute("select title from data where title = ?", [name])
        if curs.fetchall():
            return flask.jsonify({ "exist" : "1" })
        else:
            return flask.jsonify({})
    else:
        if acl_check(name, 'render') != 1:
            if flask.request.method == 'POST':
                g_data = render_set(title = name, data = flask.request.form.get('data', ''), num = 2)
                json_data = { "title" : name, "data" : g_data[0], "js_data" : g_data[1] }
                
                return flask.jsonify(json_data)
            else:
                curs.execute("select data from data where title = ?", [name])
                data = curs.fetchall()
                if data:
                    g_data = render_set(title = name, data = data[0][0], num = 2)
                    json_data = { "title" : name, "data" : g_data[0], "js_data" : g_data[1] }
                
                    return flask.jsonify(json_data)
                else:
                    return flask.jsonify({})
        else:
            return flask.jsonify({})