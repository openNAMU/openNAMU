from .tool.func import *

async def api_bbs_w_tabom(sub_code = ''):
    other_set = {}
    other_set["sub_code"] = sub_code

    func_name = sys._getframe().f_code.co_name
    if flask.request.method == 'POST':
        func_name += '_post'

    return flask.jsonify(await python_to_golang(func_name, other_set))