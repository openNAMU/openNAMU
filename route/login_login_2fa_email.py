from .tool.func import *

def login_login_2fa_email_2():
    with get_db_connect() as conn:
        curs = conn.cursor()

        # email 2fa
        # pw 2fa
        # q_a 2fa
        if not (flask.session and 'login_id' in flask.session):
            return redirect('/user')

        ip = ip_check()
        if ip_or_user(ip) == 0:
            return redirect('/user')

        if ban_check(None, 'login') == 1:
            return re_error('/ban')

        if flask.request.method == 'POST':
            if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return re_error('/error/13')
            else:
                captcha_post('', 0)

            user_agent = flask.request.headers.get('User-Agent', '')
            user_id = flask.session['b_id']
            user_pw = flask.request.form.get('pw', '')

            curs.execute(db_change('select data from user_set where name = "2fa_pw" and id = ?'), [user_id])
            user_1 = curs.fetchall()
            if user_1:
                curs.execute(db_change('select data from user_set where name = "2fa_pw_encode" and id = ?'), [user_id])
                user_1 = user_1[0][0]
                user_2 = curs.fetchall()[0][0]

                pw_check_d = pw_check(user_pw, user_1, user_2, user_id)
                if pw_check_d != 1:
                    return re_error('/error/10')

            flask.session['id'] = user_id

            ua_plus(user_id, ip, user_agent, get_time())
            conn.commit()

            flask.session.pop('b_id', None)

            return redirect('/user')
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('login'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data =  '''
                        <form method="post">
                            <input placeholder="''' + load_lang('2fa_password') + '''" name="pw" type="password">
                            <hr class=\"main_hr\">
                            ''' + captcha_get() + '''
                            <button type="submit">''' + load_lang('login') + '''</button>
                            ''' + http_warning() + '''
                        </form>
                        ''',
                menu = [['user', load_lang('return')]]
            ))