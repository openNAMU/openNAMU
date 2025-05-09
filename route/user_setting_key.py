from .tool.func import *

async def user_setting_key():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if ip_or_user(ip) == 0:
            while 1:
                key = load_random_key()
                curs.execute(db_change('select data from user_set where name = "random_key" and data = ?'), [key])
                if not curs.fetchall():
                    break

            curs.execute(db_change("delete from user_set where name = 'random_key' and id = ?"), [ip])
            curs.execute(db_change("insert into user_set (name, id, data) values ('random_key', ?, ?)"), [ip, key])

        return redirect(conn, '/change')