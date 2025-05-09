from .tool.func import *

async def login_login_2fa_email():
    with get_db_connect() as conn:
        curs = conn.cursor()

        # email 2fa
        # pw 2fa
        # q_a 2fa
        if not (flask.session and 'login_id' in flask.session):
            return redirect(conn, '/user')

        ip = ip_check()
        if ip_or_user(ip) == 0:
            return redirect(conn, '/user')

        if (await ban_check(None, 'login'))[0] == 1:
            return await re_error(conn, 0)

        if flask.request.method == 'POST':
            if await captcha_post(conn, flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return await re_error(conn, 13)

            user_agent = flask.request.headers.get('User-Agent', '')
            user_id = flask.session['b_id']
            user_pw = flask.request.form.get('pw', '')

            curs.execute(db_change('select data from user_set where name = "2fa_pw" and id = ?'), [user_id])
            user_1 = curs.fetchall()
            if user_1:
                curs.execute(db_change('select data from user_set where name = "2fa_pw_encode" and id = ?'), [user_id])
                user_1 = user_1[0][0]
                user_2 = curs.fetchall()[0][0]

                pw_check_d = pw_check(conn, user_pw, user_1, user_2, user_id)
                if pw_check_d != 1:
                    return await re_error(conn, 10)

            flask.session['id'] = user_id

            ua_plus(conn, user_id, ip, user_agent, get_time())

            flask.session.pop('b_id', None)

            return redirect(conn, '/user')
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'login'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data =  '''
                        <form method="post">
                            <input placeholder="''' + get_lang(conn, '2fa_password') + '''" name="pw" type="password">
                            <hr class=\"main_hr\">
                            ''' + await captcha_get(conn) + '''
                            <button type="submit">''' + get_lang(conn, 'login') + '''</button>
                            ''' + http_warning(conn) + '''
                        </form>
                        ''',
                menu = [['user', get_lang(conn, 'return')]]
            ))