from .tool.func import *
import pymysql

def watch_list_name_2(conn, name):
    curs = conn.cursor()
    
    if custom()[2] == 0:
        return redirect('/login')

    ip = ip_check()

    curs.execute("select count(title) from scan where user = %s", [ip])
    count = curs.fetchall()
    if count and count[0][0] > 9:
        return redirect('/watch_list')

    curs.execute("select title from scan where user = %s and title = %s", [ip, name])
    if curs.fetchall():
        curs.execute("delete from scan where user = %s and title = %s", [ip, name])
    else:
        curs.execute("insert into scan (user, title) values (%s, %s)", [ip, name])
    
    conn.commit()

    return redirect('/watch_list')