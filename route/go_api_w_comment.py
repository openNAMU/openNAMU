from .tool.func import *

async def api_w_comment(page = 1, name = 'Test'):
    other_set = {}
    other_set["num"] = str(page)
    other_set["doc_name"] = name
    other_set["ip"] = ip_check()

    func_name = sys._getframe().f_code.co_name
    if flask.request.method == 'POST':
        func_name += '_post'
    elif flask.request.method == 'DELETE':
        func_name += '_delete'

    return flask.Response(response = (await python_to_golang(func_name, other_set)), status = 200, mimetype = 'application/json')