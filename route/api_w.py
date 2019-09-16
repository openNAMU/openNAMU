from .tool.func import *

def api_w_2(name):
    

    if acl_check(name, 'render') != 1:
        if flask.request.method == 'POST':
            json_data = { "title" : name, "data" : render_set(title = name, data = flask.request.form.get('data', '')) }
            
            return flask.jsonify(json_data)
        else:
            sqlQuery("select data from data where title = ?", [name])
            data = sqlQuery("fetchall")
            if data:
                json_data = { "title" : name, "data" : render_set(title = name, data = data[0][0]) }
            
                return flask.jsonify(json_data)
            else:
                return flask.jsonify({})
    else:
        return flask.jsonify({})