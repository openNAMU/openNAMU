from .tool.func import *
import pymysql

def give_history_hidden_2(conn, name):
    curs = conn.cursor()

    num = number_check(flask.request.args.get('num', '1'))

    if admin_check(6, 'history_hidden (' + name + '#' + num + ')') == 1:
        curs.execute("select title from history where title = %s and id = %s and hide = 'O'", [name, num])
        if curs.fetchall():
            curs.execute("update history set hide = '' where title = %s and id = %s", [name, num])
        else:
            curs.execute("update history set hide = 'O' where title = %s and id = %s", [name, num])
            
        
    
    return redirect('/history/' + url_pas(name))