from .tool.func import *

async def api_bbs_set():
    other_set = {}

    func_name = sys._getframe().f_code.co_name
    func_name += '_put'

    return flask.Response(response = (await python_to_golang(func_name, other_set)), status = 200, mimetype = 'application/json')