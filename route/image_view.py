from .tool.func import *

def image_view_2(conn, name):
    curs = conn.cursor()
    
    if os.path.exists(os.path.join('data/images', name)):
        return flask.send_from_directory('./data/images', name)
    else:
        return redirect()