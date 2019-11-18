from .tool.func import *

def give_admin_2(conn, name):
    curs = conn.cursor()

    owner = admin_check()
    
    curs.execute(db_change("select acl from user where id = ?"), [name])
    user = curs.fetchall()
    if not user:
        return re_error('/error/2')
    else:
        if owner != 1:
            curs.execute(db_change('select name from alist where name = ? and acl = "owner"'), [user[0][0]])
            if curs.fetchall():
                return re_error('/error/3')

            if ip_check() == name:
                return re_error('/error/3')

    if flask.request.method == 'POST':
        if admin_check(7, 'admin (' + name + ')') != 1:
            return re_error('/error/3')

        if owner != 1:
            curs.execute(db_change('select name from alist where name = ? and acl = "owner"'), [flask.request.form.get('select', None)])
            if curs.fetchall():
                return re_error('/error/3')

        if flask.request.form.get('select', None) == 'X':
            curs.execute(db_change("update user set acl = 'user' where id = ?"), [name])
        else:
            curs.execute(db_change("update user set acl = ? where id = ?"), [flask.request.form.get('select', None), name])
        
        conn.commit()
        
        return redirect('/admin/' + url_pas(name))            
    else:
        if admin_check(7) != 1:
            return re_error('/error/3')            

        div = '<option value="X">X</option>'
        
        curs.execute(db_change('select distinct name from alist order by name asc'))
        for data in curs.fetchall():
            if user[0][0] == data[0]:
                div += '<option value="' + data[0] + '" selected="selected">' + data[0] + '</option>'
            else:
                if owner != 1:
                    curs.execute(db_change('select name from alist where name = ? and acl = "owner"'), [data[0]])
                    if not curs.fetchall():
                        div += '<option value="' + data[0] + '">' + data[0] + '</option>'
                else:
                    div += '<option value="' + data[0] + '">' + data[0] + '</option>'
        
        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('authorize') + ')', 0])],
            data =  '''
                    <form method="post">
                        <select name="select">''' + div + '''</select>
                        <hr class=\"main_hr\">
                        <button type="submit">''' + load_lang('save') + '''</button>
                    </form>
                    ''',
            menu = [['manager', load_lang('return')]]
        ))