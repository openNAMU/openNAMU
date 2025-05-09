from .tool.func import *

async def user_setting_email():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if ip_or_user(ip) != 0:
            return redirect(conn, '/login')

        if flask.request.method == 'POST':
            # c_key 같은 이름 대신 한 기능에 고유 명칭 부여 필요
            re_set_list = ['c_key']
            flask.session['c_key'] = load_random_key(32)

            user_email = re.sub(r'\\', '', flask.request.form.get('email', ''))
            email_data = re.search(r'@([^@]+)$', user_email)
            if email_data:
                curs.execute(db_change("select html from html_filter where html = ? and kind = 'email'"), [email_data.group(1)])
                if not curs.fetchall():
                    for i in re_set_list:
                        flask.session.pop(i, None)

                    return redirect(conn, '/filter/email_filter')
            else:
                for i in re_set_list:
                    flask.session.pop(i, None)

                return await re_error(conn, 36)

            curs.execute(db_change('select data from other where name = "email_title"'))
            sql_d = curs.fetchall()
            t_text = html.escape(sql_d[0][0]) if sql_d and sql_d[0][0] != '' else ((await wiki_set())[0] + ' key')

            curs.execute(db_change('select data from other where name = "email_text"'))
            sql_d = curs.fetchall()
            if sql_d and sql_d[0][0] != '':
                i_text = html.escape(sql_d[0][0]) + '\n\nKey : ' + flask.session['c_key']
            else:
                i_text = 'Key : ' + flask.session['c_key']

            curs.execute(db_change('select id from user_set where name = "email" and data = ?'), [user_email])
            if curs.fetchall():
                for i in re_set_list:
                    flask.session.pop(i, None)

                return await re_error(conn, 35)

            if await send_email(conn, user_email, t_text, i_text) == 0:
                for i in re_set_list:
                    flask.session.pop(i, None)

                return await re_error(conn, 18)

            flask.session['c_email'] = user_email

            return redirect(conn, '/change/email/check')
        else:
            curs.execute(db_change('select data from other where name = "email_insert_text"'))
            sql_d = curs.fetchall()
            b_text = (sql_d[0][0] + '<hr class="main_hr">') if sql_d and sql_d[0][0] != '' else ''

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'email'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <a href="/filter/email_filter">(''' + get_lang(conn, 'email_filter_list') + ''')</a>
                    <hr class="main_hr">
                    ''' + b_text + '''
                    <form method="post">
                        <input placeholder="''' + get_lang(conn, 'email') + '''" name="email" type="text">
                        <hr class="main_hr">
                        <button type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                ''',
                menu = [['user', get_lang(conn, 'return')]]
            ))