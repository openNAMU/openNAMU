from .tool.func import *

def api_title_index():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if flask.request.args.get('count', '1') == '1':
            curs.execute(db_change('select data from other where name = "count_all_title"'))
            title_count = curs.fetchall()
            if title_count:
                return flask.jsonify({ 'count' : title_count[0][0] })
            else:
                return flask.jsonify({ 'count' : '0' })
        else:
            return flask.jsonify({})