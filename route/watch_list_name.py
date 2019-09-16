from .tool.func import *

def watch_list_name_2(name):
    
    
    if custom()[2] == 0:
        return redirect('/login')

    ip = ip_check()

    sqlQuery("select count(title) from scan where user = ?", [ip])
    count = sqlQuery("fetchall")
    if count and count[0][0] > 9:
        return redirect('/watch_list')

    sqlQuery("select title from scan where user = ? and title = ?", [ip, name])
    if sqlQuery("fetchall"):
        sqlQuery("delete from scan where user = ? and title = ?", [ip, name])
    else:
        sqlQuery("insert into scan (user, title) values (?, ?)", [ip, name])
    
    sqlQuery("commit")

    return redirect('/watch_list')