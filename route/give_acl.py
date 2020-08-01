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
        curs.execute(db_change("select title from acl where title = ?"), [name])
        if curs.fetchall():
            curs.execute(db_change("update acl set decu = ? where title = ?"), [flask.request.form.get('decu', ''), name])
            curs.execute(db_change("update acl set dis = ? where title = ?"), [flask.request.form.get('dis', ''), name])
            curs.execute(db_change("update acl set why = ? where title = ?"), [flask.request.form.get('why', ''), name])
            curs.execute(db_change("update acl set view = ? where title = ?"), [flask.request.form.get('view', ''), name])
        else:
            curs.execute(db_change("insert into acl (title, decu, dis, why, view) values (?, ?, ?, ?, ?)"), [
                name,
                flask.request.form.get('decu', ''),
                flask.request.form.get('dis', ''),
                flask.request.form.get('why', ''),
                flask.request.form.get('view', '')
            ])

        curs.execute(db_change("select title from acl where title = ? and decu = '' and dis = '' and view = ''"), [name])
        if curs.fetchall():
            curs.execute(db_change("delete from acl where title = ?"), [name])

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
        data = '' + \
            '<h2>' + load_lang('document_acl') + '</h2>' + \
            '<hr class="main_hr">' + \
            '<select name="decu" ' + check_ok + '>' + \
        ''

        if re.search(r'^user:', name):
            acl_list = get_acl_list('user')
        else:
            acl_list = get_acl_list()

        curs.execute(db_change("select decu from acl where title = ?"), [name])
        acl_data = curs.fetchall()
        for data_list in acl_list:
            if acl_data and acl_data[0][0] == data_list:
                check = 'selected="selected"'
            else:
                check = ''

            data += '<option value="' + data_list + '" ' + check + '>' + (data_list if data_list != '' else 'normal') + '</option>'

        data += '</select>'

        if not re.search(r'^user:', name):
            data += '' + \
                '<hr class="main_hr">' + \
                '<h2>' + load_lang('discussion_acl') + '</h2>' + \
                '<hr class="main_hr">' + \
                '<select name="dis" ' + check_ok + '>' + \
            ''

            curs.execute(db_change("select dis, why, view from acl where title = ?"), [name])
            acl_data = curs.fetchall()
            for data_list in acl_list:
                if acl_data and acl_data[0][0] == data_list:
                    check = 'selected="selected"'
                else:
                    check = ''

                data += '<option value="' + data_list + '" ' + check + '>' + (data_list if data_list != '' else 'normal') + '</option>'

            data += '</select>'

            data += '' + \
                '<hr class="main_hr">' + \
                '<h2>' + load_lang('view_acl') + '</h2>' + \
                '<hr class="main_hr">' + \
                '<select name="view" ' + check_ok + '>' + \
            ''
            for data_list in acl_list:
                if acl_data and acl_data[0][2] == data_list:
                    check = 'selected="selected"'
                else:
                    check = ''

                data += '<option value="' + data_list + '" ' + check + '>' + (data_list if data_list != '' else 'normal') + '</option>'

            data += '''
                </select>
                <hr class="main_hr">
                <h2 id="exp">''' + load_lang('explanation') + '''</h2>
                <ul>
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
                </ul>
            '''

            if acl_data:
                data += '' + \
                    '<hr class="main_hr">' + \
                    '<input value="' + html.escape(acl_data[0][1]) + '" placeholder="' + load_lang('why') + '" name="why" type="text" ' + check_ok + '>' + \
                ''
            else:
                data += '' + \
                    '<hr class="main_hr">' + \
                    '<input placeholder="' + load_lang('why') + '" name="why" type="text" ' + check_ok + '>' + \
                ''
                

        return easy_minify(flask.render_template(skin_check(),
            imp = [name, wiki_set(), custom(), other2(['(' + load_lang('acl') + ')', 0])],
            data = '''
                <form method="post">
                    <a href="/setting/8">(''' + load_lang('main_acl_setting') + ''')</a>
                    ''' + data + '''
                    <hr class="main_hr">
                    <button type="submit" ''' + check_ok + '''>''' + load_lang('save') + '''</button>
                </form>
            ''',
            menu = [['w/' + url_pas(name), load_lang('document')], ['manager', load_lang('admin')], ['admin_log?search=' + url_pas('acl (' + name + ')'), load_lang('acl_record')]]
        ))
