from .tool.func import *

def api_w_raw(name = 'Test', rev = '', exist_check = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if exist_check != '':
            curs.execute(db_change("select title from data where title = ?"), [name])
            if data:
                return flask.jsonify({ 'exist' : '1' })
            else:
                return flask.jsonify({})
        else:
            if rev != '':
                curs.execute(db_change("select data from history where title = ? and id = ?"), [name, rev])
            else:
                curs.execute(db_change("select data from data where title = ?"), [name])

            data = curs.fetchall()
            if data:
                return flask.jsonify({
                    "title" : name, 
                    "data" : render_set(
                        doc_name = name, 
                        doc_data = data[0][0],
                        data_type = 'raw'
                    )
                })
            else:
                return flask.jsonify({})