from .tool.func import *

async def login_find_key():
    with get_db_connect() as conn:
        curs = conn.cursor()
        if flask.request.method == 'POST':
            if await captcha_post(conn, flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return await re_error(conn, 13)
            
            input_key = flask.request.form.get('key', '')
            curs.execute(db_change('select id from user_set where name = "random_key" and data = ?'), [input_key])
            db_data = curs.fetchall()
            if not db_data:
                return redirect(conn, '/user')
            else:
                user_id = db_data[0][0]
            
            key = load_random_key(32)
            curs.execute(db_change("update user_set set data = ? where name = 'pw' and id = ?"), [
                pw_encode(conn, key), 
                user_id
            ])
            
            curs.execute(db_change('select data from user_set where name = "2fa" and id = ?'), [user_id])
            if curs.fetchall():
                curs.execute(db_change("update user_set set data = '' where name = '2fa' and id = ?"), [user_id])
            
            curs.execute(db_change('select data from other where name = "reset_user_text"'))
            sql_d = curs.fetchall()
            b_text = (sql_d[0][0] + '<hr class="main_hr">') if sql_d and sql_d[0][0] != '' else ''
            
            return easy_minify(conn, flask.render_template(skin_check(conn),
                    imp = [get_lang(conn, 'reset_user_ok'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                    data = '' + \
                        b_text + \
                        get_lang(conn, 'id') + ' : ' + user_id + \
                        '<hr class="main_hr">' + \
                        get_lang(conn, 'password') + ' : ' + key + \
                    '',
                    menu = [['user', get_lang(conn, 'return')]]
                ))
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'password_search'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <input placeholder="''' + get_lang(conn, 'key') + '''" name="key" type="password">
                        <hr class="main_hr">
                        ''' + await captcha_get(conn) + '''
                        <button type="submit">''' + get_lang(conn, 'send') + '''</button>
                    </form>
                ''',
                menu = [['user', get_lang(conn, 'return')]]
            ))