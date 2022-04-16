from .tool.func import *

def give_acl_2(conn, name):
    curs = conn.cursor()

    check_ok = ''
    ip = ip_check()

    if flask.request.method == 'POST':
        check_data = 'acl (' + name + ')'
    else:
        check_data = None

    user_data = re.search(r'^user:(.+)$', name)
    if user_data:
        if check_data and ip_or_user(ip) != 0:
            return redirect('/login')

        if user_data.group(1) != ip_check():
            if admin_check(5) != 1:
                if check_data:
                    return re_error('/error/3')
                else:
                    check_ok = 'disabled'
    else:
        if admin_check(5) != 1:
            if check_data:
                return re_error('/error/3')
            else:
                check_ok = 'disabled'

    if flask.request.method == 'POST':
        acl_data = [['decu', flask.request.form.get('decu', '')]]
        acl_data += [['dis', flask.request.form.get('dis', '')]]
        acl_data += [['view', flask.request.form.get('view', '')]]
        acl_data += [['why', flask.request.form.get('why', '')]]
        
        curs.execute(db_change("select title from acl where title = ?"), [name])
        if curs.fetchall():
            for i in acl_data:
                curs.execute(db_change("update acl set data = ? where title = ? and type = ?"), [i[1], name, i[0]])
        else:
            for i in acl_data:
                curs.execute(db_change("insert into acl (title, data, type) values (?, ?, ?)"), [name, i[1], i[0]])

        all_d = ''
        for i in ['decu', 'dis', 'view']:
            if flask.request.form.get(i, '') == '':
                all_d += 'normal'
                if i != 'view':
                    all_d += ' | '
            else:
                all_d += flask.request.form.get(i, '')
                if i != 'view':
                    all_d += ' | '

        admin_check(5, check_data + ' (' + all_d + ')')

        conn.commit()

        return redirect('/acl/' + url_pas(name))
    else:
        data = ''
        acl_list = get_acl_list('user') if re.search(r'^user:', name) else get_acl_list()
        if not re.search(r'^user:', name):
            acl_get_list = [
                [load_lang('document_acl'), 'decu'], 
                [load_lang('discussion_acl'), 'dis'], 
                [load_lang('view_acl'), 'view']
            ]
        else:
            acl_get_list = [
                [load_lang('document_acl'), 'decu']
            ]
            
        for i in acl_get_list:
            data += '' + \
                '<h2>' + i[0] + '</h2>' + \
                '<hr class="main_hr">' + \
                '<select name="' + i[1] + '" ' + check_ok + '>' + \
            ''
    
            curs.execute(db_change("select data from acl where title = ? and type = ?"), [name, i[1]])
            acl_data = curs.fetchall()
            for data_list in acl_list:
                check = 'selected="selected"' if acl_data and acl_data[0][0] == data_list else ''
                data += '<option value="' + data_list + '" ' + check + '>' + (data_list if data_list != '' else 'normal') + '</option>'
    
            data += '</select>'
            data += '<hr class="main_hr">'

        curs.execute(db_change("select data from acl where title = ? and type = ?"), [name, 'why'])
        acl_data = curs.fetchall()
        acl_why = html.escape(acl_data[0][0]) if acl_data else ''
        data += '' + \
            '<hr class="main_hr">' + \
            '<input value="' + acl_why + '" placeholder="' + load_lang('why') + '" name="why" type="text" ' + check_ok + '>' + \
        ''

        data += '''
            <h2 id="exp">''' + load_lang('explanation') + '''</h2>
            <ul class="inside_ul">
                <li>normal : ''' + load_lang('unset') + '''</li>
                <li>admin : ''' + load_lang('admin_acl') + '''</li>
                <li>user : ''' + load_lang('member_acl') + '''</li>
                <li>50_edit : ''' + load_lang('50_edit_acl') + '''</li>
                <li>all : ''' + load_lang('all_acl') + '''</li>
                <li>email : ''' + load_lang('email_acl') + '''</li>
                <li>owner : ''' + load_lang('owner_acl') + '''</li>
                <li>ban : ''' + load_lang('ban_acl') + '''</li>
                <li>before : ''' + load_lang('before_acl') + '''</li>
                <li>30_day : ''' + load_lang('30_day_acl') + '''</li>
                <li>ban_admin : ''' + load_lang('ban_admin_acl') + '''</li>
                <li>not_all : ''' + load_lang('not_all_acl') + '''</li>
            </ul>
        '''

        return easy_minify(flask.render_template(skin_check(),
            imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('acl') + ')', 0])],
            data = '''
                <form method="post">
                    <a href="/setting/8">(''' + load_lang('main_acl_setting') + ''')</a>
                    ''' + data + '''
                    <button type="submit" ''' + check_ok + '''>''' + load_lang('save') + '''</button>
                </form>
            ''',
            menu = [
                ['w/' + url_pas(name), load_lang('document')], 
                ['manager', load_lang('admin')], 
                ['admin_log?search=' + url_pas('acl (' + name + ')'), load_lang('acl_record')]
            ]
        ))
