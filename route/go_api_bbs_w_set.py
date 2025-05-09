from .tool.func import *

async def api_bbs_w_set(name = 'Test', bbs_num = 1):
    other_set = {}
    other_set["set_name"] = name
    other_set["set_id"] = str(bbs_num)

    func_name = sys._getframe().f_code.co_name
    if flask.request.method == 'PUT':
        func_name += '_put'
        other_set['data'] = flask.request.form.get('data', 'Test')
    
    return flask.jsonify(await python_to_golang(func_name, other_set))