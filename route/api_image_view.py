from .tool.func import *

async def api_image_view(name = 'Test'):
    with get_db_connect() as conn:
        if os.path.exists(os.path.join(load_image_url(conn), name)):
            return flask.jsonify({ "exist" : "1" })
        else:
            return flask.jsonify({})