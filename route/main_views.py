from .tool.func import *

def main_views_2(conn, name):
    curs = conn.cursor()
    
    if re.search('\/', name):
        m = re.search('^(.*)\/(.*)$', name)
        if m:
            n = m.groups()
            plus = '/' + n[0]
            rename = n[1]
        else:
            plus = ''
            rename = name
    else:
        plus = ''
        rename = name

    return flask.send_from_directory('./views' + plus, rename)