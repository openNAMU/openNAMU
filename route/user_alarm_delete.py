from .tool.func import *

def user_alarm_delete():
    with get_db_connect() as conn:
        curs = conn.cursor()
    
        curs.execute(db_change("delete from alarm where name = ?"), [ip_check()])
        conn.commit()

        return redirect('/alarm')