from .tool.func import *

async def api_w_xref(name = 'Test', page = 1, xref_type = '1'):
    other_set = {}
    other_set["name"] = name
    other_set["page"] = str(page)
    other_set["do_type"] = xref_type

    return flask.Response(response = (await python_to_golang(sys._getframe().f_code.co_name, other_set)), status = 200, mimetype = 'application/json')