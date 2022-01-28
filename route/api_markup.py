from .tool.func import *

def api_markup():
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change('select data from other where name = "markup"'))
        rep_data = curs.fetchall()
        if rep_data and rep_data[0][0] != '':
            return flask.jsonify({ "markup" : rep_data[0][0] })
        else:
            return flask.jsonify({})