from .tool.func import *

def api_func_auth_list(user_name = ''):
    other_set = {}
    other_set["ip"] = ip_check() if user_name == '' else user_name

    return flask.Response(response = python_to_golang(sys._getframe().f_code.co_name, other_set), status = 200, mimetype = 'application/json')