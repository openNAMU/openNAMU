from .tool.func import *

def give_history_hidden_2(conn, name):
    curs = conn.cursor()

    num = number_check(flask.request.args.get('num', '1'))

    if admin_check(6, 'history_hidden (' + name + '#' + num + ')') == 1:
        curs.execute(db_change("select title from history where title = ? and id = ? and hide = 'O'"), [name, num])
        if curs.fetchall():
            curs.execute(db_change("update history set hide = '' where title = ? and id = ?"), [name, num])
        else:
            curs.execute(db_change("update history set hide = 'O' where title = ? and id = ?"), [name, num])

        conn.commit()

    return redirect('/history/' + url_pas(name))