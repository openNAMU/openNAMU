from .tool.func import *

def main_search_goto(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if flask.request.form.get('search', None):
            data = flask.request.form.get('search', 'Test')
        else:
            data = name

        curs.execute(db_change("select title from data where title = ? collate nocase"), [data])
        db_data = curs.fetchall()
        if db_data:
            return redirect('/w/' + url_pas(db_data[0][0]))
        else:
            return redirect('/search/' + url_pas(data))