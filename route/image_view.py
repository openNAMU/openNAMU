from .tool.func import *

def image_view_2(conn, name):
    curs = conn.cursor()
    
    if os.path.exists(os.path.join('image', name)):
        return flask.send_from_directory('./image', name)
    else:
        return redirect()