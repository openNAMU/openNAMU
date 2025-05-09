from .tool.func import *

async def api_version(version_list):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change('select data from other where name = "update"'))
        up_data = curs.fetchall()
        up_data = up_data[0][0] if up_data and up_data[0][0] in ['stable', 'beta', 'dev'] else 'stable'

        json_data = {
            "version" : version_list['r_ver'], 
            "db_version" : version_list['c_ver'],
            "skin_version" : version_list['s_ver'],
            "build" : up_data
        }

        return flask.jsonify(json_data)