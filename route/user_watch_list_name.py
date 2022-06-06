from .tool.func import *

def user_watch_list_name(tool, name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if ip_or_user(ip) != 0:
            return redirect('/login')

        if tool == 'watch_list':
            curs.execute(db_change("select count(*) from scan where user = ?"), [ip])
            count = curs.fetchall()
            if count and count[0][0] > 9:
                return re_error('/error/28')

            type_data = ''
        else:
            type_data = 'star'

        curs.execute(db_change("select title from scan where user = ? and title = ? and type = ?"), [ip, name, type_data])
        if curs.fetchall():
            curs.execute(db_change("delete from scan where user = ? and title = ? and type = ?"), [ip, name, type_data])
        else:
            curs.execute(db_change("insert into scan (user, title, type) values (?, ?, ?)"), [ip, name, type_data])

        conn.commit()

        if tool == 'watch_list':
            return redirect('/watch_list')
        else:
            return redirect('/star_doc')