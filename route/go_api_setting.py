from .tool.func import *

async def api_setting(name = 'Test'):
    other_set = {}
    other_set["set_name"] = name
    other_set["ip"] = ip_check()

    func_name = sys._getframe().f_code.co_name
    if flask.request.method == 'PUT':
        func_name += '_put'
        other_set['data'] = flask.request.form.get('data', 'Test')
    
    return flask.Response(response = (await python_to_golang(func_name, other_set)), status = 200, mimetype = 'application/json')