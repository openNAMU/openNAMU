from .tool.func import *

def api_sha224_2(conn, data):
    curs = conn.cursor()

    return flask.jsonify({ "data" : sha224_replace(data) })