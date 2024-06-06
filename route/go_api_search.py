from .tool.func import *

def api_search(name = 'Test', search_type = 'title', num = 1):
    other_set = {}
    other_set["name"] = name
    other_set["search_type"] = search_type
    other_set["num"] = str(num)

    return flask.Response(response = python_to_golang(sys._getframe().f_code.co_name, other_set), status = 200, mimetype = 'application/json')