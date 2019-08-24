from .tool.func import *
import pymysql

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
            'Google'
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
            6 : 'edit', 
            7 : 'reg', 
            8 : 'ip_view', 
            9 : 'back_up', 
            10 : 'port', 
            11 : 'key', 
            12 : 'update', 
            13 : 'email_have', 
            14 : 'discussion', 
            15 : 'encode', 
            16 : 'host'
        }
        n_list = {
            0 : 'Wiki', 
            1 : '', 
            2 : 'FrontPage', 
            3 : 'CC 0', 
            4 : '2', 
            5 : '', 
            6 : 'normal', 
            7 : '', 
            8 : '', 
            9 : '0', 
            10 : '3000', 
            11 : 'test', 
            12 : 'stable', 
            13 : '', 
            14 : 'normal', 
            15 : 'sha3', 
            16 : '0.0.0.0'
        }
        
        if flask.request.method == 'POST':
            for i in i_list:
                curs.execute("update other set data = %s where name = %s", [
                    flask.request.form.get(i_list[i], n_list[i]), 
                    i_list[i]]
                )

            conn.commit()

            admin_check(None, 'edit_set')

            return redirect('/setting/1')
        else:
            d_list = []
            
            for i in i_list:
                curs.execute('select data from other where name = %s', [i_list[i]])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list += [sql_d[0][0]]
                else:
                    curs.execute('insert into other (name, data) values (%s, %s)', [i_list[i], n_list[i]])
                    
                    d_list += [n_list[i]]

            conn.commit()
            
            div = ''
            acl_list = [
                [load_lang('member'), 'login'], 
                [load_lang('ip'), 'normal'], 
                [load_lang('admin'), 'admin']
            ]
            for i in acl_list:
                if i[1] == d_list[6]:
                    div = '<option value="' + i[1] + '">' + i[0] + '</option>' + div
                else:
                    div += '<option value="' + i[1] + '">' + i[0] + '</option>'

            div4 = ''
            for i in acl_list:
                if i[1] == d_list[14]:
                    div4 = '<option value="' + i[1] + '">' + i[0] + '</option>' + div4
                else:
                    div4 += '<option value="' + i[1] + '">' + i[0] + '</option>'

            ch_1 = ''
            if d_list[7]:
                ch_1 = 'checked="checked"'

            ch_2 = ''
            if d_list[8]:
                ch_2 = 'checked="checked"'
            
            ch_3 = ''
            if d_list[13]:
                ch_3 = 'checked="checked"'

            div2 = load_skin(d_list[5])

            div3 =''
            if d_list[12] == 'stable':
                div3 += '<option value="stable">stable</option>'
                div3 += '<option value="master">master</option>'
            else:
                div3 += '<option value="master">master</option>'
                div3 += '<option value="stable">stable</option>'
                
            div5 =''
            encode_data = ['sha256', 'sha3']
            for i in encode_data:
                if d_list[15] == i:
                    div5 = '<option value="' + i + '">' + i + '</option>' + div5
                else:
                    div5 += '<option value="' + i + '">' + i + '</option>'

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
                        <span>''' + load_lang('max_file_size') + ''' [MB]</span>
                        <hr class=\"main_hr\">
                        <input type="text" name="upload" value="''' + html.escape(d_list[4]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('backup_interval') + ' [' + load_lang('hour') + '''] (off : 0) {restart}</span>
                        <hr class=\"main_hr\">
                        <input type="text" name="back_up" value="''' + html.escape(d_list[9]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('wiki_skin') + '''</span>
                        <hr class=\"main_hr\">
                        <select name="skin">''' + div2 + '''</select>
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('default_acl') + '''</span>
                        <hr class=\"main_hr\">
                        <select name="edit">''' + div + '''</select>
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('default_discussion_acl') + '''</span>
                        <hr class=\"main_hr\">
                        <select name="discussion">''' + div4 + '''</select>
                        <hr class=\"main_hr\">
                        <input type="checkbox" name="reg" ''' + ch_1 + '''> ''' + load_lang('no_register') + '''
                        <hr class=\"main_hr\">
                        <input type="checkbox" name="ip_view" ''' + ch_2 + '''> ''' + load_lang('hide_ip') + '''
                        <hr class=\"main_hr\">
                        <input type="checkbox" name="email_have" ''' + ch_3 + '''> ''' + load_lang('email_required') + ''' {<a href="/setting/5">''' + load_lang('google_imap_required') + '''</a>}
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
                        <select name="update">''' + div3 + '''</select>
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('encryption_method') + '''</span>
                        <hr class=\"main_hr\">
                        <select name="encode">''' + div5 + '''</select>
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
            'error_404'
        ]
        if flask.request.method == 'POST':
            for i in i_list:
                curs.execute("update other set data = %s where name = %s", [
                    flask.request.form.get(i, ''), 
                    i
                ])

            conn.commit()
            
            admin_check(None, 'edit_set')

            return redirect('/setting/2')
        else:
            d_list = []
            
            for i in i_list:
                curs.execute('select data from other where name = %s', [i])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list += [sql_d[0][0]]
                else:
                    curs.execute('insert into other (name, data) values (%s, %s)', [i, ''])
                    
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
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['setting', load_lang('return')]]
            ))
    elif num == 3 or num == 4:
        if flask.request.method == 'POST':
            if num == 4:
                info_d = 'body'
                end_r = '4'
                coverage = ''
            else:
                info_d = 'head'
                end_r = '3'
                if flask.request.args.get('skin', '') == '':
                    coverage = ''
                else:
                    coverage = flask.request.args.get('skin', '')
                
            curs.execute("select name from other where name = %s and coverage = %s", [info_d, coverage])
            if curs.fetchall():
                curs.execute("update other set data = %s where name = %s and coverage = %s", [
                    flask.request.form.get('content', ''),
                    info_d,
                    coverage
                ])
            else:
                curs.execute("insert into other (name, data, coverage) values (%s, %s, %s)", [info_d, flask.request.form.get('content', ''), coverage])
            
            conn.commit()

            admin_check(None, 'edit_set')

            return redirect('/setting/' + end_r + '?skin=' + flask.request.args.get('skin', ''))
        else:
            if num == 4:
                curs.execute("select data from other where name = 'body'")
                title = '_body'
                start = ''
            else:
                curs.execute("select data from other where name = 'head' and coverage = %s", [flask.request.args.get('skin', '')])
                title = '_head'
                start = '<a href="?">(' + load_lang('all') + ')</a> ' + \
                ' '.join(['<a href="?skin=' + i + '">(' + i + ')</a>' for i in load_skin('', 1)]) + \
                '''
                    <hr class=\"main_hr\">
                    <span>&lt;style&gt;CSS&lt;/style&gt;<br>&lt;script&gt;JS&lt;/script&gt;</span>
                    <hr class=\"main_hr\">
                '''
                
            head = curs.fetchall()
            if head:
                data = head[0][0]
            else:
                data = ''

            return easy_minify(flask.render_template(skin_check(), 
                imp = [load_lang(data = 'main' + title, safe = 1), wiki_set(), custom(), other2([0, 0])],
                data = '''
                    <form method="post">
                        ''' + start + '''
                        <textarea rows="25" name="content">''' + html.escape(data) + '''</textarea>
                        <hr class=\"main_hr\">
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['setting', load_lang('return')]]
            ))
    elif num == 5:
        if flask.request.method == 'POST':
            curs.execute("select name from other where name = 'robot'")
            if curs.fetchall():
                curs.execute("update other set data = %s where name = 'robot'", [flask.request.form.get('content', '')])
            else:
                curs.execute("insert into other (name, data) values ('robot', %s)", [flask.request.form.get('content', '')])
            
            conn.commit()
            
            fw = open('./robots.txt', 'w')
            fw.write(re.sub('\r\n', '\n', flask.request.form.get('content', '')))
            fw.close()
            
            admin_check(None, 'edit_set')

            return redirect('/setting/4')
        else:
            curs.execute("select data from other where name = 'robot'")
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
                    <a href="/robots.txt">(view)</a>
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

                curs.execute("update other set data = %s where name = %s", [into_data, data])

            conn.commit()
            
            admin_check(None, 'edit_set')

            return redirect('/setting/6')
        else:
            d_list = []
            
            x = 0
            
            for i in i_list:
                curs.execute('select data from other where name = %s', [i])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list += [sql_d[0][0]]
                else:
                    curs.execute('insert into other (name, data) values (%s, %s)', [i, ''])
                    
                    d_list += ['']

                x += 1

            conn.commit()

            return easy_minify(flask.render_template(skin_check(), 
                imp = ['Google', wiki_set(), custom(), other2([0, 0])],
                data = '''
                    <form method="post">
                        <h2><a href="https://www.google.com/recaptcha/admin">recaptcha</a></h2>
                        <span>''' + load_lang('recaptcha') + ''' (HTML)</span>
                        <hr class=\"main_hr\">
                        <input name="recaptcha" value="''' + html.escape(d_list[0]) + '''">
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('recaptcha') + ' (' + load_lang('secret_key') + ''')</span>
                        <hr class=\"main_hr\">
                        <input name="sec_re" value="''' + html.escape(d_list[1]) + '''">
                        <hr class=\"main_hr\">
                        <h2><a href="https://support.google.com/mail/answer/7126229">''' + load_lang('google_imap') + '</a> {' + load_lang('restart_required') + '''}</h1>
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
    else:
        return redirect()