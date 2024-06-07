from .tool.func import *

async def api_w_raw(name = 'Test', rev = '', exist_check = ''):
    with get_db_connect() as conn:
        if acl_check(conn, name, 'render') != 1:
            other_set = {}
            other_set["name"] = name
            other_set["rev"] = str(rev)
            other_set["exist_check"] = exist_check

            return flask.Response(response = (await python_to_golang(sys._getframe().f_code.co_name, other_set)), status = 200, mimetype = 'application/json')
        else:
            return flask.jsonify({})