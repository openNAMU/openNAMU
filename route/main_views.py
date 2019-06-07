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

    m = re.search('\.(.+)$', name)
    if m:
        g = m.groups()
    else:
        g = ['']

    if g == 'css':
        return easy_minify(flask.send_from_directory('./views' + plus, rename), 'css')   
    elif g == 'js':
        return easy_minify(flask.send_from_directory('./views' + plus, rename), 'js')
    elif g == 'html':
        return easy_minify(flask.send_from_directory('./views' + plus, rename))   
    else:
        return flask.send_from_directory('./views' + plus, rename)