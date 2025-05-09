from .tool.func import *

async def user_watch_list_name(tool, name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if ip_or_user(ip) != 0:
            return redirect(conn, '/login')
        
        name_from = 0
        if tool == 'watch_list_from':
            name_from = 1
            tool = 'watch_list'
        elif tool == 'star_doc_from':
            name_from = 1
            tool = 'star_doc'

        if tool == 'watch_list':
            type_data = 'watchlist'
        else:
            type_data = 'star_doc'

        curs.execute(db_change("select data from user_set where name = ? and id = ? and data = ?"), [type_data, ip, name])
        if curs.fetchall():
            curs.execute(db_change("delete from user_set where name = ? and id = ? and data = ?"), [type_data, ip, name])
        else:
            if tool == 'watch_list':
                curs.execute(db_change("select count(*) from user_set where id = ? and name = ?"), [ip, type_data])
                count = curs.fetchall()
                if count and count[0][0] > 10:
                    return await re_error(conn, 28)

            curs.execute(db_change("insert into user_set (id, name, data) values (?, ?, ?)"), [ip, type_data, name])

        if name_from == 1:
            return redirect(conn, '/w/' + url_pas(name))
        else:
            if tool == 'watch_list':
                return redirect(conn, '/watch_list')
            else:
                return redirect(conn, '/star_doc')