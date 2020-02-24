from .tool.func import *

def api_image_view_2(conn, name, app_var):
    curs = conn.cursor()

    if os.path.exists(os.path.join(app_var['path_data_image'], name)):
        return flask.jsonify({ "exist" : "1" })
    else:
        return flask.jsonify({})