from .tool.func import *

async def api_func_sha224(data = 'Test'):
    other_set = {}
    other_set["data"] = data

    return flask.jsonify(await python_to_golang(sys._getframe().f_code.co_name, other_set))