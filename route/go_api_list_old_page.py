from .tool.func import *

async def api_list_old_page(num = 1, set_type = 'old'):
    other_set = {}
    other_set["num"] = str(num)
    other_set["set_type"] = set_type
    
    return flask.Response(response = (await python_to_golang(sys._getframe().f_code.co_name, other_set)), status = 200, mimetype = 'application/json')