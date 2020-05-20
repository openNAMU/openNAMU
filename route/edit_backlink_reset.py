from .tool.func import *

def edit_backlink_reset_2(conn, name):
    curs = conn.cursor()

    curs.execute(db_change("select data from data where title = ?"), [name])
    old = curs.fetchall()
    if old:
        curs.execute(db_change("delete from back where link = ?"), [name])
        curs.execute(db_change("delete from back where title = ? and type = 'no'"), [name])
        
        render_set(
            title = name,
            data = old[0][0],
            num = 1
        )

    return redirect('/xref/' + url_pas(name))