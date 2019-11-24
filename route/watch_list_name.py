from .tool.func import *

def watch_list_name_2(conn, name):
    curs = conn.cursor()
    
    ip = ip_check()
    if ip_or_user(ip) != 0:
        return redirect('/login')

    curs.execute(db_change("select count(title) from scan where user = ?"), [ip])
    count = curs.fetchall()
    if count and count[0][0] > 9:
        return redirect('/watch_list')

    curs.execute(db_change("select title from scan where user = ? and title = ?"), [ip, name])
    if curs.fetchall():
        curs.execute(db_change("delete from scan where user = ? and title = ?"), [ip, name])
    else:
        curs.execute(db_change("insert into scan (user, title) values (?, ?)"), [ip, name])
    
    conn.commit()

    return redirect('/watch_list')