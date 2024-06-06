from .tool.func import *

def api_w_set_reset(name = 'Test'):
    other_set = {}
    other_set["name"] = name
    other_set["ip"] = ip_check()

    return flask.Response(response = python_to_golang(sys._getframe().f_code.co_name, other_set), status = 200, mimetype = 'application/json')