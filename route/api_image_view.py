from .tool.func import *

def api_image_view_2(conn, name):
    curs = conn.cursor()

    if os.path.exists(os.path.join(load_image_url(), name)):
        return flask.jsonify({ "exist" : "1" })
    else:
        return flask.jsonify({ "exist" : "0" })