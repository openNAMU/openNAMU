from .tool.func import *

async def api_w_set_reset(name = 'Test'):
    other_set = {}
    other_set["name"] = name

    return flask.jsonify(await python_to_golang(sys._getframe().f_code.co_name, other_set))