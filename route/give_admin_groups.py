from .tool.func import *

def give_admin_plus_2(conn, name):
    curs = conn.cursor()
    
    if flask.request.method == 'POST':
        if admin_check(None, 'admin_plus (' + name + ')') != 1:
            return re_error('/error/3')

        curs.execute("delete from alist where name = ?", [name])
        
        if flask.request.form.get('ban', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'ban')", [name])

        if flask.request.form.get('toron', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'toron')", [name])
            
        if flask.request.form.get('check', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'check')", [name])

        if flask.request.form.get('acl', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'acl')", [name])

        if flask.request.form.get('hidel', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'hidel')", [name])

        if flask.request.form.get('give', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'give')", [name])

        if flask.request.form.get('owner', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'owner')", [name])
            
        conn.commit()
        
        return redirect('/admin_plus/' + url_pas(name))
    else:        
        data = '<ul>'
        
        exist_list = ['', '', '', '', '', '', '', '']

        curs.execute('select acl from alist where name = ?', [name])
        acl_list = curs.fetchall()    
        for go in acl_list:
            if go[0] == 'ban':
                exist_list[0] = 'checked="checked"'
            elif go[0] == 'toron':
                exist_list[2] = 'checked="checked"'
            elif go[0] == 'check':
                exist_list[3] = 'checked="checked"'
            elif go[0] == 'acl':
                exist_list[4] = 'checked="checked"'
            elif go[0] == 'hidel':
                exist_list[5] = 'checked="checked"'
            elif go[0] == 'give':
                exist_list[6] = 'checked="checked"'
            elif go[0] == 'owner':
                exist_list[7] = 'checked="checked"'

        if admin_check() != 1:
            state = 'disabled'
        else:
            state = ''

        data += '''
                <li><input type="checkbox" ''' + state +  ' name="ban" ' + exist_list[0] + '> ' + load_lang('ban_authority') + '''</li>
                <li><input type="checkbox" ''' + state +  ' name="toron" ' + exist_list[2] + '> ' + load_lang('discussion_authority') + '''</li>
                <li><input type="checkbox" ''' + state +  ' name="check" ' + exist_list[3] + '> ' + load_lang('user_check_authority') + '''</li>
                <li><input type="checkbox" ''' + state +  ' name="acl" ' + exist_list[4] + '> ' + load_lang('document_acl_authority') + '''</li>
                <li><input type="checkbox" ''' + state +  ' name="hidel" ' + exist_list[5] + '> ' + load_lang('history_hide_authority') + '''</li>
                <li><input type="checkbox" ''' + state +  ' name="give" ' + exist_list[6] + '> ' + load_lang('authorization_authority') + '''</li>
                <li><input type="checkbox" ''' + state +  ' name="owner" ' + exist_list[7] + '> ' + load_lang('owner_authority') + '''</li>
            </ul>
        '''

        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('admin_group_add'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                <form method="post">
                    ''' + data + '''
                    <hr class=\"main_hr\">
                    <button id="save" ''' + state +  ''' type="submit">''' + load_lang('save') + '''</button>
                </form>
            ''',
            menu = [['manager', load_lang('return')]]
        ))     