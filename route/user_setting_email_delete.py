from .tool.func import *

def user_setting_email_delete():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if ip_or_user(ip) == 0:
            curs.execute(db_change("delete from user_set where name = 'email' and id = ?"), [ip])
            conn.commit()
    
        return redirect('/change')