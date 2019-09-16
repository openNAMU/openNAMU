from .tool.func import *

def topic_top_2(name, sub, num):
    
    
    if admin_check(3, 'notice (' + name + ' - ' + sub + '#' + str(num) + ')') != 1:
        return re_error('/error/3')

    sqlQuery("select title from topic where title = ? and sub = ? and id = ?", [name, sub, str(num)])
    if sqlQuery("fetchall"):
        sqlQuery("select top from topic where id = ? and title = ? and sub = ?", [str(num), name, sub])
        top_data = sqlQuery("fetchall")
        if top_data:
            if top_data[0][0] == 'O':
                sqlQuery("update topic set top = '' where title = ? and sub = ? and id = ?", [name, sub, str(num)])
            else:
                sqlQuery("update topic set top = 'O' where title = ? and sub = ? and id = ?", [name, sub, str(num)])
        
        rd_plus(name, sub, get_time())

        sqlQuery("commit")

    return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num))        
