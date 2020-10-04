from .tool.func import *

def login_register_2(conn):
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
        if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        user_id = flask.request.form.get('id', '')
        user_pw = flask.request.form.get('pw', '')
        user_repeat = flask.request.form.get('pw2', '')
        if user_id == '' or user_pw == '':
            return re_error('/error/27')

        if user_pw != user_repeat:
            return re_error('/error/20')

        if re.search(r'(?:[^A-Za-zㄱ-힣0-9])', user_id):
            return re_error('/error/8')

        curs.execute(db_change('select html from html_filter where kind = "name"'))
        set_d = curs.fetchall()
        for i in set_d:
            check_r = re.compile(i[0], re.I)
            if check_r.search(user_id):
                return re_error('/error/8')

        if len(user_id) > 32:
            return re_error('/error/7')

        curs.execute(db_change("select id from user where id = ?"), [user_id])
        if curs.fetchall():
            return re_error('/error/6')
    
        curs.execute(db_change("select id from user_application where id = ?"), [user_id])
        if curs.fetchall():
            return re_error('/error/6')

        hashed = pw_encode(user_pw)
        ans_q = flask.request.form.get('approval_question_answer', '')

        curs.execute(db_change('select data from other where name = "requires_approval"'))
        requires_approval = curs.fetchall()
        requires_approval = requires_approval and requires_approval[0][0] == 'on'
        requires_approval = None if admin == 1 else requires_approval
        if requires_approval:
            curs.execute(db_change('select data from other where name = "approval_question"'))
            approval_question = curs.fetchall()
            approval_question = approval_question[0][0] if approval_question and approval_question[0][0] else ''
        else:
            approval_question = ''

        # c_id, c_pw, c_ans, c_que, c_key, c_type
        flask.session['c_id'] = user_id
        flask.session['c_pw'] = hashed
        flask.session['c_type'] = 'register'
        if requires_approval:
            flask.session['c_ans'] = flask.request.form.get('approval_question_answer', '')
            flask.session['c_que'] = approval_question
        
        curs.execute(db_change('select data from other where name = "email_have"'))
        sql_data = curs.fetchall()
        if sql_data and sql_data[0][0] != '' and admin != 1:
            flask.session['c_key'] = load_random_key(32)

            return redirect('/need_email')
        else:
            flask.session['c_key'] = 'email_pass'

            return redirect('/check_key')
    else:
        curs.execute(db_change('select data from other where name = "contract"'))
        data = curs.fetchall()
        contract = (data[0][0] + '<hr class="main_hr">') if data and data[0][0] != '' else ''

        http_warring = '<hr class="main_hr"><span>' + load_lang('http_warring') + '</span>'
        approval_question = ''
        
        curs.execute(db_change('select data from other where name = "requires_approval"'))
        requires_approval = curs.fetchall()
        requires_approval = requires_approval and requires_approval[0][0] == 'on'
        requires_approval = None if admin == 1 else requires_approval
        if requires_approval:
            curs.execute(db_change('select data from other where name = "approval_question"'))
            data = curs.fetchall()
            if data and data[0][0] != '':
                approval_question = '''
                    <hr class="main_hr">
                    <span>''' + load_lang('approval_question') + ' : ' + data[0][0] + '''<span>
                    <hr class="main_hr">
                    <input placeholder="''' + load_lang('approval_question') + '''" name="approval_question_answer" type="text">
                    <hr class="main_hr">
                '''

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('register'), wiki_set(), custom(), other2([0, 0])],
            data = '''
                <form method="post">
                    ''' + contract + '''
                    <input placeholder="''' + load_lang('id') + '''" name="id" type="text">
                    <hr class="main_hr">
                    <input placeholder="''' + load_lang('password') + '''" name="pw" type="password">
                    <hr class="main_hr">
                    <input placeholder="''' + load_lang('password_confirm') + '''" name="pw2" type="password">
                    <hr class="main_hr">
                    ''' + approval_question + '''
                    ''' + captcha_get() + '''
                    <button type="submit">''' + load_lang('save') + '''</button>
                    ''' + http_warring + '''
                </form>
            ''',
            menu = [['user', load_lang('return')]]
        ))
