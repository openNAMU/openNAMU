from .tool.func import *

async def api_list_recent_block(num = 1, set_type = 'all', user_name = 'Test', why = ''):
    other_set = {}
    other_set["num"] = str(num)
    other_set["set_type"] = set_type
    other_set["user_name"] = user_name
    other_set["why"] = why
    other_set["ip"] = ip_check()

    return flask.Response(response = (await python_to_golang(sys._getframe().f_code.co_name, other_set)), status = 200, mimetype = 'application/json')