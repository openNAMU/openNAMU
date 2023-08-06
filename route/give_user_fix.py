from .tool.func import *

def give_user_fix(user_name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change("select data from user_set where id = ? and name = 'pw'"), [user_name])
        if not curs.fetchall():
            return re_error('/error/2')

        if admin_check() != 1:
            return re_error('/error/3')

        if flask.request.method == 'POST':
            select = flask.request.form.get('select', '')

            admin_check(None, 'user_fix (' + user_name + ') (' + select + ')')
            if select == 'password_change':
                password = flask.request.form.get('new_password', '')
                check_password = flask.request.form.get('password_check', '')

                if password == check_password:
                    hashed = pw_encode(password)
                    curs.execute(db_change("update user_set set data = ? where id = ? and name = 'pw'"), [
                        hashed,
                        user_name
                    ])
                else:
                    return re_error('/error/20')
            elif select == '2fa_password_change':
                password = flask.request.form.get('new_password', '')
                check_password = flask.request.form.get('password_check', '')

                if password == check_password:
                    hashed = pw_encode(password)
                    curs.execute(db_change('select data from user_set where name = "2fa_pw" and id = ?'), [user_name])
                    if curs.fetchall():
                        curs.execute(db_change("update user_set set data = ? where name = '2fa_pw' and id = ?"), [hashed, user_name])
                    else:
                        curs.execute(db_change("insert into user_set (name, id, data) values ('2fa_pw', ?, ?)"), [user_name, hashed])
                else:
                    return re_error('/error/20')
            elif select == '2fa_off':
                curs.execute(db_change('select data from user_set where name = "2fa" and id = ?'), [user_name])
                if curs.fetchall():
                    curs.execute(db_change("update user_set set data = '' where name = '2fa' and id = ?"), [user_name])

            conn.commit()

            return redirect('/user/' + url_pas(user_name))
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('user_fix'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <div id="opennamu_get_user_info">''' + html.escape(user_name) + '''</div>
                        <hr class="main_hr">
                        <a href="/change/user_name/''' + url_pas(user_name) + '''">(''' + load_lang('change_user_name') + ''')</a>
                        <hr class="main_hr">
                        <select name="select">
                            <option value="password_change">''' + load_lang('password_change') + '''</option>
                            <option value="2fa_password_change">''' + load_lang('2fa_password_change') + '''</option>
                            <option value="2fa_off">''' + load_lang('2fa_off') + '''</option>
                        </select>
                        <hr class="main_hr">
                        ''' + load_lang('password_change') + ''' | ''' + load_lang('2fa_password_change') + '''
                        <hr class="main_hr">
                        <input placeholder="''' + load_lang('new_password') + '''" name="new_password" type="password">
                        <hr class="main_hr">
                        <input placeholder="''' + load_lang('password_confirm') + '''" name="password_check" type="password">
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('go') + '''</button>
                    </form>
                ''',
                menu = [['manager', load_lang('return')]]
            ))