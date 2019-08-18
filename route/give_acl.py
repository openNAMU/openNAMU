from .tool.func import *
import pymysql

def give_acl_2(conn, name):
    curs = conn.cursor()

    check_ok = ''
    
    if flask.request.method == 'POST':
        check_data = 'acl (' + name + ')'
    else:
        check_data = None
    
    user_data = re.search('^user:(.+)$', name)
    if user_data:
        if check_data and custom()[2] == 0:
            return redirect('/login')
        
        if user_data.groups()[0] != ip_check():
            if admin_check(5, check_data) != 1:
                if check_data:
                    return re_error('/error/3')
                else:
                    check_ok = 'disabled'
    else:
        if admin_check(5, check_data) != 1:
            if check_data:
                return re_error('/error/3')
            else:
                check_ok = 'disabled'

    if flask.request.method == 'POST':
        decu = flask.request.form.get('decu', '')
        view = flask.request.form.get('view', '')

        curs.execute("select title from acl where title = %s", [name])
        if curs.fetchall():
            curs.execute("update acl set decu = %s where title = %s", [decu, name])
            curs.execute("update acl set dis = %s where title = %s", [flask.request.form.get('dis', ''), name])
            curs.execute("update acl set why = %s where title = %s", [flask.request.form.get('why', ''), name])
            curs.execute("update acl set view = %s where title = %s", [view, name])
        else:
            curs.execute("insert into acl (title, decu, dis, why, view) values (%s, %s, %s, %s, %s)", [
                name, 
                decu, 
                flask.request.form.get('dis', ''), 
                flask.request.form.get('why', ''), 
                view
            ])
        
        curs.execute("select title from acl where title = %s and decu = '' and dis = '' and view = ''", [name])
        if curs.fetchall():
            curs.execute("delete from acl where title = %s", [name])

        conn.commit()
            
        return redirect('/acl/' + url_pas(name))            
    else:
        data = '<h2>' + load_lang('document_acl') + '</h2><hr class=\"main_hr\"><select name="decu" ' + check_ok + '>'
    
        if re.search('^user:', name):
            acl_list = [['', 'normal'], ['user', 'member'], ['all', 'all']]
        else:
            acl_list = [['', 'normal'], ['user', 'member'], ['admin', 'admin'], ['50_edit', '50 edit'], ['email', 'email']]
        
        curs.execute("select decu from acl where title = %s", [name])
        acl_data = curs.fetchall()
        for data_list in acl_list:
            if acl_data and acl_data[0][0] == data_list[0]:
                check = 'selected="selected"'
            else:
                check = ''
            
            data += '<option value="' + data_list[0] + '" ' + check + '>' + data_list[1] + '</option>'
            
        data += '</select>'
        
        if not re.search('^user:', name):
            data += '<hr class=\"main_hr\"><h2>' + load_lang('discussion_acl') + '</h2><hr class=\"main_hr\"><select name="dis" ' + check_ok + '>'
        
            curs.execute("select dis, why, view from acl where title = %s", [name])
            acl_data = curs.fetchall()
            for data_list in acl_list:
                if acl_data and acl_data[0][0] == data_list[0]:
                    check = 'selected="selected"'
                else:
                    check = ''
                    
                data += '<option value="' + data_list[0] + '" ' + check + '>' + data_list[1] + '</option>'
                
            data += '</select>'

            data += '<hr class=\"main_hr\"><h2>' + load_lang('view_acl') + '</h2><hr class=\"main_hr\"><select name="view" ' + check_ok + '>'
            for data_list in acl_list:
                if acl_data and acl_data[0][2] == data_list[0]:
                    check = 'selected="selected"'
                else:
                    check = ''
                    
                data += '<option value="' + data_list[0] + '" ' + check + '>' + data_list[1] + '</option>'
                
            data += '''
                </select>
                <h2>''' + load_lang('explanation') + '''</h2>
                <ul>
                    <li>normal : ''' + load_lang('default') + '''</li>
                    <li>admin : ''' + load_lang('admin_acl') + '''</li>
                    <li>member : ''' + load_lang('member_acl') + '''</li>
                    <li>50 edit : ''' + load_lang('50_edit_acl') + '''</li>
                    <li>all : ''' + load_lang('all_acl') + '''</li>
                    <li>email : ''' + load_lang('email_acl') + '''</li>
                </ul>
            '''
                
            if check_ok == '':
                if acl_data:
                    data += '<hr class=\"main_hr\"><input value="' + html.escape(acl_data[0][1]) + '" placeholder="' + load_lang('why') + '" name="why" type="text" ' + check_ok + '>'
                else:
                    data += '<hr class=\"main_hr\"><input placeholder="' + load_lang('why') + '" name="why" type="text" ' + check_ok + '>'
            
        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('acl') + ')', 0])],
            data =  '''
                <form method="post">
                    ''' + data + '''
                    <hr class=\"main_hr\">
                    <button type="submit" ''' + check_ok + '''>''' + load_lang('save') + '''</button>
                </form>
            ''',
            menu = [['w/' + url_pas(name), load_lang('document')], ['manager', load_lang('admin')]]
        ))