from .tool.func import *

def main_func_setting(db_set, num = 0):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if not (num == 0 or num == 8) and admin_check() != 1:
            return re_error('/ban')

        if num == 0:
            li_list = [
                [1, load_lang('main_setting')],
                [2, load_lang('text_setting')],
                [5, 'robots.txt'],
                [6, load_lang('ext_api_req_set')],
                [3, load_lang('main_head')],
                [4, load_lang('main_body')],
                [7, load_lang('main_bottom_body')]
            ]

            li_data = ''.join(['<li><a href="/setting/' + str(li[0]) + '">' + li[1] + '</a></li>' for li in li_list])

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('setting'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '<h2>' + load_lang('list') + '</h2><ul class="inside_ul">' + li_data + '</ul>',
                menu = [['manager', load_lang('return')]]
            ))
        elif num == 1:
            i_list = {
                0 : 'name',
                2 : 'frontpage',
                4 : 'upload',
                5 : 'skin',
                7 : 'reg',
                8 : 'ip_view',
                9 : 'back_up',
                10 : 'port',
                11 : 'key',
                12 : 'update',
                15 : 'encode',
                16 : 'host',
                19 : 'slow_edit',
                20 : 'requires_approval',
                21 : 'backup_where',
                22 : 'domain',
                23 : 'ua_get'
            }
            n_list = {
                0 : 'Wiki',
                2 : 'FrontPage',
                4 : '2',
                5 : '',
                7 : '',
                8 : '',
                9 : '0',
                10 : '3000',
                11 : 'test',
                12 : 'stable',
                15 : 'sha3',
                16 : '0.0.0.0',
                19 : '0',
                20 : '',
                21 : '',
                22 : flask.request.host_url,
                23 : ''
            }

            if flask.request.method == 'POST':
                for i in i_list:
                    curs.execute(db_change("update other set data = ? where name = ?"), [
                        flask.request.form.get(i_list[i], n_list[i]),
                        i_list[i]
                    ])

                conn.commit()

                admin_check(None, 'edit_set (' + str(num) + ')')

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

                check_box_div = ['', '', '', '', '']
                for i in range(0, len(check_box_div)):
                    if i == 0:
                        acl_num = 7
                    elif i == 1:
                        acl_num = 8
                    elif i == 3:
                        acl_num = 20
                    else:
                        acl_num = 23

                    if d_list[acl_num]:
                        check_box_div[i] = 'checked="checked"'

                branch_div = ''
                branch_list = ['stable', 'dev', 'beta']
                for i in branch_list:
                    if d_list[12] == i:
                        branch_div = '<option value="' + i + '">' + i + '</option>' + branch_div
                    else:
                        branch_div += '<option value="' + i + '">' + i + '</option>'

                sqlite_only = 'style="display:none;"' if db_set != 'sqlite' else ''

                return easy_minify(flask.render_template(skin_check(),
                    imp = [load_lang('main_setting'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                    data = '''
                        <form method="post" id="main_set_data">
                            <h2>1. ''' + load_lang('basic_set') + '''</h2>
                            <span>''' + load_lang('wiki_name') + '''</span>
                            <hr class="main_hr">
                            <input name="name" value="''' + html.escape(d_list[0]) + '''">
                            <hr class="main_hr">

                            <span><a href="/setting/9">(''' + load_lang('wiki_logo') + ''')</a></span>
                            <hr class="main_hr">

                            <span>''' + load_lang('main_page') + '''</span>
                            <hr class="main_hr">
                            <input name="frontpage" value="''' + html.escape(d_list[2]) + '''">
                            <hr class="main_hr">

                            <span>''' + load_lang('domain') + '''</span> (EX : http://2du.pythonanywhere.com/)
                            <hr class="main_hr">
                            <input name="''' + i_list[22] + '''" value="''' + html.escape(d_list[22]) + '''">
                            <hr class="main_hr">

                            <span>''' + load_lang('wiki_host') + '''</span>
                            <hr class="main_hr">
                            <input name="host" value="''' + html.escape(d_list[16]) + '''">
                            <hr class="main_hr">

                            <span>''' + load_lang('wiki_port') + '''</span>
                            <hr class="main_hr">
                            <input name="port" value="''' + html.escape(d_list[10]) + '''">
                            <hr class="main_hr">

                            <span>''' + load_lang('wiki_secret_key') + '''</span>
                            <hr class="main_hr">
                            <input type="password" name="key" value="''' + html.escape(d_list[11]) + '''">
                            <hr class="main_hr">

                            <span>''' + load_lang('encryption_method') + '''</span>
                            <hr class="main_hr">
                            <select name="encode">''' + acl_div[0] + '''</select>

                            <h2>2. ''' + load_lang('design_set') + '''</h2>
                            <span>''' + load_lang('wiki_skin') + '''</span>
                            <hr class="main_hr">
                            <select name="skin">''' + load_skin(d_list[5] if d_list[5] != '' else 'tenshi') + '''</select>

                            <h2>3. ''' + load_lang('login_set') + '''</h2>
                            <input type="checkbox" name="reg" ''' + check_box_div[0] + '''> ''' + load_lang('no_register') + '''
                            <hr class="main_hr">

                            <input type="checkbox" name="ip_view" ''' + check_box_div[1] + '''> ''' + load_lang('hide_ip') + '''
                            <hr class="main_hr">

                            <input type="checkbox" name="requires_approval" ''' + check_box_div[3] + '''> ''' + load_lang('requires_approval') + '''
                            <hr class="main_hr">

                            <input type="checkbox" name="ua_get" ''' + check_box_div[4] + '''> ''' + load_lang('ua_get_off') + '''

                            <h2>4. ''' + load_lang('server_set') + '''</h2>
                            <span>''' + load_lang('max_file_size') + ''' (MB)</span>
                            <hr class="main_hr">
                            <input name="upload" value="''' + html.escape(d_list[4]) + '''">
                            <hr class="main_hr">

                            <span ''' + sqlite_only + '''>
                                <span>
                                    ''' + load_lang('backup_interval') + ' (' + load_lang('hour') + ') (' + load_lang('off') + ' : 0) ' + \
                                    '(' + load_lang('restart_required') + ''')</span>
                                <hr class="main_hr">
                                <input name="back_up" value="''' + html.escape(d_list[9]) + '''">
                                <hr class="main_hr">

                                <span>
                                    ''' + load_lang('backup_where') + ' (' + load_lang('empty') + ' : ' + load_lang('default') + ') ' + \
                                    '(' + load_lang('restart_required') + ''') (EX : ./data/backup.db)
                                </span>
                                <hr class="main_hr">
                                <input name="backup_where" value="''' + html.escape(d_list[21]) + '''">
                                <hr class="main_hr">
                            </span>

                            <span>''' + load_lang('update_branch') + '''</span>
                            <hr class="main_hr">
                            <select name="update">''' + branch_div + '''</select>

                            <h2>5. ''' + load_lang('edit_set') + '''</h2>
                            <span><a href="/setting/8">(''' + load_lang('main_acl_setting') + ''')</a></span>
                            <hr class="main_hr">

                            <span>''' + load_lang('slow_edit') + ' (' + load_lang('second') + ') (' + load_lang('off') + ''' : 0)</span>
                            <hr class="main_hr">
                            <input name="''' + i_list[19] + '''" value="''' + html.escape(d_list[19]) + '''">

                            <hr class="main_hr">
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        <script>simple_render('main_set_data');</script>
                    ''',
                    menu = [['setting', load_lang('return')]]
                ))
        elif num == 2:
            i_list = [
                'contract',
                'no_login_warning',
                'edit_bottom_text',
                'copyright_checkbox_text',
                'check_key_text',
                'email_title',
                'email_text',
                'email_insert_text',
                'password_search_text',
                'reset_user_text',
                'error_401',
                'error_404',
                'approval_question',
                'edit_help',
                'upload_help',
                'upload_default',
                'license',
                'topic_text'
            ]
            if flask.request.method == 'POST':
                for i in i_list:
                    curs.execute(db_change("update other set data = ? where name = ?"), [
                        flask.request.form.get(i, ''),
                        i
                    ])

                conn.commit()

                admin_check(None, 'edit_set (' + str(num) + ')')

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
                    imp = [load_lang('text_setting'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                    data = '''
                        <form method="post" id="main_set_data">
                            <h2>1. ''' + load_lang('register_text') + ''' (HTML)</h2>
                            <textarea rows="3" name="''' + i_list[0] + '''">''' + html.escape(d_list[0]) + '''</textarea>

                            <h2>2. ''' + load_lang('non_login_alert') + ''' (HTML)</h2>
                            <textarea rows="3" name="''' + i_list[1] + '''">''' + html.escape(d_list[1]) + '''</textarea>

                            <h2>3. ''' + load_lang('edit_bottom_text') + ''' (HTML)</h2>
                            <textarea rows="3" name="''' + i_list[2] + '''">''' + html.escape(d_list[2]) + '''</textarea>

                            <h2>4. ''' + load_lang('copyright_checkbox_text') + ''' (HTML)</h2>
                            <textarea rows="3" name="''' + i_list[3] + '''">''' + html.escape(d_list[3]) + '''</textarea>

                            <h2>5. ''' + load_lang('check_key_text') + ''' (HTML)</h2>
                            <textarea rows="3" name="''' + i_list[4] + '''">''' + html.escape(d_list[4]) + '''</textarea>

                            <h2>6. ''' + load_lang('email_title') + '''</h2>
                            <textarea rows="3" name="''' + i_list[5] + '''">''' + html.escape(d_list[5]) + '''</textarea>

                            <h2>7. ''' + load_lang('email_text') + '''</h2>
                            <textarea rows="3" name="''' + i_list[6] + '''">''' + html.escape(d_list[6]) + '''</textarea>

                            <h2>8. ''' + load_lang('email_insert_text') + '''</h2>
                            <textarea rows="3" name="''' + i_list[7] + '''">''' + html.escape(d_list[7]) + '''</textarea>

                            <h2>9. ''' + load_lang('password_search_text') + '''</h2>
                            <textarea rows="3" name="''' + i_list[8] + '''">''' + html.escape(d_list[8]) + '''</textarea>

                            <h2>10. ''' + load_lang('reset_user_text') + '''</h2>
                            <textarea rows="3" name="''' + i_list[9] + '''">''' + html.escape(d_list[9]) + '''</textarea>

                            <h2>11. ''' + load_lang('error_401') + '''</h2>
                            <textarea rows="3" name="''' + i_list[10] + '''">''' + html.escape(d_list[10]) + '''</textarea>

                            <h2>12. ''' + load_lang('error_404') + '''</h2>
                            <textarea rows="3" name="''' + i_list[11] + '''">''' + html.escape(d_list[11]) + '''</textarea>

                            <h2>13. ''' + load_lang('approval_question') + '''</h2>
                            <sup>(1)</sup>
                            <hr class="main_hr">
                            <textarea rows="3" name="''' + i_list[12] + '''">''' + html.escape(d_list[12]) + '''</textarea>

                            <h2>14. ''' + load_lang('edit_help') + '''</h2>
                            <textarea rows="3" name="''' + i_list[13] + '''">''' + html.escape(d_list[13]) + '''</textarea>

                            <h2>15. ''' + load_lang('upload_help') + ''' (HTML)</h2>
                            <textarea rows="3" name="''' + i_list[14] + '''">''' + html.escape(d_list[14]) + '''</textarea>

                            <h2>16. ''' + load_lang('upload_default') + '''</h2>
                            <textarea rows="3" name="''' + i_list[15] + '''">''' + html.escape(d_list[15]) + '''</textarea>

                            <h2>17. ''' + load_lang('bottom_text') + ''' (HTML)</h2>
                            <textarea rows="3" name="''' + i_list[16] + '''">''' + html.escape(d_list[16]) + '''</textarea>
                            
                            <h2>18. ''' + load_lang('topic_text') + '''</h2>
                            <textarea rows="3" name="''' + i_list[17] + '''">''' + html.escape(d_list[17]) + '''</textarea>

                            <hr class="main_hr">
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        <ul id="footnote_data">
                            <li><a href="#note_1" id="note_1_end">(1)</a> ''' + load_lang('approval_question_visible_only_when_approval_on') + '''</li>
                        </ul>
                        <script>simple_render('main_set_data');</script>
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

                admin_check(None, 'edit_set (' + str(num) + ')')

                return redirect('/setting/' + end_r + '?skin=' + flask.request.args.get('skin', ''))
            else:
                if num == 4:
                    curs.execute(db_change("select data from other where name = 'body'"))
                    title = '_body'
                    start = ''
                    plus = '''
                        <button id="preview" type="button" onclick="load_raw_preview(\'content\', \'see_preview\')">''' + load_lang('preview') + '''</button>
                        <hr class="main_hr">
                        <div id="see_preview"></div>
                    '''
                elif num == 7:
                    curs.execute(db_change("select data from other where name = 'bottom_body'"))
                    title = '_bottom_body'
                    start = ''
                    plus = '''
                        <button id="preview" type="button" onclick="load_raw_preview(\'content\', \'see_preview\')">''' + load_lang('preview') + '''</button>
                        <hr class="main_hr">
                        <div id="see_preview"></div>
                    '''
                else:
                    curs.execute(db_change("select data from other where name = 'head' and coverage = ?"), [flask.request.args.get('skin', '')])
                    title = '_head'
                    start = '' + \
                        '<a href="?">(' + load_lang('all') + ')</a> ' + \
                        ' '.join(['<a href="?skin=' + i + '">(' + i + ')</a>' for i in load_skin('', 1)]) + '''
                        <hr class="main_hr">
                        <span>&lt;style&gt;CSS&lt;/style&gt;<br>&lt;script&gt;JS&lt;/script&gt;</span>
                        <hr class="main_hr">
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
                    imp = [load_lang(data = 'main' + title, safe = 1), wiki_set(), wiki_custom(), wiki_css(['(HTML)' + sub_plus, 0])],
                    data = '''
                        <form method="post">
                            ''' + start + '''
                            <textarea rows="25" placeholder="''' + load_lang('enter_html') + '''" name="content" id="content">''' + html.escape(data) + '''</textarea>
                            <hr class="main_hr">
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

                fw = open('./robots.txt', 'w', encoding='utf8')
                fw.write(re.sub('\r\n', '\n', flask.request.form.get('content', '')))
                fw.close()

                admin_check(None, 'edit_set (' + str(num) + ')')

                return redirect('/setting/5')
            else:
                if not os.path.exists('robots.txt'):
                    curs.execute(db_change('select data from other where name = "robot"'))
                    robot_test = curs.fetchall()
                    if robot_test:
                        fw_test = open('./robots.txt', 'w', encoding='utf8')
                        fw_test.write(re.sub('\r\n', '\n', robot_test[0][0]))
                        fw_test.close()
                    else:
                        fw_test = open('./robots.txt', 'w', encoding='utf8')
                        fw_test.write('User-agent: *\nDisallow: /\nAllow: /$\nAllow: /w/')
                        fw_test.close()

                        curs.execute(db_change('insert into other (name, data) values ("robot", "User-agent: *\nDisallow: /\nAllow: /$\nAllow: /w/")'))

                curs.execute(db_change("select data from other where name = 'robot'"))
                robot = curs.fetchall()
                if robot:
                    data = robot[0][0]
                else:
                    data = ''

                f = open('./robots.txt', encoding='utf8')
                lines = f.readlines()
                f.close()

                if not data or data == '':
                    data = ''.join(lines)

                return easy_minify(flask.render_template(skin_check(),
                    imp = ['robots.txt', wiki_set(), wiki_custom(), wiki_css([0, 0])],
                    data = '''
                        <a href="/robots.txt">(''' + load_lang('view') + ''')</a>
                        <hr class="main_hr">
                        <form method="post">
                            <textarea rows="25" name="content">''' + html.escape(data) + '''</textarea>
                            <hr class="main_hr">
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                    ''',
                    menu = [['setting', load_lang('return')]]
                ))
        elif num == 6:
            i_list = [
                'recaptcha',
                'sec_re',
                'smtp_server',
                'smtp_port',
                'smtp_security',
                'smtp_email',
                'smtp_pass',
                'recaptcha_ver',
                'oauth_client_id',
                'email_have'
            ]

            if flask.request.method == 'POST':
                for data in i_list:
                    into_data = flask.request.form.get(data, '')

                    curs.execute(db_change("update other set data = ? where name = ?"), [into_data, data])

                conn.commit()

                admin_check(None, 'edit_set (' + str(num) + ')')

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

                security_radios = ''
                for i in ['tls', 'starttls', 'plain']:
                    security_radios += '<input name="smtp_security" type="radio" value="' + i + '" ' + ('checked' if d_list[4] == i else '') + '>' + i + '<hr class="main_hr">'

                re_ver = ''
                if d_list[7] == '':
                    re_ver += '<option value="">v2</option><option value="v3">v3</option>'
                else:
                    re_ver += '<option value="v3">v3</option><option value="">v2</option>'

                return easy_minify(flask.render_template(skin_check(),
                    imp = [load_lang('ext_api_req_set'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                    data = '''
                        <form method="post" id="main_set_data">
                            <h2>1. ''' + load_lang('recaptcha') + '''</h2>
                            <a href="https://www.google.com/recaptcha/admin">(Google)</a>
                            <hr class="main_hr">

                            <span>''' + load_lang('public_key') + '''</span>
                            <hr class="main_hr">
                            <input name="recaptcha" value="''' + html.escape(d_list[0]) + '''">
                            <hr class="main_hr">

                            <span>''' + load_lang('secret_key') + '''</span>
                            <hr class="main_hr">
                            <input name="sec_re" value="''' + html.escape(d_list[1]) + '''">
                            <hr class="main_hr">

                            <select name="recaptcha_ver">
                                ''' + re_ver + '''
                            </select>

                            <h2>2. ''' + load_lang('email_setting') + '''</h1>
                            <input type="checkbox" name="email_have" ''' + ('checked' if d_list[9] != '' else '')  + '''> ''' + \
                                 load_lang('email_required') + '''

                            <h2>2.1. ''' + load_lang('smtp_setting') + '''</h1>
                            <a href="https://support.google.com/mail/answer/7126229">(Google)</a>
                            <hr class="main_hr">

                            <span>''' + load_lang('smtp_server') + '''</span>
                            <hr class="main_hr">
                            <input name="smtp_server" value="''' + html.escape(d_list[2]) + '''">
                            <hr class="main_hr">

                            <span>''' + load_lang('smtp_port') + '''</span>
                            <hr class="main_hr">
                            <input name="smtp_port" value="''' + html.escape(d_list[3]) + '''">
                            <hr class="main_hr">

                            <span>''' + load_lang('smtp_security') + '''</span>
                            <hr class="main_hr">
                            ''' + security_radios + '''
                            <hr class="main_hr">

                            <span>''' + load_lang('smtp_username') + '''</span>
                            <hr class="main_hr">
                            <input name="smtp_email" value="''' + html.escape(d_list[5]) + '''">
                            <hr class="main_hr">

                            <span>''' + load_lang('smtp_password') + '''</span>
                            <hr class="main_hr">
                            <input type="password" name="smtp_pass" value="''' + html.escape(d_list[6]) + '''">

                            <h2>3. ''' + load_lang('oauth') + '''</h2>
                            <a href="https://developers.google.com/identity/protocols/oauth2">(Google)</a>
                            <hr class="main_hr">

                            <span>''' + load_lang('oauth_client_id') + '''</span>
                            <hr class="main_hr">
                            <input name="oauth_client_id" value="''' + html.escape(d_list[8]) + '''">
                            <hr class="main_hr">

                            <hr class="main_hr">
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        <script>simple_render('main_set_data');</script>
                    ''',
                    menu = [['setting', load_lang('return')]]
                ))
        elif num == 8:
            i_list = {
                1 : 'edit',
                2 : 'discussion',
                3 : 'upload_acl',
                4 : 'all_view_acl',
                5 : 'many_upload_acl',
                6 : 'vote_acl'
            }

            if flask.request.method == 'POST':
                if admin_check(None, 'edit_set (' + str(num) + ')') != 1:
                    return re_error('/ban')
                else:
                    for i in i_list:
                        curs.execute(db_change("update other set data = ? where name = ?"), [
                            flask.request.form.get(i_list[i], 'normal'),
                            i_list[i]
                        ])

                    conn.commit()

                    return redirect('/setting/8')
            else:
                d_list = {}

                if admin_check() != 1:
                    disable = 'disabled'
                else:
                    disable = ''

                for i in i_list:
                    curs.execute(db_change('select data from other where name = ?'), [i_list[i]])
                    sql_d = curs.fetchall()
                    if sql_d:
                        d_list[i] = sql_d[0][0]
                    else:
                        curs.execute(db_change('insert into other (name, data) values (?, ?)'), [i_list[i], 'normal'])

                        d_list[i] = 'normal'

                conn.commit()

                acl_div = []
                for i in range(0, len(i_list)):
                    acl_div += ['']

                acl_list = get_acl_list()
                for i in range(0, len(i_list)):
                    for data_list in acl_list:
                        if data_list == d_list[i + 1]:
                            check = 'selected="selected"'
                        else:
                            check = ''

                        acl_div[i] += '<option value="' + data_list + '" ' + check + '>' + (data_list if data_list != '' else 'normal') + '</option>'

                return easy_minify(flask.render_template(skin_check(),
                    imp = [load_lang('main_acl_setting'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                    data = '''
                        <form method="post">
                            <a href="/acl/TEST#exp">(''' + load_lang('reference') + ''')</a>
                            <hr class="main_hr">
                            <span>''' + load_lang('document_acl') + '''</span> 
                            <hr class="main_hr">
                            <select ''' + disable + ''' name="edit">''' + acl_div[0] + '''</select>
                            <hr class="main_hr">
                            <span>''' + load_lang('discussion_acl') + '''</span>
                            <hr class="main_hr">
                            <select ''' + disable + ''' name="discussion">''' + acl_div[1] + '''</select>
                            <hr class="main_hr">
                            <span>''' + load_lang('upload_acl') + '''</span>
                            <hr class="main_hr">
                            <select ''' + disable + ''' name="upload_acl">''' + acl_div[2] + '''</select>
                            <hr class="main_hr">
                            <span>''' + load_lang('view_acl') + '''</span>
                            <hr class="main_hr">
                            <select ''' + disable + ''' name="all_view_acl">''' + acl_div[3] + '''</select>
                            <hr class="main_hr">
                            <span>''' + load_lang('many_upload_acl') + '''</span>
                            <hr class="main_hr">
                            <select ''' + disable + ''' name="many_upload_acl">''' + acl_div[4] + '''</select>
                            <hr class="main_hr">
                            <span>''' + load_lang('vote_acl') + '''</span>
                            <hr class="main_hr">
                            <select ''' + disable + ''' name="vote_acl">''' + acl_div[5] + '''</select>
                            <hr class="main_hr">
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                    ''',
                    menu = [['setting', load_lang('return')]]
                ))
        elif num == 9:
            skin_list = [0] + load_skin('', 1)
            i_list = []
            for i in skin_list:
                i_list += [['logo', '' if i == 0 else i]]

            if flask.request.method == 'POST':
                for i in i_list:
                    curs.execute(db_change("update other set data = ? where name = ? and coverage = ?"), [
                        flask.request.form.get(('main_css' if i[1] == '' else i[1]), ''),
                        i[0], 
                        i[1]
                    ])

                conn.commit()

                admin_check(None, 'edit_set (' + str(num) + ')')

                return redirect('/setting/9')
            else:
                d_list = []
                for i in i_list:
                    curs.execute(db_change('select data from other where name = ? and coverage = ?'), [i[0], i[1]])
                    sql_d = curs.fetchall()
                    if sql_d:
                        d_list += [sql_d[0][0]]
                    else:
                        curs.execute(db_change('insert into other (name, data, coverage) values (?, ?, ?)'), [i[0], '', i[1]])

                        d_list += ['']
                        
                conn.commit()

                end_data = ''
                for i in range(0, len(skin_list)):
                    end_data += '' + \
                        '<span>' + load_lang('wiki_logo') + ' ' + ('(' + skin_list[i] + ')' if skin_list[i] != 0 else '') + ' (HTML)' + \
                        '<hr class="main_hr">' + \
                        '<input name="' + (skin_list[i] if skin_list[i] != 0 else 'main_css') + '" value="' + html.escape(d_list[i]) + '">' + \
                        '<hr class="main_hr">' + \
                    ''

                return easy_minify(flask.render_template(skin_check(),
                    imp = [load_lang('wiki_logo'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                    data = '''
                        <form method="post">
                            ''' + end_data + '''
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                    ''',
                    menu = [['setting', load_lang('return')]]
                ))
        else:
            return redirect()