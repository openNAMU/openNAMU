from .tool.func import *

async def api_user_setting_editor():
    other_set = {}
    
    func_name = sys._getframe().f_code.co_name
    if flask.request.method == 'POST':
        func_name += '_post'
        other_set['data'] = flask.request.form.get('data', 'Test')
    elif flask.request.method == 'DELETE':
        func_name += '_delete'
        other_set['data'] = flask.request.form.get('data', 'Test')

    return flask.jsonify(await python_to_golang(func_name, other_set))