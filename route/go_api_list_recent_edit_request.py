from .tool.func import *

async def api_list_recent_edit_request(num = 1, set_type = 'normal', limit = 50):
    other_set = {}
    other_set["num"] = str(num)
    other_set["limit"] = str(limit)
    other_set["set_type"] = set_type

    return await python_to_golang(sys._getframe().f_code.co_name, other_set)

async def api_list_recent_edit_request_exter(num = 1, set_type = 'normal', limit = 50):
    return flask.jsonify(await api_list_recent_edit_request_exter(num, set_type, limit))