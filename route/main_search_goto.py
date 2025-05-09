from .tool.func import *

from .go_api_func_search import api_func_search

async def main_search_goto(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if flask.request.form.get('search', None):
            data = flask.request.form.get('search', 'Test')
        else:
            data = name

        search_data = await api_func_search(data, 'title', 1)

        curs.execute(db_change("select title from data where title = ? collate nocase"), [data])
        db_data = curs.fetchall()
        if db_data:
            return redirect(conn, '/w/' + url_pas(db_data[0][0]))
        elif len(search_data) == 1:
            return redirect(conn, '/w/' + url_pas(search_data[0]))
        else:
            return redirect(conn, '/search/' + url_pas(data))