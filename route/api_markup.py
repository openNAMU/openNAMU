from .tool.func import *

def api_markup_2(conn, name):
    curs = conn.cursor()

    curs.execute('select data from other where name = "markup"')
    rep_data = curs.fetchall()
    if rep_data[0][0] != '':
        return flask.jsonify({ "markup" : rep_data[0][0] })
    else:
        return flask.jsonify({})