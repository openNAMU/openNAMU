from .tool.func import *

def api_version_2(conn, version_list):
    curs = conn.cursor()

    curs.execute(db_change('select data from other where name = "update"'))
    up_data = curs.fetchall()
    up_data = up_data[0][0] if up_data and up_data[0][0] in ['stable', 'beta', 'dev'] else 'stable'

    json_data = {
        "version" : version_list['beta']['r_ver'], 
        "db_version" : version_list['beta']['c_ver'],
        "skin_version" : version_list['beta']['s_ver'],
        "build" : up_data
    }

    return flask.jsonify(json_data)