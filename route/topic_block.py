from .tool.func import *
import pymysql

def topic_block_2(conn, name, sub, num):
    curs = conn.cursor()

    if admin_check(3, 'blind (' + name + ' - ' + sub + '#' + str(num) + ')') != 1:
        return re_error('/error/3')

    curs.execute("select block from topic where title = %s and sub = %s and id = %s", [name, sub, str(num)])
    block = curs.fetchall()
    if block:
        if block[0][0] == 'O':
            curs.execute("update topic set block = '' where title = %s and sub = %s and id = %s", [name, sub, str(num)])
        else:
            curs.execute("update topic set block = 'O' where title = %s and sub = %s and id = %s", [name, sub, str(num)])
        
        rd_plus(name, sub, get_time())
        
        conn.commit()
        
    return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num))