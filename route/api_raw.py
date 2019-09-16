from .tool.func import *

def api_raw_2(name):
    

    if acl_check(name, 'render') != 1:
        sqlQuery("select data from data where title = ?", [name])
        data = sqlQuery("fetchall")
        if data:
            json_data = { "title" : name, "data" : render_set(title = name, data = data[0][0], s_data = 1) }
        
            return flask.jsonify(json_data)
        else:
            return flask.jsonify({})
    else:
        return flask.jsonify({})