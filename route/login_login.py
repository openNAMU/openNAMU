from .tool.func import *

async def login_login():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if ip_or_user(ip) == 0:
            return redirect(conn, '/user')

        if (await ban_check(None, 'login'))[0] == 1:
            return await re_error(conn, 0)

        if flask.request.method == 'POST':
            if await captcha_post(conn, flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return await re_error(conn, 13)

            user_agent = flask.request.headers.get('User-Agent', '')
            user_id = flask.request.form.get('id', '')
            user_pw = flask.request.form.get('pw', '')

            curs.execute(db_change("select data from user_set where id = ? and name = 'pw'"), [user_id])
            db_data = curs.fetchall()
            if not db_data:
                return await re_error(conn, 2)
            else:
                db_user_pw = db_data[0][0]
                
            curs.execute(db_change("select data from user_set where id = ? and name = 'encode'"), [user_id])
            db_data = curs.fetchall()
            if not db_data:
                return await re_error(conn, 2)
            else:
                db_user_encode = db_data[0][0]

            if pw_check(conn, user_pw, db_user_pw, db_user_encode, user_id) != 1:
                return await re_error(conn, 10)

            curs.execute(db_change('select data from user_set where name = "2fa" and id = ?'), [user_id])
            fa_data = curs.fetchall()
            if fa_data and fa_data[0][0] != '':
                flask.session['login_id'] = user_id

                return redirect(conn, '/login/2fa')
            else:
                flask.session['id'] = user_id

                ua_plus(conn, user_id, ip, user_agent, get_time())

                return redirect(conn, '/user')
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'login'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data =  '''
                        <form method="post">
                            <input placeholder="''' + get_lang(conn, 'id') + '''" name="id" type="text">
                            <hr class="main_hr">
                            <input placeholder="''' + get_lang(conn, 'password') + '''" name="pw" type="password">
                            <hr class="main_hr">
                            <!-- <label><input type="checkbox" name="auto_login"> ''' + get_lang(conn, 'auto_login') + ''' (''' + get_lang(conn, 'not_working') + ''')</label>
                            <hr class="main_hr"> -->
                            ''' + await captcha_get(conn) + '''
                            <button type="submit">''' + get_lang(conn, 'login') + '''</button>
                            ''' + http_warning(conn) + '''
                        </form>
                        ''',
                menu = [['user', get_lang(conn, 'return')]]
            ))