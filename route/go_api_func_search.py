from .tool.func import *

async def api_func_search(name = 'Test', search_type = 'title', num = 1):
    other_set = {}
    other_set["name"] = name
    other_set["search_type"] = search_type
    other_set["num"] = str(num)

    return flask.Response(response = (await python_to_golang(sys._getframe().f_code.co_name, other_set)), status = 200, mimetype = 'application/json')