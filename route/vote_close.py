from .tool.func import *

async def vote_close(num = 1):
    num = str(num)
    
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'vote') == 1:
            return await re_error(conn, 0)

        curs.execute(db_change('select type from vote where id = ? and user = ""'), [num])
        data_list = curs.fetchall()
        if not data_list:
            return redirect(conn, '/vote')

        curs.execute(db_change('select data from vote where id = ? and name = "open_user" and type = "option"'), [num])
        db_data = curs.fetchall()
        open_user = db_data[0][0] if db_data else ''
        if open_user != ip_check() and await acl_check('', 'vote_auth', '', '') == 1:
            return await re_error(conn, 0)

        if data_list[0][0] == 'close':
            type_set = 'open'
        elif data_list[0][0] == 'n_close':
            type_set = 'n_open'
        elif data_list[0][0] == 'open':
            type_set = 'close'
        else:
            type_set = 'n_close'

        curs.execute(db_change("update vote set type = ? where user = '' and id = ? and type = ?"), [type_set, num, data_list[0][0]])
        curs.execute(db_change('delete from vote where name = "end_date" and type = "option" and id = ?'), [num])

        if data_list[0][0] == 'close' or data_list[0][0] == 'n_close':
            return redirect(conn, '/vote')
        else:
            return redirect(conn, '/vote/list/close')