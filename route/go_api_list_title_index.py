from .tool.func import *

def api_list_title_index(num = 1):
    other_set = {}
    other_set["num"] = str(num)

    return flask.Response(response = python_to_golang(sys._getframe().f_code.co_name, other_set), status = 200, mimetype = 'application/json')