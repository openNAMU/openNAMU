from .tool.func import *

async def api_w_page_view(name = 'Test'):
    other_set = {}
    other_set["doc_name"] = name

    return flask.Response(response = (await python_to_golang(sys._getframe().f_code.co_name, other_set)), status = 200, mimetype = 'application/json')