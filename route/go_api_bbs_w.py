from .tool.func import *

def api_bbs_w(sub_code = '', legacy = 'on'):
    other_set = {}
    other_set['ip'] = ip_check()
    other_set["legacy"] = legacy
    other_set['sub_code'] = sub_code

    return flask.Response(response = python_to_golang_sync(sys._getframe().f_code.co_name, other_set), status = 200, mimetype = 'application/json')