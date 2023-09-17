from .tool.func import *

def user_alarm_delete(id = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()
    
        if id != '':
            curs.execute(db_change("delete from user_notice where name = ? and id = ?"), [ip_check(), str(id)])
        else:
            curs.execute(db_change("delete from user_notice where name = ?"), [ip_check()])
        
        conn.commit()

        return redirect('/alarm')