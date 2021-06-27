from .tool.func import *

def user_setting_2(conn, server_init):
    curs = conn.cursor()

    support_language = ['default'] + server_init.server_init()['language']['list']
    ip = ip_check()

    if ban_check() == 1:
        return re_error('/ban')

    if ip_or_user(ip) == 0:
        if flask.request.method == 'POST':
            auto_list = [
                ['skin', flask.request.form.get('skin', '')], 
                ['lang', flask.request.form.get('lang', '')]
            ]
                
            twofa_turn_on = 0 
            twofa_on = flask.request.form.get('2fa', '')
            if twofa_on != '':
                twofa_turn_on = 1
                twofa_pw = flask.request.form.get('2fa_pw', '')
                if twofa_pw != '':
                    twofa_pw = pw_encode(twofa_pw)
                    curs.execute(db_change("select data from user_set where id = ? and name = 'encode'"), [ip])
                    twofa_encode = curs.fetchall()[0][0]
                    auto_list += [['2fa', 'on'], ['2fa_pw', twofa_pw], ['2fa_pw_encode', twofa_encode]]
                else:
                    auto_list += [['2fa', 'on']]
            else:
                auto_list += [['2fa', '']]
                
            for auto_data in auto_list:
                curs.execute(db_change('select data from user_set where name = ? and id = ?'), [auto_data[0], ip])
                if curs.fetchall():
                    curs.execute(db_change("update user_set set data = ? where name = ? and id = ?"), [auto_data[1], auto_data[0], ip])
                else:
                    curs.execute(db_change("insert into user_set (name, id, data) values (?, ?, ?)"), [auto_data[0], ip, auto_data[1]])

            conn.commit()

            return redirect('/change')
        else:
            curs.execute(db_change('select data from user_set where name = "email" and id = ?'), [ip])
            data = curs.fetchall()
            email = data[0][0] if data else '-'

            div2 = load_skin('', 0, 1)
            div3 = ''

            curs.execute(db_change('select data from user_set where name = "lang" and id = ?'), [ip_check()])
            data = curs.fetchall()
            data = [['default']] if not data else data

            for lang_data in support_language:
                see_data = lang_data if lang_data != 'default' else load_lang('default')
                
                if data and data[0][0] == lang_data:
                    div3 = '<option value="' + lang_data + '">' + see_data + '</option>' + div3
                else:
                    div3 += '<option value="' + lang_data + '">' + see_data + '</option>'

            curs.execute(db_change('select data from user_set where name = "2fa" and id = ?'), [ip])
            fa_data = curs.fetchall()
            fa_data = 'checked' if fa_data and fa_data[0][0] != '' else ''

            curs.execute(db_change('select data from user_set where name = "2fa_pw" and id = ?'), [ip])
            fa_data_pw = curs.fetchall()
            fa_data_pw = load_lang('2fa_password_change') if fa_data_pw else load_lang('2fa_password')

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('user_setting'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <div id="get_user_info"></div>
                        <script>load_user_info("''' + ip + '''");</script>
                        <hr class="main_hr">
                        <a href="/pw_change">(''' + load_lang('password_change') + ''')</a>
                        <hr class="main_hr">
                        <span>''' + load_lang('email') + ''' : ''' + email + '''</span> <a href="/email_change">(''' + load_lang('email_change') + ''')</a>
                        <h2>''' + load_lang('main') + '''</h2>
                        <span>''' + load_lang('skin') + '''</span>
                        <hr class="main_hr">
                        <select name="skin">''' + div2 + '''</select>
                        <hr class="main_hr">
                        <span>''' + load_lang('language') + '''</span>
                        <hr class="main_hr">
                        <select name="lang">''' + div3 + '''</select>
                        <h2>''' + load_lang('2fa') + '''</h2>
                        <input type="checkbox" id="twofa_check_input" onclick="do_twofa_check(0);" name="2fa" value="on" ''' + fa_data + '''> ''' + load_lang('on') + '''
                        <div id="fa_plus_content">
                            <hr class="main_hr">
                            <input type="password" name="2fa_pw" placeholder="''' + fa_data_pw + '''">
                        </div>
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('save') + '''</button>
                        ''' + http_warring() + '''
                        <script>do_twofa_check(1);</script>
                    </form>
                ''',
                menu = [['user', load_lang('return')]]
            ))
    else:
        if flask.request.method == 'POST':
            flask.session['skin'] = flask.request.form.get('skin', '')
            flask.session['lang'] = flask.request.form.get('lang', '')
            
            return redirect('/change')
        else:
            div2 = load_skin(('' if not 'skin' in flask.session else flask.session['skin']), 0, 1)
            div3 = ''

            data = [['default']] if not 'lang' in flask.session else [[flask.session['lang']]]

            for lang_data in support_language:
                see_data = lang_data if lang_data != 'default' else load_lang('default')
                
                if data and data[0][0] == lang_data:
                    div3 = '<option value="' + lang_data + '">' + see_data + '</option>' + div3
                else:
                    div3 += '<option value="' + lang_data + '">' + see_data + '</option>'
            
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('user_setting'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <div id="get_user_info"></div>
                        <script>load_user_info("''' + ip + '''");</script>
                        <hr class="main_hr">
                        <h2>''' + load_lang('main') + '''</h2>
                        <span>''' + load_lang('skin') + '''</span>
                        <hr class="main_hr">
                        <select name="skin">''' + div2 + '''</select>
                        <hr class="main_hr">
                        <span>''' + load_lang('language') + '''</span>
                        <hr class="main_hr">
                        <select name="lang">''' + div3 + '''</select>
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('save') + '''</button>
                        ''' + http_warring() + '''
                        <hr class="main_hr">
                        <span>''' + load_lang('user_head_warring') + '''</span>
                    </form>
                ''',
                menu = [['user', load_lang('return')]]
            ))