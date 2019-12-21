from .tool.func import *

def setting_2(conn, num):
    curs = conn.cursor()

    if num != 0 and admin_check() != 1:
        return re_error('/ban')

    if num == 0:
        li_list = [
            load_lang('main_setting'),
            load_lang('text_setting'),
            load_lang('main_head'),
            load_lang('main_body'),
            'robots.txt',
            'Google',
            load_lang('main_bottom_body'),
            load_lang('main_acl_setting')
        ]

        x = 0
        li_data = ''

        for li in li_list:
            x += 1
            li_data += '<li><a href="/setting/' + str(x) + '">' + li + '</a></li>'

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('setting'), wiki_set(), custom(), other2([0, 0])],
            data = '<h2>' + load_lang('list') + '</h2><ul>' + li_data + '</ul>',
            menu = [['manager', load_lang('return')]]
        ))
    elif num == 1:
        i_list = {
            0 : 'name',
            1 : 'logo',
            2 : 'frontpage',
            3 : 'license',
            4 : 'upload',
            5 : 'skin',
            7 : 'reg',
            8 : 'ip_view',
            9 : 'back_up',
            10 : 'port',
            11 : 'key',
            12 : 'update',
            13 : 'email_have',
            15 : 'encode',
            16 : 'host',
            19 : 'slow_edit'
        }
        n_list = {
            0 : 'Wiki',
            1 : '',
            2 : 'FrontPage',
            3 : 'CC 0',
            4 : '2',
            5 : '',
            7 : '',
            8 : '',
            9 : '0',
            10 : '3000',
            11 : 'test',
            12 : 'stable',
            13 : '',
            15 : 'sha3',
            16 : '0.0.0.0',
            19 : '0'
        }

        if flask.request.method == 'POST':
            for i in i_list:
                curs.execute(db_change("update other set data = ? where name = ?"), [
                    flask.request.form.get(i_list[i], n_list[i]),
                    i_list[i]
                ])

            conn.commit()

            admin_check(None, 'edit_set')

            return redirect('/setting/1')
        else:
            d_list = {}

            for i in i_list:
                curs.execute(db_change('select data from other where name = ?'), [i_list[i]])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list[i] = sql_d[0][0]
                else:
                    curs.execute(db_change('insert into other (name, data) values (?, ?)'), [i_list[i], n_list[i]])

                    d_list[i] = n_list[i]

            conn.commit()

            acl_div = ['']
            encode_data = ['sha256', 'sha3']
            for acl_data in encode_data:
                if acl_data == d_list[15]:
                    acl_div[0] = '<option value="' + acl_data + '">' + acl_data + '</option>' + acl_div[0]
                else:
                    acl_div[0] += '<option value="' + acl_data + '">' + acl_data + '</option>'

            check_box_div = ['', '', '']
            for i in range(0, 3):
                if i == 0:
                    acl_num = 7
                elif i == 1:
                    acl_num = 8
                else:
                    acl_num = 13

                if d_list[acl_num]:
                    check_box_div[i] = 'checked="checked"'

            branch_div =''
            if d_list[12] == 'stable':
                branch_div += '<option value="stable">stable</option>'
                branch_div += '<option value="master">master</option>'
            else:
                branch_div += '<option value="master">master</option>'
                branch_div += '<option value="stable">stable</option>'

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('main_setting'), wiki_set(), custom(), other2([0, 0])],
                data = '''
                    <form method="post">
                        <span>''' + load_lang('wiki_name') + '''</span>
                        <hr class=\"main_hr\">
                        <input type="text" name="name" value="''' + html.escape(d_list[0]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('wiki_logo') + ''' (HTML)</span>
                        <hr class=\"main_hr\">
                        <input type="text" name="logo" value="''' + html.escape(d_list[1]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('main_page') + '''</span>
                        <hr class=\"main_hr\">
                        <input type="text" name="frontpage" value="''' + html.escape(d_list[2]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('bottom_text') + ''' (HTML)</span>
                        <hr class=\"main_hr\">
                        <input type="text" name="license" value="''' + html.escape(d_list[3]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('max_file_size') + ''' (MB)</span>
                        <hr class=\"main_hr\">
                        <input type="text" name="upload" value="''' + html.escape(d_list[4]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('backup_interval') + ' (' + load_lang('hour') + ') (' + load_lang('off') + ' : 0) (' + load_lang('sqlite_only') + ') (' + load_lang('restart_required') + ''')</span>
                        <hr class=\"main_hr\">
                        <input type="text" name="back_up" value="''' + html.escape(d_list[9]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('wiki_skin') + '''</span>
                        <hr class=\"main_hr\">
                        <select name="skin">''' + load_skin(d_list[5]) + '''</select>
                        <hr class=\"main_hr\">
                        <input type="checkbox" name="reg" ''' + check_box_div[0] + '''> ''' + load_lang('no_register') + '''
                        <hr class=\"main_hr\">
                        <input type="checkbox" name="ip_view" ''' + check_box_div[1] + '''> ''' + load_lang('hide_ip') + '''
                        <hr class=\"main_hr\">
                        <input type="checkbox" name="email_have" ''' + check_box_div[2] + '''> ''' + load_lang('email_required') + ' <a href="/setting/6">(' + load_lang('google_imap_required') + ''')</a>
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('wiki_host') + '''</span>
                        <hr class=\"main_hr\">
                        <input type="text" name="host" value="''' + html.escape(d_list[16]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('wiki_port') + '''</span>
                        <hr class=\"main_hr\">
                        <input type="text" name="port" value="''' + html.escape(d_list[10]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('wiki_secret_key') + '''</span>
                        <hr class=\"main_hr\">
                        <input type="password" name="key" value="''' + html.escape(d_list[11]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('update_branch') + '''</span>
                        <hr class=\"main_hr\">
                        <select name="update">''' + branch_div + '''</select>
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('encryption_method') + '''</span>
                        <hr class=\"main_hr\">
                        <select name="encode">''' + acl_div[0] + '''</select>
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('slow_edit') + ' (' + load_lang('second') + ') (' + load_lang('off') + ''' : 0)</span>
                        <hr class=\"main_hr\">
                        <input name="''' + i_list[19] + '''" value="''' + html.escape(d_list[19]) + '''">
                        <hr class=\"main_hr\">
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['setting', load_lang('return')]]
            ))
    elif num == 2:
        i_list = [
            'contract',
            'no_login_warring',
            'edit_bottom_text',
            'check_key_text',
            'email_title',
            'email_text',
            'email_insert_text',
            'password_search_text',
            'reset_user_text',
            'error_401',
            'error_404',
            'edit_help'
        ]
        if flask.request.method == 'POST':
            for i in i_list:
                curs.execute(db_change("update other set data = ? where name = ?"), [
                    flask.request.form.get(i, ''),
                    i
                ])

            conn.commit()

            admin_check(None, 'edit_set')

            return redirect('/setting/2')
        else:
            d_list = []

            for i in i_list:
                curs.execute(db_change('select data from other where name = ?'), [i])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list += [sql_d[0][0]]
                else:
                    curs.execute(db_change('insert into other (name, data) values (?, ?)'), [i, ''])

                    d_list += ['']

            conn.commit()

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('text_setting'), wiki_set(), custom(), other2([0, 0])],
                data = '''
                    <form method="post">
                        <span>''' + load_lang('register_text') + ''' (HTML)</span>
                        <hr class=\"main_hr\">
                        <input name="''' + i_list[0] + '''" value="''' + html.escape(d_list[0]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('non_login_alert') + ''' (HTML)</span>
                        <hr class=\"main_hr\">
                        <input name="''' + i_list[1] + '''" value="''' + html.escape(d_list[1]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('edit_bottom_text') + ''' (HTML)</span>
                        <hr class=\"main_hr\">
                        <input name="''' + i_list[2] + '''" value="''' + html.escape(d_list[2]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('check_key_text') + ''' (HTML)</span>
                        <hr class=\"main_hr\">
                        <input name="''' + i_list[3] + '''" value="''' + html.escape(d_list[3]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('email_title') + '''</span>
                        <hr class=\"main_hr\">
                        <input name="''' + i_list[4] + '''" value="''' + html.escape(d_list[4]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('email_text') + '''</span>
                        <hr class=\"main_hr\">
                        <input name="''' + i_list[5] + '''" value="''' + html.escape(d_list[5]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('email_insert_text') + '''</span>
                        <hr class=\"main_hr\">
                        <input name="''' + i_list[6] + '''" value="''' + html.escape(d_list[6]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('password_search_text') + '''</span>
                        <hr class=\"main_hr\">
                        <input name="''' + i_list[7] + '''" value="''' + html.escape(d_list[7]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('reset_user_text') + '''</span>
                        <hr class=\"main_hr\">
                        <input name="''' + i_list[8] + '''" value="''' + html.escape(d_list[8]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('error_401') + '''</span>
                        <hr class=\"main_hr\">
                        <input name="''' + i_list[9] + '''" value="''' + html.escape(d_list[9]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('error_404') + '''</span>
                        <hr class=\"main_hr\">
                        <input name="''' + i_list[10] + '''" value="''' + html.escape(d_list[10]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('edit_help') + '''</span>
                        <hr class=\"main_hr\">
                        <input name="''' + i_list[11] + '''" value="''' + html.escape(d_list[11]) + '''">
                        <hr class=\"main_hr\">
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['setting', load_lang('return')]]
            ))
    elif num == 3 or num == 4 or num == 7:
        if flask.request.method == 'POST':
            if num == 4:
                info_d = 'body'
                end_r = '4'
                coverage = ''
            elif num == 7:
                info_d = 'bottom_body'
                end_r = '7'
                coverage = ''
            else:
                info_d = 'head'
                end_r = '3'
                if flask.request.args.get('skin', '') == '':
                    coverage = ''
                else:
                    coverage = flask.request.args.get('skin', '')

            curs.execute(db_change("select name from other where name = ? and coverage = ?"), [info_d, coverage])
            if curs.fetchall():
                curs.execute(db_change("update other set data = ? where name = ? and coverage = ?"), [
                    flask.request.form.get('content', ''),
                    info_d,
                    coverage
                ])
            else:
                curs.execute(db_change("insert into other (name, data, coverage) values (?, ?, ?)"), [info_d, flask.request.form.get('content', ''), coverage])

            conn.commit()

            admin_check(None, 'edit_set')

            return redirect('/setting/' + end_r + '?skin=' + flask.request.args.get('skin', ''))
        else:
            if num == 4:
                curs.execute(db_change("select data from other where name = 'body'"))
                title = '_body'
                start = ''
                plus = '''
                    <button id="preview" type="button" onclick="load_raw_preview(\'content\', \'see_preview\')">''' + load_lang('preview') + '''</button>
                    <hr class=\"main_hr\">
                    <div id="see_preview"></div>
                '''
            elif num == 7:
                curs.execute(db_change("select data from other where name = 'bottom_body'"))
                title = '_bottom_body'
                start = ''
                plus = '''
                    <button id="preview" type="button" onclick="load_raw_preview(\'content\', \'see_preview\')">''' + load_lang('preview') + '''</button>
                    <hr class=\"main_hr\">
                    <div id="see_preview"></div>
                '''
            else:
                curs.execute(db_change("select data from other where name = 'head' and coverage = ?"), [flask.request.args.get('skin', '')])
                title = '_head'
                start = '' + \
                    '<a href="?">(' + load_lang('all') + ')</a> ' + \
                    ' '.join(['<a href="?skin=' + i + '">(' + i + ')</a>' for i in load_skin('', 1)]) + '''
                    <hr class=\"main_hr\">
                    <span>&lt;style&gt;CSS&lt;/style&gt;<br>&lt;script&gt;JS&lt;/script&gt;</span>
                    <hr class=\"main_hr\">
                '''
                plus = ''

            head = curs.fetchall()
            if head:
                data = head[0][0]
            else:
                data = ''

            if flask.request.args.get('skin', '') != '':
                sub_plus = ' (' + flask.request.args.get('skin', '') + ')'
            else:
                sub_plus = ''

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang(data = 'main' + title, safe = 1), wiki_set(), custom(), other2([' (HTML)' + sub_plus, 0])],
                data = '''
                    <form method="post">
                        ''' + start + '''
                        <textarea rows="25" name="content" id="content">''' + html.escape(data) + '''</textarea>
                        <hr class=\"main_hr\">
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        ''' + plus + '''
                    </form>
                ''',
                menu = [['setting', load_lang('return')]]
            ))
    elif num == 5:
        if flask.request.method == 'POST':
            curs.execute(db_change("select name from other where name = 'robot'"))
            if curs.fetchall():
                curs.execute(db_change("update other set data = ? where name = 'robot'"), [flask.request.form.get('content', '')])
            else:
                curs.execute(db_change("insert into other (name, data) values ('robot', ?)"), [flask.request.form.get('content', '')])

            conn.commit()

            fw = open('./robots.txt', 'w')
            fw.write(re.sub('\r\n', '\n', flask.request.form.get('content', '')))
            fw.close()

            admin_check(None, 'edit_set')

            return redirect('/setting/5')
        else:
            if not os.path.exists('robots.txt'):
                curs.execute(db_change('select data from other where name = "robot"'))
                robot_test = curs.fetchall()
                if robot_test:
                    fw_test = open('./robots.txt', 'w')
                    fw_test.write(re.sub('\r\n', '\n', robot_test[0][0]))
                    fw_test.close()
                else:
                    fw_test = open('./robots.txt', 'w')
                    fw_test.write('User-agent: *\nDisallow: /\nAllow: /$\nAllow: /w/')
                    fw_test.close()

                    curs.execute(db_change('insert into other (name, data) values ("robot", "User-agent: *\nDisallow: /\nAllow: /$\nAllow: /w/")'))

            curs.execute(db_change("select data from other where name = 'robot'"))
            robot = curs.fetchall()
            if robot:
                data = robot[0][0]
            else:
                data = ''

            f = open('./robots.txt', 'r')
            lines = f.readlines()
            f.close()

            if not data or data == '':
                data = ''.join(lines)

            return easy_minify(flask.render_template(skin_check(),
                imp = ['robots.txt', wiki_set(), custom(), other2([0, 0])],
                data = '''
                    <a href="/robots.txt">(''' + load_lang('view') + ''')</a>
                    <hr class=\"main_hr\">
                    <form method="post">
                        <textarea rows="25" name="content">''' + html.escape(data) + '''</textarea>
                        <hr class=\"main_hr\">
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['setting', load_lang('return')]]
            ))
    elif num == 6:
        i_list = [
            'recaptcha',
            'sec_re',
            'g_email',
            'g_pass'
        ]

        if flask.request.method == 'POST':
            for data in i_list:
                if data == 'g_email':
                    into_data = re.sub('@.*$', '', flask.request.form.get(data, ''))
                else:
                    into_data = flask.request.form.get(data, '')

                curs.execute(db_change("update other set data = ? where name = ?"), [into_data, data])

            conn.commit()

            admin_check(None, 'edit_set')

            return redirect('/setting/6')
        else:
            d_list = []

            x = 0

            for i in i_list:
                curs.execute(db_change('select data from other where name = ?'), [i])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list += [sql_d[0][0]]
                else:
                    curs.execute(db_change('insert into other (name, data) values (?, ?)'), [i, ''])

                    d_list += ['']

                x += 1

            conn.commit()

            return easy_minify(flask.render_template(skin_check(),
                imp = ['Google', wiki_set(), custom(), other2([0, 0])],
                data = '''
                    <form method="post">
                        <h2><a href="https://www.google.com/recaptcha/admin">''' + load_lang('recaptcha') + '''</a></h2>
                        <span>HTML</span>
                        <hr class=\"main_hr\">
                        <input name="recaptcha" placeholder='&lt;div class="g-recaptcha" data-sitekey="''' + load_lang('public_key') + '''"&gt;&lt;/div&gt;' value="''' + html.escape(d_list[0]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('secret_key') + '''</span>
                        <hr class=\"main_hr\">
                        <input name="sec_re" value="''' + html.escape(d_list[1]) + '''">
                        <hr class=\"main_hr\">
                        <h2><a href="https://support.google.com/mail/answer/7126229">''' + load_lang('google_imap') + '</a> (' + load_lang('restart_required') + ''')</h1>
                        <span>''' + load_lang('google_email') + '''</span>
                        <hr class=\"main_hr\">
                        <input name="g_email" value="''' + html.escape(d_list[2]) + '''">
                        <hr class=\"main_hr\">
                        <span><a href="https://security.google.com/settings/security/apppasswords">''' + load_lang('google_app_password') + '''</a></span>
                        <hr class=\"main_hr\">
                        <input type="password" name="g_pass" value="''' + html.escape(d_list[3]) + '''">
                        <hr class=\"main_hr\">
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['setting', load_lang('return')]]
            ))
    elif num == 8:
        i_list = {
            1 : 'edit',
            2 : 'discussion',
            3 : 'upload_acl',
            4 : 'all_view_acl',
            5 : 'edit_req_acl'
        }
        n_list = {
            1 : 'normal',
            2 : 'normal',
            3 : 'normal',
            4 : 'normal',
            5 : 'normal'
        }

        if flask.request.method == 'POST':
            for i in i_list:
                curs.execute(db_change("update other set data = ? where name = ?"), [
                    flask.request.form.get(i_list[i], n_list[i]),
                    i_list[i]
                ])

            conn.commit()

            admin_check(None, 'edit_set')

            return redirect('/setting/8')
        else:
            d_list = {}

            for i in i_list:
                curs.execute(db_change('select data from other where name = ?'), [i_list[i]])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list[i] = sql_d[0][0]
                else:
                    curs.execute(db_change('insert into other (name, data) values (?, ?)'), [i_list[i], n_list[i]])

                    d_list[i] = n_list[i]

            conn.commit()

            acl_div = ['', '', '', '', '']
            acl_list = ['normal', 'user', 'admin', 'owner', '50_edit', 'email']
            for i in range(0, 5):
                for acl_data in acl_list:
                    if acl_data == d_list[i + 1]:
                        acl_div[i] = '<option value="' + acl_data + '">' + acl_data + '</option>' + acl_div[i]
                    else:
                        acl_div[i] += '<option value="' + acl_data + '">' + acl_data + '</option>'

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('main_acl_setting'), wiki_set(), custom(), other2([0, 0])],
                data = '''
                    <form method="post">
                        <span>''' + load_lang('document_acl') + '</span> <a href="/acl/TEST">(' + load_lang('reference') + ''')</a>
                        <hr class=\"main_hr\">
                        <select name="edit">''' + acl_div[0] + '''</select>
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('discussion_acl') + '''</span>
                        <hr class=\"main_hr\">
                        <select name="discussion">''' + acl_div[1] + '''</select>
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('upload_acl') + '''</span>
                        <hr class=\"main_hr\">
                        <select name="upload_acl">''' + acl_div[2] + '''</select>
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('view_acl') + '''</span>
                        <hr class=\"main_hr\">
                        <select name="all_view_acl">''' + acl_div[3] + '''</select>
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('edit_req_acl') + '''</span>
                        <hr class=\"main_hr\">
                        <select name="edit_req_acl">''' + acl_div[4] + '''</select>
                        <hr class=\"main_hr\">
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['setting', load_lang('return')]]
            ))
    else:
        return redirect()
