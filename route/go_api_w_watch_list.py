from .tool.func import *

async def api_w_watch_list(name = 'Test', do_type = 'watch_list', num = 1):
    other_set = {}
    other_set["name"] = name
    other_set["do_type"] = do_type
    other_set["num"] = str(num)

    return flask.jsonify(await python_to_golang(sys._getframe().f_code.co_name, other_set))