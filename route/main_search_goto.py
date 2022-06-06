from .tool.func import *

def main_search_goto(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if flask.request.form.get('search', None):
            data = flask.request.form.get('search', 'Test')
        else:
            data = name

        curs.execute(db_change("select title from data where title = ?"), [data])
        t_data = curs.fetchall()
        if t_data:
            return redirect('/w/' + url_pas(data))
        else:
            return redirect('/search/' + url_pas(data))