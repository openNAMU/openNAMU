from .tool.func import *

def api_image_view(name = 'Test'):
    with get_db_connect() as conn:
        if os.path.exists(os.path.join(load_image_url(), name)):
            return flask.jsonify({ "exist" : "1" })
        else:
            return flask.jsonify({})