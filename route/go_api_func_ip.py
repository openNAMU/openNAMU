from .tool.func import *

async def api_func_ip(data = 'Test'):
    other_set = {}

    func_name = sys._getframe().f_code.co_name
    if flask.request.method == 'POST':
        func_name += '_post'

        for_a = 1
        while 1:
            data = flask.request.form.get('data_' + str(for_a), '')
            
            if data == '':
                break
            else:
                other_set['data_' + str(for_a)] = data

            for_a += 1
    else:
        other_set["data"] = data

    return flask.jsonify(await python_to_golang(func_name, other_set))