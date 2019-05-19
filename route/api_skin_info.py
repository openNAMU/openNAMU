from .tool.func import *

def api_skin_info_2(conn):
    curs = conn.cursor()

    json_address = re.sub("(((?!\.|\/).)+)\.html$", "info.json", skin_check())
    try:
        json_data = json.loads(open(json_address).read())
    except:
        json_data = None

    if json_data:    
        return flask.jsonify(json_data)
    else:
        return flask.jsonify({}), 404