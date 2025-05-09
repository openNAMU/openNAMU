from .tool.func import *

async def api_func_search(name = 'Test', search_type = 'title', num = 1):
    other_set = {}
    other_set["name"] = name
    other_set["search_type"] = search_type
    other_set["num"] = str(num)

    return await python_to_golang(sys._getframe().f_code.co_name, other_set)

async def api_func_search_exter(name = 'Test', search_type = 'title', num = 1):
    return flask.jsonify(await api_func_search(name, search_type, num))