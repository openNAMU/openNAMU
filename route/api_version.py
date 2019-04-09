from .tool.func import *

def api_version_2(conn, r_ver, c_ver):
    curs = conn.cursor()

    n_ver = load_version()
        
    json_data = { "version" : r_ver, "db_version" : c_ver, "lastest_version" : n_ver  }

    return flask.jsonify(json_data)