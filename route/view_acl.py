from .tool.func import *

def view_acl(name):
    with get_db_connect() as conn:
        curs = conn.cursor()

        check_ok = ''
        user_page = 0
        ip = ip_check()

        if flask.request.method == 'POST':
            check_data = 'document_set (' + name + ')'
        else:
            check_data = None

        user_data = re.search(r'^user:(.+)$', name)
        if user_data:
            if check_data and ip_or_user(ip) != 0:
                return redirect('/login')

            if user_data.group(1) != ip:
                if admin_check(5) != 1:
                    if check_data:
                        return re_error('/error/3')
                    else:
                        check_ok = 'disabled'
            else:
                user_page = 1
        else:
            if admin_check(5) != 1:
                if check_data:
                    return re_error('/error/3')
                else:
                    check_ok = 'disabled'

        if flask.request.method == 'POST':
            acl_data = ['decu', 'document_edit_acl', 'document_move_acl', 'document_delete_acl', 'dis', 'view', 'why']

            for i in acl_data:
                form_data = flask.request.form.get(i, '')
            
                curs.execute(db_change("delete from acl where title = ? and type = ?"), [name, i])
                curs.execute(db_change("insert into acl (title, data, type) values (?, ?, ?)"), [name, form_data, i])
                
                curs.execute(db_change("delete from data_set where doc_name = ? and doc_rev = ? and set_name = 'acl_date'"), [name, i])
                    
                time_limit = flask.request.form.get(i + '_date', '')
                if re.search(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$', time_limit):
                    curs.execute(db_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, 'acl_date', ?)"), [name, i, time_limit])

            all_d = ''
            for i in acl_data[:-1]:
                if flask.request.form.get(i, '') == '':
                    all_d += 'normal'
                    if i != 'view':
                        all_d += ' | '
                else:
                    all_d += flask.request.form.get(i, '')
                    if i != 'view':
                        all_d += ' | '

            markup_data = flask.request.form.get('document_markup', '')

            curs.execute(db_change("select set_data from data_set where doc_name = ? and set_name = 'document_markup'"), [name])
            db_data = curs.fetchall()

            curs.execute(db_change("delete from data_set where doc_name = ? and set_name = 'document_markup'"), [name])
            curs.execute(db_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'document_markup', ?)"), [name, markup_data])

            if not db_data or db_data[0][0] != markup_data:
                curs.execute(db_change("select data from data where title = ?"), [name])
                db_data_2 = curs.fetchall()
                if db_data_2:
                    render_set(
                        doc_name = name,
                        doc_data = db_data_2[0][0],
                        data_type = 'backlink'
                    )

            markup_data = markup_data if markup_data != '' else 'normal'

            if user_page == 1:
                admin_check(5, check_data + ' (' + all_d + ')' + ' (' + markup_data + ')')

            conn.commit()

            return redirect('/acl/' + url_pas(name))
        else:
            data = '<h2>' + load_lang('acl') + '</h2>'
            acl_list = get_acl_list('user') if re.search(r'^user:', name) else get_acl_list()
            if not re.search(r'^user:', name):
                acl_get_list = [
                    [load_lang('view_acl'), 'view', '3'],
                    [load_lang('document_acl'), 'decu', '4'],
                    [load_lang('document_edit_acl'), 'document_edit_acl', '5'],
                    [load_lang('document_move_acl'), 'document_move_acl', '5'],
                    [load_lang('document_delete_acl'), 'document_delete_acl', '5'],
                    [load_lang('discussion_acl'), 'dis', '3'],
                ]
            else:
                acl_get_list = [
                    [load_lang('document_acl'), 'decu', '2']
                ]

            for i in acl_get_list:
                data += '' + \
                    '<h' + i[2] + '>' + i[0] + '</h' + i[2] + '>' + \
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
                
                date_value = ''
                
                curs.execute(db_change("select set_data from data_set where doc_name = ? and doc_rev = ? and set_name = 'acl_date'"), [name, i[1]])
                db_data = curs.fetchall()
                if db_data:
                    date_value = db_data[0][0]
                
                data += '<input type="date" ' + check_ok + ' value="' + date_value + '" name="' + i[1] + '_date" pattern="\\d{4}-\\d{2}-\\d{2}">'
                data += '<hr class="main_hr">'

            curs.execute(db_change("select data from acl where title = ? and type = ?"), [name, 'why'])
            acl_data = curs.fetchall()
            acl_why = html.escape(acl_data[0][0]) if acl_data else ''
            data += '' + \
                '<h3>' + load_lang('why') + '</h3>' + \
                '<input value="' + acl_why + '" ' + check_ok + ' placeholder="' + load_lang('why') + '" name="why" ' + check_ok + '>' + \
                '<hr class="main_hr">' + \
            ''

            data += '''
                <h3>''' + load_lang('explanation') + '''</h3>
                <span id="exp"></span>
                <ul class="opennamu_ul">
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
                    <li>90_day : ''' + load_lang('90_day_acl') + '''</li>
                </ul>
                <hr class="main_hr">
                <h2>''' + load_lang('markup') + '''</h2>
            '''

            curs.execute(db_change("select set_data from data_set where doc_name = ? and set_name = 'document_markup'"), [name])
            db_data = curs.fetchall()
            markup_load = db_data[0][0] if db_data and db_data[0][0] != '' else ''

            markup_list = ['normal'] + get_init_set_list('markup')['list']
            markup_html = ''
            for for_a in markup_list:
                if markup_load == for_a:
                    markup_html = '<option value="' + (for_a if for_a != 'normal' else '') + '">' + for_a + '</option>' + markup_html
                else:
                    markup_html += '<option value="' + (for_a if for_a != 'normal' else '') + '">' + for_a + '</option>'
            
            markup_html = '<select name="document_markup" ' + check_ok + '>' + markup_html + '</select>'

            data += markup_html
            data += '<hr class="main_hr">'

            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('acl') + ')', 0])],
                data = '''
                    <form method="post">
                        <a href="/setting/acl">(''' + load_lang('main_acl_setting') + ''')</a>
                        <hr class="main_hr">
                        ''' + render_simple_set(data) + '''
                        <button type="submit" ''' + check_ok + '''>''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [
                    ['w/' + url_pas(name), load_lang('return')], 
                    ['manager', load_lang('admin')], 
                    ['list/admin/auth_use_page/1/' + url_pas('acl (' + name + ')'), load_lang('acl_record')]
                ]
            ))