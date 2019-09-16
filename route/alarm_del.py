from .tool.func import *

def alarm_del_2():
    
    
    sqlQuery("delete from alarm where name = ?", [ip_check()])
    sqlQuery("commit")

    return redirect('/alarm')