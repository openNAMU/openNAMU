from .tool.func import *

def api_skin_info_2(conn, name):
    curs = conn.cursor()

    if name == '':
        name = skin_check()
    else:
        name = './views/' + name + '/index.html'

    print(skin_check())
    print(name)

    json_address = re.sub("(((?!\.|\/).)+)\.html$", "info.json", name)
    try:
        json_data = json.loads(open(json_address).read())
    except:
        json_data = None

    if json_data:
        return flask.jsonify(json_data)
    else:
        return flask.jsonify({}), 404