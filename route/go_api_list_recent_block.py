from .tool.func import *

async def api_list_recent_block(num = 1, set_type = 'all', user_name = 'Test', why = ''):
    other_set = {}
    other_set["num"] = str(num)
    other_set["set_type"] = set_type
    other_set["user_name"] = user_name
    other_set["why"] = why

    return flask.jsonify(await python_to_golang(sys._getframe().f_code.co_name, other_set))