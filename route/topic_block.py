from .tool.func import *

def topic_block_2(name, sub, num):
    

    if admin_check(3, 'blind (' + name + ' - ' + sub + '#' + str(num) + ')') != 1:
        return re_error('/error/3')

    sqlQuery("select block from topic where title = ? and sub = ? and id = ?", [name, sub, str(num)])
    block = sqlQuery("fetchall")
    if block:
        if block[0][0] == 'O':
            sqlQuery("update topic set block = '' where title = ? and sub = ? and id = ?", [name, sub, str(num)])
        else:
            sqlQuery("update topic set block = 'O' where title = ? and sub = ? and id = ?", [name, sub, str(num)])
        
        rd_plus(name, sub, get_time())
        
        sqlQuery("commit")
        
    return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num))