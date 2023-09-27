from .tool.func import *

def login_register_2():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if ban_check(None, 'login') == 1:
            return re_error('/ban')

        ip = ip_check()
        admin = admin_check()
        if admin != 1 and ip_or_user(ip) == 0:
            return redirect('/user')

        if admin != 1:
            curs.execute(db_change('select data from other where name = "reg"'))
            set_d = curs.fetchall()
            if set_d and set_d[0][0] == 'on':
                return re_error('/ban')

        if flask.request.method == 'POST':
            # 리캡차
            if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return re_error('/error/13')
            else:
                captcha_post('', 0)

            user_id = flask.request.form.get('id', '')
            user_pw = flask.request.form.get('pw', '')
            user_repeat = flask.request.form.get('pw2', '')

            # PW 검증
            if user_id == '' or user_pw == '':
                return re_error('/error/27')

            if user_pw != user_repeat:
                return re_error('/error/20')

            # PW 길이 제한
            curs.execute(db_change("select data from other where name = 'password_min_length'"))
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '':
                password_min_length = int(number_check(db_data[0][0]))
                if password_min_length > len(user_pw):
                    return re_error('/error/40')

            if do_user_name_check(user_id) == 1:
                return re_error('/error/8')

            if admin != 1:
                # 이메일 필요시 /register/email로 발송
                curs.execute(db_change('select data from other where name = "email_have"'))
                sql_data = curs.fetchall()
                if sql_data and sql_data[0][0] != '':
                    # 임시로 세션에 저장
                    flask.session['reg_id'] = user_id
                    flask.session['reg_pw'] = user_pw

                    return redirect('/register/email')

                # 가입 승인 필요시 /register/submit으로 발송
                curs.execute(db_change('select data from other where name = "requires_approval"'))
                sql_data = curs.fetchall()
                if sql_data and sql_data[0][0] != '':
                    flask.session['submit_id'] = user_id
                    flask.session['submit_pw'] = user_pw

                    return redirect('/register/submit')

            # 전부 아니면 바로 가입 후 /login으로 발송
            add_user(user_id, user_pw)

            conn.commit()

            return redirect('/login')
        else:
            curs.execute(db_change('select data from other where name = "contract"'))
            data = curs.fetchall()
            contract = (data[0][0] + '<hr class="main_hr">') if data and data[0][0] != '' else ''

            curs.execute(db_change("select data from other where name = 'password_min_length'"))
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '':
                password_min_length = ' (' + load_lang('password_min_length') + ' : ' + db_data[0][0] + ')'
            else:
                password_min_length = ''

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('register'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        ''' + contract + '''

                        <input placeholder="''' + load_lang('id') + '''" name="id" type="text">
                        <hr class="main_hr">

                        <input placeholder="''' + load_lang('password') + password_min_length + '''" name="pw" type="password">
                        <hr class="main_hr">

                        <input placeholder="''' + load_lang('password_confirm') + '''" name="pw2" type="password">
                        <hr class="main_hr">

                        ''' + captcha_get() + '''

                        <button type="submit">''' + load_lang('save') + '''</button>

                        ''' + http_warning() + '''
                    </form>
                ''',
                menu = [['user', load_lang('return')]]
            ))