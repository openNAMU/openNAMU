from .tool.func import *

def del_inter_2(conn, tools, name):
    curs = conn.cursor()
    
    if admin_check(None, tools) == 1:
        if tools == 'del_inter_wiki':
            curs.execute("delete from inter where title = ?", [name])
        elif tools == 'del_edit_filter':
            curs.execute("delete from filter where name = ?", [name])
        elif tools == 'del_name_filter':
            curs.execute("delete from html_filter where html = ? and kind = 'name'", [name])
        else:
            curs.execute("delete from html_filter where html = ? and kind = 'email'", [name])
        
        conn.commit()

        return redirect('/' + re.sub('^del_', '', tools))
    else:
        return re_error('/error/3')