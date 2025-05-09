from .tool.func import *

async def api_setting(name = 'Test', method = '', data = ''):
    other_set = {}
    other_set["set_name"] = name

    if method == '':
        method = flask.request.method

    func_name = sys._getframe().f_code.co_name
    if method == 'PUT':
        func_name += '_put'

        if data == '':
            other_set['data'] = flask.request.form.get('data', 'Test')
        else:
            other_set['data'] = data

    return await python_to_golang(func_name, other_set)

async def api_setting_exter(name = 'Test'):
    return flask.jsonify(await api_setting(name))