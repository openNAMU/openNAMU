from .tool.func import *

def api_skin_info_2(conn, name):
    curs = conn.cursor()

    if name == '':
        name = skin_check()
    else:
        name = './views/' + name + '/index.html'

    if not flask.request.args.get('all', None):
        json_address = re.sub("(((?!\.|\/).)+)\.html$", "info.json", name)
        try:
            json_data = json.loads(open(json_address).read())
        except:
            json_data = None

        if json_data:
            return flask.jsonify(json_data)
        else:
            return flask.jsonify({}), 404
    else:
        a_data = {}
        for i in load_skin(skin_check(1), 1):
            json_address = re.sub("(((?!\.|\/).)+)\.html$", "info.json", './views/' + i + '/index.html')
            try:
                json_data = json.loads(open(json_address).read())
            except:
                json_data = None

            if json_data:
                if i == skin_check(1):
                    json_data = {**json_data, **{ "main" : "true" }}

                a_data = {**a_data, **{ i : json_data }}

        if a_data == {}:
            return flask.jsonify({})
        else:
            return flask.jsonify(a_data)