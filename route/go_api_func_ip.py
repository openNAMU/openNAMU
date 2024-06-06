from .tool.func import *

def api_func_ip(data = 'Test'):
    other_set = {}
    other_set["data"] = data
    other_set["ip"] = ip_check()

    return flask.Response(response = python_to_golang(sys._getframe().f_code.co_name, other_set), status = 200, mimetype = 'application/json')