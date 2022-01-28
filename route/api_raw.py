from .tool.func import *

def api_raw(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if acl_check(name, 'render') != 1:
            rev = flask.request.args.get('num', '')
            if rev != '':
                curs.execute(db_change("select data from history where title = ? and id = ?"), [name, rev])
            else:
                curs.execute(db_change("select data from data where title = ?"), [name])
            data = curs.fetchall()
            if data:
                json_data = {
                    "title" : name, 
                    "data" : render_set(
                        doc_name = name, 
                        doc_data = data[0][0],
                        data_type = 'raw'
                    )
                }

                return flask.jsonify(json_data)

        return flask.jsonify({})