from .tool.func import *

APPVAR = json.loads(open('data/app_variables.json', encoding='utf-8').read())

def image_view_2(conn, name):
    curs = conn.cursor()
    
    if os.path.exists(os.path.join(APPVAR['PATH_DATA_IMAGES'], name)):
        return flask.send_from_directory('.'+APPVAR['PATH_DATA_IMAGES'], name)
    else:
        return redirect()