from .tool.func import *

async def api_bbs_list():
    return flask.jsonify(await python_to_golang(sys._getframe().f_code.co_name))