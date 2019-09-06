from .tool.func import *
import pymysql

def topic_top_2(conn, name, sub, num):
    curs = conn.cursor()
    
    if admin_check(3, 'notice (' + name + ' - ' + sub + '#' + str(num) + ')') != 1:
        return re_error('/error/3')

    curs.execute("select title from topic where title = %s and sub = %s and id = %s", [name, sub, num])
    if curs.fetchall():
        curs.execute("select top from topic where id = %s and title = %s and sub = %s", [num, name, sub])
        top_data = curs.fetchall()
        if top_data:
            if top_data[0][0] == 'O':
                curs.execute("update topic set top = '' where title = %s and sub = %s and id = %s", [name, sub, num])
            else:
                curs.execute("update topic set top = 'O' where title = %s and sub = %s and id = %s", [name, sub, num])
        
        rd_plus(name, sub, get_time())

        

    return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num))        
