from .tool.func import *

def login_check_key_2(conn, tool):
    curs = conn.cursor()

    # 난잡한 코드 정리 필요
    if  flask.request.method == 'POST' or \
        ('c_key' in flask.session and flask.session['c_key'] == 'email_pass'):
        re_set_list = ['c_id', 'c_pw', 'c_ans', 'c_que', 'c_key', 'c_type', 'c_email']
        ip = ip_check()
        input_key = flask.request.form.get('key', '')
        user_agent = flask.request.headers.get('User-Agent', '')


        if  'c_type' in flask.session and \
            flask.session['c_type'] == 'pass_find' and \
            flask.session['c_key'] == input_key:
            curs.execute(db_change("update user set pw = ? where id = ?"), [pw_encode(flask.session['c_key']), flask.session['c_id']])
            conn.commit()

            user_id = flask.session['c_id']
            user_pw = flask.session['c_key']

            for i in re_set_list:
                flask.session.pop(i, None)

            curs.execute(db_change('select data from other where name = "reset_user_text"'))
            sql_d = curs.fetchall()
            b_text = (sql_d[0][0] + '<hr class="main_hr">') if sql_d and sql_d[0][0] != '' else ''

            curs.execute(db_change('select data from user_set where name = "2fa" and id = ?'), [user_id])
            if curs.fetchall():
                curs.execute(db_change("update user_set set data = '' where name = '2fa' and id = ?"), [user_id])

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('reset_user_ok'), wiki_set(), custom(), other2([0, 0])],
                data = b_text + load_lang('id') + ' : ' + user_id + '<br>' + load_lang('password') + ' : ' + user_pw,
                menu = [['user', load_lang('return')]]
            ))
        elif    'c_type' in flask.session and \
                (flask.session['c_key'] == input_key or flask.session['c_key'] == 'email_pass'):
            curs.execute(db_change('select data from other where name = "encode"'))
            db_data = curs.fetchall()

            if flask.session['c_type'] == 'register':
                if flask.session['c_key'] == 'email_pass':
                    flask.session['c_email'] = ''

                curs.execute(db_change("select id from user limit 1"))
                first = 1 if not curs.fetchall() else 0

                curs.execute(db_change("select id from user where id = ?"), [flask.session['c_id']])
                if curs.fetchall():
                    for i in re_set_lire:
                        flask.session.pop(i, None)

                    return re_error('/error/6')
            
                curs.execute(db_change("select id from user_application where id = ?"), [flask.session['c_id']])
                if curs.fetchall():
                    for i in re_set_lire:
                        flask.session.pop(i, None)

                    return re_error('/error/6')

                curs.execute(db_change('select data from other where name = "requires_approval"'))
                requires_approval = curs.fetchall()
                if requires_approval and requires_approval[0][0] == 'on':
                    application_token = load_random_key(32)
                    curs.execute(db_change(
                        "insert into user_application (id, pw, date, encode, question, answer, token, ip, ua, email) " + \
                        "values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
                    ), [
                        flask.session['c_id'],
                        flask.session['c_pw'],
                        get_time(),
                        db_data[0][0],
                        flask.session['c_que'],
                        flask.session['c_ans'],
                        application_token,
                        ip,
                        user_agent,
                        flask.session['c_email']
                    ])
                    conn.commit()
    
                    for i in re_set_list:
                        flask.session.pop(i, None)

                    return redirect('/application_submitted')
                else:
                    curs.execute(db_change("insert into user (id, pw, acl, date, encode) values (?, ?, ?, ?, ?)"), [
                        flask.session['c_id'],
                        flask.session['c_pw'],
                        'user' if first == 0 else 'owner',
                        get_time(),
                        db_data[0][0]
                    ])

                curs.execute(db_change("insert into user_set (name, id, data) values ('email', ?, ?)"), [
                    flask.session['c_id'],
                    flask.session['c_email']
                ])
                curs.execute(db_change("insert into ua_d (name, ip, ua, today, sub) values (?, ?, ?, ?, '')"), [
                    flask.session['c_id'],
                    ip,
                    user_agent,
                    get_time()
                ])

                flask.session['id'] = flask.session['c_id']
                flask.session['head'] = ''

                conn.commit()
            else:
                curs.execute(db_change('delete from user_set where name = "email" and id = ?'), [ip])
                curs.execute(db_change('insert into user_set (name, id, data) values ("email", ?, ?)'), [ip, flask.session['c_email']])

                first = 0

            for i in re_set_list:
                flask.session.pop(i, None)

            return redirect('/change') if first == 0 else redirect('/setting') 
        else:
            for i in re_set_list:
                flask.session.pop(i, None)

            return redirect('/user')
    else:
        curs.execute(db_change('select data from other where name = "check_key_text"'))
        sql_d = curs.fetchall()
        b_text = (sql_d[0][0] + '<hr class="main_hr">') if sql_d and sql_d[0][0] != '' else ''

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('check_key'), wiki_set(), custom(), other2([0, 0])],
            data = '''
                <form method="post">
                    ''' + b_text + '''
                    <input placeholder="''' + load_lang('key') + '''" name="key" type="text">
                    <hr class="main_hr">
                    <button type="submit">''' + load_lang('save') + '''</button>
                </form>
            ''',
            menu = [['user', load_lang('return')]]
        ))