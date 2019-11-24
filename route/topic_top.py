from .tool.func import *

def topic_top_2(conn, name, sub, num):
    curs = conn.cursor()
    
    if admin_check(3, 'notice (' + name + ' - ' + sub + '#' + str(num) + ')') != 1:
        return re_error('/error/3')

    curs.execute(db_change("select title from topic where title = ? and sub = ? and id = ?"), [name, sub, str(num)])
    if curs.fetchall():
        curs.execute(db_change("select top from topic where id = ? and title = ? and sub = ?"), [str(num), name, sub])
        top_data = curs.fetchall()
        if top_data:
            if top_data[0][0] == 'O':
                curs.execute(db_change("update topic set top = '' where title = ? and sub = ? and id = ?"), [name, sub, str(num)])
            else:
                curs.execute(db_change("update topic set top = 'O' where title = ? and sub = ? and id = ?"), [name, sub, str(num)])
        
        rd_plus(name, sub, get_time())

        conn.commit()

    return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num))        
