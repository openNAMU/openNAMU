from .tool.func import *

async def view_w_watch_list(name, num = 1, do_type = 'watch_list'):
    other_set = {}
    other_set['name'] = name
    other_set['num'] = str(num)
    other_set['do_type'] = do_type

    data = await python_to_golang(sys._getframe().f_code.co_name, other_set)

    return data["data"]