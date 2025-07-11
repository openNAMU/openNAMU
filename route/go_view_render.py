from .tool.func import *

async def view_render():
    other_set = {}

    data = await python_to_golang(sys._getframe().f_code.co_name, other_set)

    return data["data"]