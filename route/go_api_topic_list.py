from .tool.func import *

async def api_topic_list(name = 'Test', set_type = 'normal', num = 1):
    other_set = {}
    other_set["name"] = str(name)
    other_set["set_type"] = set_type
    other_set["num"] = str(num)

    return flask.jsonify(await python_to_golang(sys._getframe().f_code.co_name, other_set))