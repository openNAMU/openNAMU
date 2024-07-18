from .tool.func import *

async def api_bbs_w_comment(sub_code = '', tool = "", legacy = 'on'):
    other_set = {}
    other_set["sub_code"] = sub_code
    other_set["tool"] = tool
    other_set["legacy"] = legacy

    return flask.Response(response = (await python_to_golang(sys._getframe().f_code.co_name, other_set)), status = 200, mimetype = 'application/json')