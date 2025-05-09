from .tool.func import *

async def api_give_auth():
    if flask.request.method == 'PATCH':
        func_name = sys._getframe().f_code.co_name
        func_name += '_patch'

        other_set = {}
        other_set["user_name"] = flask.request.form.get('user_name', '')
        other_set['auth'] = flask.request.form.get('auth', '')
        other_set['change_auth'] = flask.request.form.get('change_auth', '')

        return flask.jsonify(await python_to_golang(func_name, other_set))
    else:
        return flask.jsonify({}) 