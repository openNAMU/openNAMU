from .tool.func import *

async def give_user_fix(user_name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change("select data from user_set where id = ? and name = 'pw'"), [user_name])
        if not curs.fetchall():
            return await re_error(conn, 2)

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 3)

        if flask.request.method == 'POST':
            select = flask.request.form.get('select', '')

            await acl_check(tool = 'owner_auth', memo = 'user_fix (' + user_name + ') (' + select + ')')
            if select == 'password_change':
                password = flask.request.form.get('new_password', '')
                check_password = flask.request.form.get('password_check', '')

                if password == check_password:
                    hashed = pw_encode(conn, password)
                    curs.execute(db_change("update user_set set data = ? where id = ? and name = 'pw'"), [
                        hashed,
                        user_name
                    ])
                else:
                    return await re_error(conn, 20)
            elif select == '2fa_password_change':
                password = flask.request.form.get('new_password', '')
                check_password = flask.request.form.get('password_check', '')

                if password == check_password:
                    hashed = pw_encode(conn, password)
                    curs.execute(db_change('select data from user_set where name = "2fa_pw" and id = ?'), [user_name])
                    if curs.fetchall():
                        curs.execute(db_change("update user_set set data = ? where name = '2fa_pw' and id = ?"), [hashed, user_name])
                    else:
                        curs.execute(db_change("insert into user_set (name, id, data) values ('2fa_pw', ?, ?)"), [user_name, hashed])
                else:
                    return await re_error(conn, 20)
            elif select == '2fa_off':
                curs.execute(db_change('select data from user_set where name = "2fa" and id = ?'), [user_name])
                if curs.fetchall():
                    curs.execute(db_change("update user_set set data = '' where name = '2fa' and id = ?"), [user_name])

            return redirect(conn, '/user/' + url_pas(user_name))
        else:
            return await render_template(
                await get_lang('user_fix'),
                '''
                    <form method="post">
                        <div id="opennamu_get_user_info">''' + html.escape(user_name) + '''</div>
                        <hr class="main_hr">
                        <a href="/change/user_name/''' + url_pas(user_name) + '''">(''' + await get_lang('change_user_name') + ''')</a>
                        <hr class="main_hr">
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="select">
                            <option value="password_change">''' + await get_lang('password_change') + '''</option>
                            <option value="2fa_password_change">''' + await get_lang('2fa_password_change') + '''</option>
                            <option value="2fa_off">''' + await get_lang('2fa_off') + '''</option>
                        </select></span>
                        <hr class="main_hr">
                        ''' + await get_lang('password_change') + ''' | ''' + await get_lang('2fa_password_change') + '''
                        <hr class="main_hr">
                        <input class="__ON_INPUT__" placeholder="''' + await get_lang('new_password') + '''" name="new_password" type="password">
                        <hr class="main_hr">
                        <input class="__ON_INPUT__" placeholder="''' + await get_lang('password_confirm') + '''" name="password_check" type="password">
                        <hr class="main_hr">
                        <button class="__ON_BUTTON__" type="submit">''' + await get_lang('go') + '''</button>
                    </form>
                ''',
                0,
                [['manager', await get_lang('return')]]
            )
