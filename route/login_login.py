from .tool.func import *

def login_login_2():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if ip_or_user(ip) == 0:
            return redirect('/user')

        if ban_check(None, 'login') == 1:
            return re_error('/ban')

        if flask.request.method == 'POST':
            if 'login_count' in flask.session:
                count = int(number_check(flask.session['login_count']))
                if count > 3:
                    if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                        return re_error('/error/13')
                    else:
                        captcha_post('', 0)

            user_agent = flask.request.headers.get('User-Agent', '')
            user_id = flask.request.form.get('id', '')
            user_data = {}

            curs.execute(db_change(
                'select name, data from user_set where id = ? and name = "pw" or name = "encode"'
            ), [user_id])
            sql_data = curs.fetchall()
            if not sql_data:
                return re_error('/error/2')

            for i in sql_data:
                user_data[i[0]] = i[1]

            if len(user_data) < 2:
                return re_error('/error/2')

            if pw_check(
                flask.request.form.get('pw', ''),
                user_data['pw'],
                user_data['encode'],
                user_id
            ) != 1:
                if not 'login_count' in flask.session:
                    flask.session['login_count'] = 1
                else:
                    flask.session['login_count'] = int(number_check(flask.session['login_count'])) + 1

                return re_error('/error/10')

            curs.execute(db_change('select data from user_set where name = "2fa" and id = ?'), [user_id])
            fa_data = curs.fetchall()
            if fa_data and fa_data[0][0] != '':
                flask.session['login_id'] = user_id

                return redirect('/login/2fa')
            else:
                flask.session['id'] = user_id

                ua_plus(user_id, ip, user_agent, get_time())
                conn.commit()

                return redirect('/user')
        else:
            captcha_data = ''
            if 'login_count' in flask.session:
                count = int(number_check(flask.session['login_count']))
                if count > 3:
                    captcha_data = captcha_get()

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('login'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data =  '''
                        <form method="post">
                            <input placeholder="''' + load_lang('id') + '''" name="id" type="text">
                            <hr class="main_hr">
                            <input placeholder="''' + load_lang('password') + '''" name="pw" type="password">
                            <hr class="main_hr">
                            <!-- <input type="checkbox" name="auto_login"> ''' + load_lang('auto_login') + ''' (''' + load_lang('not_working') + ''')
                            <hr class="main_hr"> -->
                            ''' + captcha_data + '''
                            <button type="submit">''' + load_lang('login') + '''</button>
                            ''' + http_warning() + '''
                        </form>
                        ''',
                menu = [['user', load_lang('return')]]
            ))