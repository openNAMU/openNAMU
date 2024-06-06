from .tool.func import *

def api_func_ip_menu(ip = "Test", option = ""):
    other_set = {}
    other_set["ip"] = ip
    other_set["my_ip"] = ip_check()
    other_set["option"] = option

    return flask.Response(response = python_to_golang(sys._getframe().f_code.co_name, other_set), status = 200, mimetype = 'application/json')