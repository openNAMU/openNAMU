from .tool.func import *

async def api_w_watch_list(name = 'Test', do_type = 'watch_list', num = 1):
    other_set = {}
    other_set["name"] = name
    other_set["do_type"] = do_type
    other_set["ip"] = ip_check()
    other_set["num"] = str(num)

    return flask.Response(response = (await python_to_golang(sys._getframe().f_code.co_name, other_set)), status = 200, mimetype = 'application/json')