from .tool.func import *

def inter_wiki_del_2(conn, tools, name):
    curs = conn.cursor()

    if admin_check(None, tools) == 1:
        if tools == 'del_inter_wiki':
            curs.execute(db_change("delete from inter where title = ?"), [name])
        elif tools == 'del_edit_filter':
            curs.execute(db_change("delete from filter where name = ?"), [name])
        elif tools == 'del_name_filter':
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'name'"), [name])
        elif tools == 'del_file_filter':
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'file'"), [name])
        elif tools == 'del_email_filter':
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'email'"), [name])
        elif tools == 'del_image_license':
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'image_license'"), [name])
        else:
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'edit_top'"), [name])

        conn.commit()

        return redirect('/' + re.sub('^del_', '', tools))
    else:
        return re_error('/error/3')