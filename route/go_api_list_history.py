from .tool.func import *

async def api_list_history(num = 1, set_type = 'normal', doc_name = 'Test'):
    other_set = {}
    other_set["num"] = str(num)
    other_set["doc_name"] = doc_name
    other_set["set_type"] = set_type
    
    return await python_to_golang(sys._getframe().f_code.co_name, other_set)

async def api_list_history_exter(num = 1, set_type = 'normal', doc_name = 'Test'):
    return flask.jsonify(await api_list_history(num, set_type, doc_name))