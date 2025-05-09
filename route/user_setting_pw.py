from .tool.func import *

async def user_setting_pw():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if (await ban_check())[0] == 1:
            return await re_error(conn, 0)

        ip = ip_check()
        if ip_or_user(ip) != 0:
            return redirect(conn, '/login')

        if flask.request.method == 'POST':
            user_pw_now = flask.request.form.get('password_now', '')
            user_pw = flask.request.form.get('password_new', '')
            user_repeat = flask.request.form.get('password_new_repeat', '')
        
            # PW 검증
            if user_pw == '':
                return await re_error(conn, 27)

            if user_pw != user_repeat:
                return await re_error(conn, 20)
    
            # PW 길이 제한
            curs.execute(db_change("select data from other where name = 'password_min_length'"))
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '':
                password_min_length = int(number_check(db_data[0][0]))
                if password_min_length > len(user_pw):
                    return await re_error(conn, 40)

            curs.execute(db_change("select data from user_set where id = ? and name = 'pw'"), [ip])
            db_data = curs.fetchall()
            if not db_data:
                return await re_error(conn, 2)
            else:
                db_user_pw = db_data[0][0]
                
            curs.execute(db_change("select data from user_set where id = ? and name = 'encode'"), [ip])
            db_data = curs.fetchall()
            if not db_data:
                return await re_error(conn, 2)
            else:
                db_user_encode = db_data[0][0]
                
            if pw_check(conn, user_pw_now, db_user_pw, db_user_encode, ip) != 1:
                return await re_error(conn, 10)

            curs.execute(db_change("update user_set set data = ? where id = ? and name = 'pw'"), [pw_encode(conn, user_pw), ip])

            return redirect(conn, '/user')
        else:
            curs.execute(db_change("select data from other where name = 'password_min_length'"))
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '':
                password_min_length = ' (' + get_lang(conn, 'password_min_length') + ' : ' + db_data[0][0] + ')'
            else:
                password_min_length = ''
            
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'password_change'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <input placeholder="''' + get_lang(conn, 'now_password') + '''" name="password_now" type="password">
                        <hr class="main_hr">
                        
                        <input placeholder="''' + get_lang(conn, 'new_password') + password_min_length + '''" name="password_new" type="password">
                        <hr class="main_hr">
                        
                        <input placeholder="''' + get_lang(conn, 'password_confirm') + '''" name="password_new_repeat" type="password">
                        <hr class="main_hr">
                        
                        <button type="submit">''' + get_lang(conn, 'save') + '''</button>
                        
                        ''' + http_warning(conn) + '''
                    </form>
                ''',
                menu = [['change', get_lang(conn, 'return')]]
            ))