from .tool.func import *

def api_func_language(data = 'Test'):
    other_set = {}
    if flask.request.method == 'POST':
        other_set["data"] = flask.request.form.get('data', '')
        other_set["data"] = other_set["data"].split(' ')
    else:
        other_set["data"] = [data]

    return flask.Response(response = python_to_golang(sys._getframe().f_code.co_name, other_set), status = 200, mimetype = 'application/json')