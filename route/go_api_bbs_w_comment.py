from .tool.func import *

async def api_bbs_w_comment_n(bbs_num = "", post_num = "", tool = "length"):
    other_set = {}
    other_set["bbs_num"] = str(bbs_num)
    other_set["post_num"] = str(post_num)
    other_set["tool"] = tool

    return flask.Response(response = (await python_to_golang(sys._getframe().f_code.co_name, other_set)), status = 200, mimetype = 'application/json')