from .tool.func import *

async def api_func_language(legacy = 'on', data = 'Test'):
    other_set = {}
    if flask.request.method == 'POST':
        other_set["data"] = flask.request.form.get('data', '')
        other_set["data"] = other_set["data"].split(' ')
    else:
        other_set["data"] = [data]

    other_set["legacy"] = legacy

    return flask.Response(response = (await python_to_golang(sys._getframe().f_code.co_name, other_set)), status = 200, mimetype = 'application/json')