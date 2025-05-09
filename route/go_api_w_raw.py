from .tool.func import *

async def api_w_raw(name = 'Test', rev = '', exist_check = ''):
    other_set = {}
    other_set["name"] = name
    other_set["rev"] = str(rev)
    other_set["exist_check"] = exist_check

    return await python_to_golang(sys._getframe().f_code.co_name, other_set)

async def api_w_raw_exter(name = 'Test', rev = '', exist_check = ''):
    return flask.jsonify(await api_w_raw(name, rev, exist_check))