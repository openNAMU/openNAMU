from .tool.func import *

async def api_func_llm():
    if flask.request.method == 'POST':
        other_set = {}
        other_set["prompt"] = flask.request.form.get('prompt', '')

        return flask.jsonify(await python_to_golang(sys._getframe().f_code.co_name, other_set))
    else:
        return flask.jsonify({})