from .tool.func import *

def del_alarm_2(conn):
    curs = conn.cursor()
    
    curs.execute("delete from alarm where name = ?", [ip_check()])
    conn.commit()

    return redirect('/alarm')