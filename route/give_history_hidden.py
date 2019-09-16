from .tool.func import *

def give_history_hidden_2(name):
    

    num = number_check(flask.request.args.get('num', '1'))

    if admin_check(6, 'history_hidden (' + name + '#' + num + ')') == 1:
        sqlQuery("select title from history where title = ? and id = ? and hide = 'O'", [name, num])
        if sqlQuery("fetchall"):
            sqlQuery("update history set hide = '' where title = ? and id = ?", [name, num])
        else:
            sqlQuery("update history set hide = 'O' where title = ? and id = ?", [name, num])
            
        sqlQuery("commit")
    
    return redirect('/history/' + url_pas(name))