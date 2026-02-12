from .tool.func import *

async def api_func_language(legacy = 'on', data = 'Test', safe = ''):
    other_set = {}
    if flask.request.method == 'POST':
        other_set["data"] = flask.request.form.get('data', '')
    else:
        other_set["data"] = data

    other_set["legacy"] = legacy
    other_set["safe"] = safe

    return await python_to_golang(sys._getframe().f_code.co_name, other_set)

async def api_func_language_exter(legacy = 'on', data = 'Test'):
    return flask.jsonify(await api_func_language(legacy, data))