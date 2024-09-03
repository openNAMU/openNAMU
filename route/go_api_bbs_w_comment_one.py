from .tool.func import *

async def api_bbs_w_comment_one(sub_code = '', legacy = 'on', tool = ''):
    other_set = {}
    other_set["sub_code"] = sub_code
    other_set["legacy"] = legacy
    other_set["tool"] = tool
    other_set["ip"] = ip_check()

    return flask.Response(response = (await python_to_golang(sys._getframe().f_code.co_name, other_set)), status = 200, mimetype = 'application/json')