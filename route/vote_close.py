from .tool.func import *

def vote_close(num = 1):
    num = str(num)
    
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/ban')

        curs.execute(db_change('select type from vote where id = ? and user = ""'), [num])
        data_list = curs.fetchall()
        if not data_list:
            return redirect('/vote')

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
        conn.commit()

        if data_list[0][0] == 'close' or data_list[0][0] == 'n_close':
            return redirect('/vote')
        else:
            return redirect('/vote/list/close')