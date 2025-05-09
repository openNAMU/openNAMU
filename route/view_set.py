from .tool.func import *

def view_set_markup(conn, document_name = '', markup = '', addon = '', disable = ''):
    curs = conn.cursor()

    curs.execute(db_change('select data from other where name = "markup"'))
    db_data = curs.fetchall()
    default_markup = db_data[0][0] if db_data and db_data[0][0] != '' else 'namumark'

    markup_load = markup
    if markup == '':
        curs.execute(db_change("select set_data from data_set where doc_name = ? and set_name = 'document_markup'"), [document_name])
        db_data = curs.fetchall()
        markup_load = db_data[0][0] if db_data and db_data[0][0] != '' else ''

    markup_list = ['normal'] + get_init_set_list('markup')['list']
    markup_html = ''
    for for_a in markup_list:
        if markup_load == for_a:
            markup_html = '<option value="' + (for_a if for_a != 'normal' else default_markup) + '">' + for_a + '</option>' + markup_html
        else:
            markup_html += '<option value="' + (for_a if for_a != 'normal' else default_markup) + '">' + for_a + '</option>'
    
    markup_html = '<select name="document_markup" ' + disable + ' ' + addon + '>' + markup_html + '</select>'

    return markup_html

async def view_set(name = 'Test', multiple = False):
    with get_db_connect() as conn:
        curs = conn.cursor()

        check_ok = ''
        ip = ip_check()
        time = get_time()

        if multiple and flask.request.method == 'POST':
            all_title = re.findall(r'([^\n]+)\n', flask.request.form.get('title_name', '').replace('\r', '') + '\n')
            for name in all_title:
                view_set(name, False)

            return redirect(conn, '/list/document/acl')
        else:
            if flask.request.method == 'POST':
                check_data = 'document_set (' + name + ')'
            else:
                check_data = ''

            need_admin = True

            user_data = re.search(r'^user:([^/]+)', name)
            if user_data:
                if ip_or_user(ip) != 0:
                    return redirect(conn, '/login')

                if user_data.group(1) == ip:
                    need_admin = False
            
            if need_admin:
                if await acl_check(tool = 'acl_auth') == 1:
                    if flask.request.method == 'POST':
                        return await re_error(conn, 3)
                    else:
                        check_ok = 'disabled'

        if flask.request.method == 'POST':
            acl_data = ['decu', 'document_edit_acl', 'document_edit_request_acl', 'document_move_acl', 'document_delete_acl', 'dis', 'view', 'why']
            acl_result = []
            acl_text = ''

            for i in acl_data:
                form_data = flask.request.form.get(i, '')
                
                acl_result += [form_data]

                acl_text += i + '\n'
                acl_text += form_data + '\n'
            
                curs.execute(db_change("delete from acl where title = ? and type = ?"), [name, i])
                curs.execute(db_change("insert into acl (title, data, type) values (?, ?, ?)"), [name, form_data, i])
                
                curs.execute(db_change("delete from data_set where doc_name = ? and doc_rev = ? and set_name = 'acl_date'"), [name, i])
                    
                time_limit = flask.request.form.get(i + '_date', '')
                if re.search(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$', time_limit):
                    curs.execute(db_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, 'acl_date', ?)"), [name, i, time_limit])
                    
                    acl_text += time_limit + '\n'

                acl_text += '\n\n'

            markup_data = flask.request.form.get('document_markup', '')
            
            acl_text += 'document_markup\n'
            acl_text += markup_data + '\n\n'

            curs.execute(db_change("select set_data from data_set where doc_name = ? and set_name = 'document_markup'"), [name])
            db_data = curs.fetchall()

            curs.execute(db_change("delete from data_set where doc_name = ? and set_name = 'document_markup'"), [name])
            curs.execute(db_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'document_markup', ?)"), [name, markup_data])

            if not db_data or db_data[0][0] != markup_data:
                curs.execute(db_change("select data from data where title = ?"), [name])
                db_data_2 = curs.fetchall()
                if db_data_2:
                    render_set(conn, 
                        doc_name = name,
                        doc_data = db_data_2[0][0],
                        data_type = 'backlink'
                    )

            markup_data = markup_data if markup_data != '' else 'normal'

            if await acl_check('', 'owner_auth', '', '') != 1:
                document_top = flask.request.form.get('document_top', '')

                acl_text += 'document_top\n'
                acl_text += document_top + '\n\n'

                curs.execute(db_change("delete from data_set where doc_name = ? and set_name = 'document_top'"), [name])
                curs.execute(db_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'document_top', ?)"), [name, document_top])
                
                document_editor_top = flask.request.form.get('document_editor_top', '')

                acl_text += 'document_editor_top\n'
                acl_text += document_editor_top + '\n\n'

                curs.execute(db_change("delete from data_set where doc_name = ? and set_name = 'document_editor_top'"), [name])
                curs.execute(db_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'document_editor_top', ?)"), [name, document_editor_top])

            if need_admin:
                await acl_check(tool = 'acl_auth', memo = check_data)

            history_plus(conn, 
                name,
                acl_text,
                time,
                ip,
                acl_result[7],
                '0',
                mode = 'setting'
            )

            return redirect(conn, '/acl/' + url_pas(name))
        else:
            data = '<h2>' + get_lang(conn, 'acl') + '</h2>'
            acl_list = await get_acl_list()
            acl_get_list = [
                [get_lang(conn, 'view_acl'), 'view', '3'],
                [get_lang(conn, 'document_acl'), 'decu', '4'],
                [get_lang(conn, 'document_edit_acl'), 'document_edit_acl', '5'],
                [get_lang(conn, 'document_edit_request_acl'), 'document_edit_request_acl', '5'],
                [get_lang(conn, 'document_move_acl'), 'document_move_acl', '5'],
                [get_lang(conn, 'document_delete_acl'), 'document_delete_acl', '5'],
                [get_lang(conn, 'discussion_acl'), 'dis', '3'],
            ]

            for i in acl_get_list:
                data += '' + \
                    '<h' + i[2] + '>' + i[0] + '</h' + i[2] + '>' + \
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
                '<h3>' + get_lang(conn, 'why') + '</h3>' + \
                '<input value="' + acl_why + '" ' + check_ok + ' placeholder="' + get_lang(conn, 'why') + '" name="why" ' + check_ok + '>' + \
                '<hr class="main_hr">' + \
            ''

            data += '''
                <h3>''' + get_lang(conn, 'explanation') + '''</h3>
                <span id="exp"></span>
                <ul>
                    <li>normal : ''' + get_lang(conn, 'unset') + '''</li>
                    <li>admin : ''' + get_lang(conn, 'admin_acl') + '''</li>
                    <li>user : ''' + get_lang(conn, 'member_acl') + '''</li>
                    <li>50_edit : ''' + get_lang(conn, '50_edit_acl') + '''</li>
                    <li>all : ''' + get_lang(conn, 'all_acl') + '''</li>
                    <li>email : ''' + get_lang(conn, 'email_acl') + '''</li>
                    <li>owner : ''' + get_lang(conn, 'owner_acl') + '''</li>
                    <li>ban : ''' + get_lang(conn, 'ban_acl') + '''</li>
                    <li>before : ''' + get_lang(conn, 'before_acl') + '''</li>
                    <li>30_day : ''' + get_lang(conn, '30_day_acl') + '''</li>
                    <li>ban_admin : ''' + get_lang(conn, 'ban_admin_acl') + '''</li>
                    <li>not_all : ''' + get_lang(conn, 'not_all_acl') + '''</li>
                    <li>90_day : ''' + get_lang(conn, '90_day_acl') + '''</li>
                    <li>up_to_level_3 : ''' + get_lang(conn, 'up_to_level_3') + '''</li>
                    <li>up_to_level_10 : ''' + get_lang(conn, 'up_to_level_10') + '''</li>
                </ul>
                <h2>''' + get_lang(conn, 'markup') + '''</h2>
            '''

            data += view_set_markup(conn, document_name = name, disable = check_ok)

            save_button = '<button type="submit" ' + check_ok + '>' + get_lang(conn, 'save') + '</button>'
            if await acl_check('', 'owner_auth', '', '') == 1:
                check_ok = 'disabled'

            curs.execute(db_change("select set_data from data_set where doc_name = ? and set_name = 'document_top'"), [name])
            db_data = curs.fetchall()
            document_top = db_data[0][0] if db_data and db_data[0][0] != '' else ''

            curs.execute(db_change("select set_data from data_set where doc_name = ? and set_name = 'document_editor_top'"), [name])
            db_data = curs.fetchall()
            document_editor_top = db_data[0][0] if db_data and db_data[0][0] != '' else ''

            data += '''
                <h2>''' + get_lang(conn, 'document_top') + ''' (HTML)</h2>
                <textarea ''' + check_ok + ''' class="opennamu_textarea_100" name="document_top">''' + html.escape(document_top) + '''</textarea>
                
                <h2>''' + get_lang(conn, 'document_editor_top') + ''' (HTML)</h2>
                <textarea ''' + check_ok + ''' class="opennamu_textarea_100" name="document_editor_top">''' + html.escape(document_editor_top) + '''</textarea>
            '''
            data += '<hr class="main_hr">'

            text_area = ''
            if multiple == True:
                text_area = '<textarea class="opennamu_textarea_500" placeholder="' + get_lang(conn, 'many_delete_help') + '" name="title_name"></textarea><hr class="main_hr">'
                menu = [
                    ['manager', get_lang(conn, 'admin')]
                ]
                title = get_lang(conn, 'mutiple_document_setting')
                sub = 0
            else:
                menu = [
                    ['w/' + url_pas(name), get_lang(conn, 'return')], 
                    ['acl_multiple', get_lang(conn, 'mutiple_document_setting')],
                    ['manager', get_lang(conn, 'admin')]
                ]
                title = name
                sub = '(' + get_lang(conn, 'document_setting') + ')'
                save_button += ' <button type="button" onclick="w_set_reset();" ' + check_ok + '>' + get_lang(conn, 'reset') + '</button>'

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [title, await wiki_set(), await wiki_custom(conn), wiki_css([sub, 0])],
                data = '''
                    <form method="post">
                        <script defer src="/views/main_css/js/route/w_set.js''' + cache_v() + '''"></script>
                        ''' + text_area + '''
                        ''' + render_simple_set(conn, data) + '''
                        ''' + save_button + '''
                    </form>
                ''',
                menu = menu
            ))