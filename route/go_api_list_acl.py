from .tool.func import *

async def api_list_acl(data_type = ''):
    other_set = {}
    other_set = data_type

    return flask.jsonify(await python_to_golang(sys._getframe().f_code.co_name, other_set))