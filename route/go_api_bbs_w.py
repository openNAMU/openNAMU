from .tool.func import *

async def api_bbs_w(sub_code = '', legacy = 'on'):
    other_set = {}
    other_set["legacy"] = legacy
    other_set['sub_code'] = sub_code

    return await python_to_golang(sys._getframe().f_code.co_name, other_set)

async def api_bbs_w_exter(sub_code = '', legacy = 'on'):
    return flask.jsonify(await api_bbs_w(sub_code, legacy))