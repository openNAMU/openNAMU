from .tool.func import *
import pymysql

def alarm_del_2(conn):
    curs = conn.cursor()
    
    curs.execute("delete from alarm where name = %s", [ip_check()])
    

    return redirect('/alarm')