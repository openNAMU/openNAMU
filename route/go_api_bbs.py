from .tool.func import *

async def api_bbs(bbs_num = "", page = 1):
    other_set = {}
    other_set["bbs_num"] = str(bbs_num)
    other_set["page"] = str(page)

    return await python_to_golang(sys._getframe().f_code.co_name, other_set)

async def api_bbs_exter(bbs_num = "", page = 1):
    return flask.jsonify(await api_bbs(bbs_num, page))