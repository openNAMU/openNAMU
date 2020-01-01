from .tool.func import *

def alarm_del_2(conn):
    curs = conn.cursor()

    curs.execute(db_change("delete from alarm where name = ?"), [ip_check()])
    conn.commit()

    return redirect('/alarm')