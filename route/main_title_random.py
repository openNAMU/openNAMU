from .tool.func import *

def main_title_random_2(conn):
    curs = conn.cursor()

    curs.execute(db_change("" + \
        "select title from data " + \
        "where title not like 'user:%' and title not like 'category:%' and title not like 'file:%' " + \
        "order by random() limit 1" + \
    ""))
    data = curs.fetchall()
    return redirect('/w/' + url_pas(data[0][0])) if data else redirect()