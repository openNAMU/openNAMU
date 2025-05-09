from .tool.func import *

async def api_list_recent_change(num = 1, set_type = 'normal', limit = 10, legacy = 'on'):
    other_set = {}
    other_set["num"] = str(num)
    other_set["limit"] = str(limit)
    other_set["set_type"] = set_type
    other_set["legacy"] = legacy

    return await python_to_golang(sys._getframe().f_code.co_name, other_set)

async def api_list_recent_change_exter(num = 1, set_type = 'normal', limit = 10, legacy = 'on'):
    response = flask.make_response(flask.jsonify(await api_list_recent_change(num, set_type, limit, legacy)))
    
    response.headers.add("Access-Control-Allow-Origin", "*")
    response.headers.add('Access-Control-Allow-Headers', "Content-Type")
    response.headers.add('Access-Control-Allow-Methods', "GET")

    return response