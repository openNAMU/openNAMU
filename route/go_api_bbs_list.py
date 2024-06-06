from .tool.func import *

def api_bbs_list():
    return flask.Response(response = python_to_golang(sys._getframe().f_code.co_name), status = 200, mimetype = 'application/json')