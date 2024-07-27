from .tool.func import *

async def api_list_markup():
    other_set = {}

    return flask.Response(response = (await python_to_golang(sys._getframe().f_code.co_name, other_set)), status = 200, mimetype = 'application/json')