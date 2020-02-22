from .tool.func import *

def delete_admin_group_2(conn, name):
    if admin_check(None) != 1:
        return re_error('/error/3')
    
    curs.execute(db_change("delete from alist where name = ?"), [name])
    
    return redirect('/give_log')
