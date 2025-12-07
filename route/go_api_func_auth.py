from .tool.func import *

async def api_func_auth(user_name = ''):
    other_set = {}
    other_set["ip"] = ip_check() if user_name == '' else user_name

    return await python_to_golang(sys._getframe().f_code.co_name, other_set)

async def api_func_auth_exter(user_name = ''):
    return flask.jsonify(await api_func_auth(user_name))