from .tool.func import *

def title_random_2(conn):
    curs = conn.cursor()

    curs.execute("select title from data order by random() limit 1")
    data = curs.fetchall()
    if data:
        return redirect('/w/' + url_pas(data[0][0]))
    else:
        return redirect()