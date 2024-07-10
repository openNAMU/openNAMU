from .tool.func import *

async def api_give_auth():
    if flask.request.method == 'PATCH':
        func_name = sys._getframe().f_code.co_name
        func_name += '_patch'

        other_set = {}
        other_set["ip"] = ip_check()
        other_set["user_name"] = flask.request.form.get('user_name', '')
        other_set['auth'] = flask.request.form.get('auth', '')
        other_set['change_auth'] = flask.request.form.get('change_auth', '')

        return flask.Response(response = (await python_to_golang(func_name, other_set)), status = 200, mimetype = 'application/json')
    else:
        return flask.jsonify({}) 