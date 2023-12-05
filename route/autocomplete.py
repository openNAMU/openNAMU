from .tool.func import *

def autocomplete(name='Test', limit=5):
    with get_db_connect() as conn:
        curs = conn.cursor()

        limit = min(limit, 10)

        curs.execute(db_change("select title from data where title like ? order by title limit ?"),
                     ['%' + name + '%', limit])

        all_titles = curs.fetchall()

        if all_titles:
            return flask.jsonify(all_titles)
        else:
            return flask.jsonify({})