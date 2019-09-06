from .tool.func import *

def main_views_2(conn, name):
    curs = conn.cursor()

    m = re.search('\.([^.]+)$', name)
    if m:
        g = m.groups()
    else:
        g = ['']

    if g[0] == 'css':
        c_open = open('./views/' + name, encoding='utf-8')
        f_open = c_open.read()
        c_open.close()
        return flask.Response(easy_minify(f_open, 'css'), mimetype="text/css")
    elif g[0] == 'js':
        c_open = open('./views/' + name, encoding='utf-8')
        f_open = c_open.read()
        c_open.close()
        return flask.Response(easy_minify(f_open, 'js'), mimetype="text/js")
    elif g[0] == 'html':
        c_open = open('./views/' + name, encoding='utf-8')
        f_open = c_open.read()
        c_open.close()
        return flask.Response(easy_minify(f_open), mimetype="text/html")
    else:
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