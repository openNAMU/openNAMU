from .tool.func import *

async def api_user_rankup():
    other_set = {}
    
    func_name = sys._getframe().f_code.co_name
    if flask.request.method == 'PATCH':
        func_name += '_patch'

    return flask.jsonify(await python_to_golang(func_name, other_set))