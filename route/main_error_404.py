from .tool.func import *

def main_error_404_2(conn):
    curs = conn.cursor()

    if os.path.exists('404.html') and flask.request.path != '/':
        return open('404.html', 'r').read()
    else:
        return redirect('/w/' + url_pas(wiki_set(2)))