from .tool.func import *

def api_func_lang(data = 'Test'):
    with get_db_connect() as conn:
        return flask.jsonify({ "data" : get_lang(conn, data) })