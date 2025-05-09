from .tool.func import *

async def api_bbs_w_comment_one(sub_code = '', legacy = 'on', tool = ''):
    other_set = {}
    other_set["sub_code"] = sub_code
    other_set["legacy"] = legacy
    other_set["tool"] = tool

    return await python_to_golang(sys._getframe().f_code.co_name, other_set)

async def api_bbs_w_comment_one_exter(sub_code = '', legacy = 'on', tool = ''):
    return flask.jsonify(await api_bbs_w_comment_one(sub_code, legacy, tool))