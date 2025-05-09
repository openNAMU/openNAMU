from .tool.func import *

async def recent_history_hidden(name = 'Test', rev = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        num = str(rev)

        if await acl_check(tool = 'hidel_auth', memo = 'history_hidden (' + name + '#' + num + ')') != 1:
            curs.execute(db_change("select title from history where title = ? and id = ? and hide = 'O'"), [name, num])
            if curs.fetchall():
                curs.execute(db_change("update history set hide = '' where title = ? and id = ?"), [name, num])
            else:
                curs.execute(db_change("update history set hide = 'O' where title = ? and id = ?"), [name, num])

        return redirect(conn, '/history/' + url_pas(name))