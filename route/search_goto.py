from .tool.func import *

def search_goto_2(conn):
    curs = conn.cursor()

    curs.execute("select title from data where title = ?", [flask.request.form.get('search', 'test')])
    data = curs.fetchall()
    if data:
        return redirect('/w/' + url_pas(flask.request.form.get('search', 'test')))
    else:
        return redirect('/search/' + url_pas(flask.request.form.get('search', 'test')))