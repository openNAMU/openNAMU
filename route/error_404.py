from .tool.func import *

def error_404_2(conn):
    curs = conn.cursor()

    return redirect('/w/' + url_pas(wiki_set(2)))