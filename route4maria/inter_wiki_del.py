from .tool.func import *
import pymysql

def inter_wiki_del_2(conn, tools, name):
    curs = conn.cursor()
    
    if admin_check(None, tools) == 1:
        if tools == 'del_inter_wiki':
            curs.execute("delete from inter where title = %s", [name])
        elif tools == 'del_edit_filter':
            curs.execute("delete from filter where name = %s", [name])
        elif tools == 'del_name_filter':
            curs.execute("delete from html_filter where html = %s and kind = 'name'", [name])
        elif tools == 'del_file_filter':
            curs.execute("delete from html_filter where html = %s and kind = 'file'", [name])
        else:
            curs.execute("delete from html_filter where html = %s and kind = 'email'", [name])
        
        

        return redirect('/' + re.sub('^del_', '', tools))
    else:
        return re_error('/error/3')