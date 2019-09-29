from .tool.func import *

def api_w_2(name):
    

    if flask.request.args.get('exist', None):
        sqlQuery("select title from data where title = ?", [name])
        if sqlQuery("fetchall"):
            return flask.jsonify({ "exist" : "1" })
        else:
            return flask.jsonify({})
    else:
        if acl_check(name, 'render') != 1:
            if flask.request.method == 'POST':
                g_data = render_set(title = name, data = flask.request.form.get('data', ''), num = 2)
                
                return flask.jsonify({ "title" : name, "data" : g_data[0], "js_data" : g_data[1] })
            else:
                sqlQuery("select data from data where title = ?", [name])
                data = sqlQuery("fetchall")
                if data:
                    if flask.request.args.get('include', '1'):
                        include_re = re.compile('\[include\(((?:(?!\)\]).)+)\)\]', re.I)
                        json_data = include_re.sub('', data[0][0])
                    else:
                        json_data = g_data[0]

                    g_data = render_set(title = name, data = json_data, num = 2)

                    return flask.jsonify({ "title" : name, "data" : g_data[0], "js_data" : g_data[1] })
                else:
                    return flask.jsonify({})
        else:
            return flask.jsonify({})