from .tool.func import *

def api_image_view(name = 'Test'):
    with get_db_connect() as conn:
        if flask.request.method == 'POST':
            try:
                title_list = json.loads(flask.request.form.get('title_list', ''))
                title_list = list(set(title_list))
            except:
                title_list = []

            data_list = {}
            for i in title_list:
                if os.path.exists(os.path.join(load_image_url(), i)):
                    data_list[i] = '1'

            return flask.jsonify(data_list)
        else:
            if os.path.exists(os.path.join(load_image_url(), name)):
                return flask.jsonify({ "exist" : "1" })
            else:
                return flask.jsonify({})