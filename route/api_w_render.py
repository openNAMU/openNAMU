from .tool.func import *

def api_w_render(name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if flask.request.method == 'POST':
            data_org = flask.request.form.get('data', '')
            data_pas = render_set(
                doc_name = name, 
                doc_data = data_org, 
                data_type = 'api_view'
            )

            return flask.jsonify({
                "data" : data_pas[0], 
                "js_data" : data_pas[1]
            })
        else:
            return ''