from .tool.func import *

async def api_list_markup():
    return flask.jsonify(await python_to_golang(sys._getframe().f_code.co_name))