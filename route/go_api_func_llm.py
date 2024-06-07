from .tool.func import *

async def api_func_llm():
    if flask.request.method == 'POST':
        other_set = {}
        other_set["prompt"] = flask.request.form.get('prompt', '')
        other_set["ip"] = ip_check()

        return flask.Response(response = (await python_to_golang(sys._getframe().f_code.co_name, other_set)), status = 200, mimetype = 'application/json')
    else:
        return flask.jsonify({})