from .tool.func import *

def user_setting_2(conn, server_init):
    curs = conn.cursor()

    support_language = ['default'] + server_init.server_set_var['language']['list']
    ip = ip_check()

    if ban_check() == 1:
        return re_error('/ban')

    if ip_or_user(ip) == 0:
        if flask.request.method == 'POST':
            pass_list = ['2fa']
            auto_list = ['skin', 'lang'] + pass_list + ['2fa_pw', '2fa_pw_encode']

            for auto_data in auto_list:
                if auto_data == '2fa_pw':
                    if flask.request.form.get('2fa_pw', '') != '':
                        get_data = pw_encode(flask.request.form.get(auto_data, ''))
                    else:
                        get_data = ''
                elif auto_data == '2fa_pw_encode':
                    if flask.request.form.get('2fa_pw', '') != '':
                        curs.execute(db_change("select encode from user where id = ?"), [ip])
                        get_data = curs.fetchall()[0][0]
                    else:
                        get_data = ''
                else:
                    get_data = flask.request.form.get(auto_data, '')

                if auto_data in pass_list or get_data != '':
                    curs.execute(db_change('select data from user_set where name = ? and id = ?'), [auto_data, ip])
                    if curs.fetchall():
                        curs.execute(db_change("update user_set set data = ? where name = ? and id = ?"), [get_data, auto_data, ip])
                    else:
                        curs.execute(db_change("insert into user_set (name, id, data) values (?, ?, ?)"), [auto_data, ip, get_data])

            conn.commit()

            return redirect('/change')
        else:
            curs.execute(db_change('select data from user_set where name = "email" and id = ?'), [ip])
            data = curs.fetchall()
            if data:
                email = data[0][0]
            else:
                email = '-'

            div2 = load_skin('', 0, 1)
            div3 = ''

            curs.execute(db_change('select data from user_set where name = "lang" and id = ?'), [ip_check()])
            data = curs.fetchall()
            if not data:
                data = [['default']]

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
            
            http_warring = '' + \
                '<hr class="main_hr">' + \
                '<span>' + load_lang('http_warring') + '</span>' + \
            ''

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('user_setting'), wiki_set(), custom(), other2([0, 0])],
                data = '''
                    <form method="post">
                        <span>''' + load_lang('id') + ''' : ''' + ip_pas(ip) + '''</span>
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
                        <input type="checkbox" name="2fa" value="on" ''' + fa_data + '''> ''' + load_lang('on') + '''
                        <hr class="main_hr">
                        <input type="password" name="2fa_pw" placeholder="''' + fa_data_pw + '''">
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('save') + '''</button>
                        ''' + http_warring + '''
                    </form>
                ''',
                menu = [['user', load_lang('return')]]
            ))
    else:
        return redirect('/login')