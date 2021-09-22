from .tool.func import *

def give_admin_2(conn, name):
    curs = conn.cursor()

    owner = admin_check()

    curs.execute(db_change("select data from user_set where id = ? and name = 'acl'"), [name])
    user_acl = curs.fetchall()
    if not user_acl:
        return re_error('/error/2')
    else:
        user_acl = user_acl[0][0]
    
    if owner != 1:
        curs.execute(db_change('select name from alist where name = ? and acl = "owner"'), [user_acl])
        if curs.fetchall():
            return re_error('/error/3')

        if ip_check() == name:
            return re_error('/error/3')

    if flask.request.method == 'POST':
        if admin_check(7, 'admin (' + name + ')') != 1:
            return re_error('/error/3')

        if flask.request.form.get('select', 'X') == 'X':
            select_data = 'user'
        else:
            select_data = flask.request.form.get('select', 'X')
        
        curs.execute(db_change('select name from alist where name = ? and acl = "owner"'), [select_data])
        if owner != 1 and curs.fetchall():
            return re_error('/error/3')
            
        curs.execute(db_change("update user_set set data = ? where id = ? and name = 'acl'"), [
            select_data, 
            name
        ])

        conn.commit()

        return redirect('/admin/' + url_pas(name))
    else:
        if admin_check(7) != 1:
            return re_error('/error/3')

        div = '<option value="X">X</option>'

        curs.execute(db_change('select distinct name from alist order by name asc'))
        for data in curs.fetchall():
            if user_acl == data[0]:
                div = '<option value="' + data[0] + '">' + data[0] + '</option>' + div
            else:
                div += '<option value="' + data[0] + '">' + data[0] + '</option>'

        return easy_minify(flask.render_template(skin_check(),
            imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('authorize') + ')', 0])],
            data =  '''
                    <form method="post">
                        <select name="select">''' + div + '''</select>
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('save') + '''</button>
                    </form>
                    ''',
            menu = [['manager', load_lang('return')]]
        ))