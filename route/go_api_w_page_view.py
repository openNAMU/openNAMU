from .tool.func import *

async def api_w_page_view(name = 'Test'):
    other_set = {}
    other_set["doc_name"] = name

    return flask.jsonify(await python_to_golang(sys._getframe().f_code.co_name, other_set))