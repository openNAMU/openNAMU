from .tool.func import *

async def api_list_old_page(num = 1, set_type = 'old'):
    other_set = {}
    other_set["num"] = str(num)
    other_set["set_type"] = set_type

    return await python_to_golang(sys._getframe().f_code.co_name, other_set)

async def api_list_old_page_exter(num = 1, set_type = 'old'):
    return flask.jsonify(await api_list_old_page(num, set_type))