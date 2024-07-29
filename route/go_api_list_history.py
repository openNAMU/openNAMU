from .tool.func import *

async def api_list_history(num = 1, set_type = 'normal', doc_name = 'Test'):
    other_set = {}
    other_set["num"] = str(num)
    other_set["doc_name"] = doc_name
    other_set["set_type"] = set_type
    other_set["ip"] = ip_check()
    
    return flask.Response(response = (await python_to_golang(sys._getframe().f_code.co_name, other_set)), status = 200, mimetype = 'application/json')