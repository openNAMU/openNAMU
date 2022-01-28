from .tool.func import *

def api_sha224(data = 'Test'):
    if flask.request.method == 'POST':
        try:
            title_list = json.loads(flask.request.form.get('title_list', ''))
            title_list = list(set(title_list))
        except:
            title_list = []

        data_list = {}
        for i in title_list:
            data_list[i] = sha224_replace(i)

        return flask.jsonify(data_list)
    else:
        return flask.jsonify({ "data" : sha224_replace(data) })