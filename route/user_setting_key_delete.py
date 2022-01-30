from .tool.func import *

def user_setting_key_delete():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if ip_or_user(ip) == 0:
            curs.execute(db_change("delete from user_set where name = 'random_key' and id = ?"), [ip])
            conn.commit()
    
        return redirect('/change')