from .tool.func import *

async def api_bbs(bbs_num = "", page = 1):
    other_set = {}
    other_set["bbs_num"] = str(bbs_num)
    other_set["page"] = str(page)
    other_set["ip"] = ip_check()

    return flask.Response(response = (await python_to_golang(sys._getframe().f_code.co_name, other_set)), status = 200, mimetype = 'application/json')