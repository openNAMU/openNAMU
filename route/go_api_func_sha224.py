from .tool.func import *

def api_func_sha224(data = 'Test'):
    other_set = {}
    other_set["data"] = data

    return flask.Response(response = python_to_golang(sys._getframe().f_code.co_name, other_set), status = 200, mimetype = 'application/json')