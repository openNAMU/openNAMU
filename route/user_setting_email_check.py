from .tool.func import *

async def user_setting_email_check():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if ip_or_user(ip) != 0:
            return redirect(conn, '/login')

        re_set_list = ['c_key', 'c_email']
        if  not 'c_key' in flask.session or \
            not 'c_email' in flask.session:
            for i in re_set_list:
                flask.session.pop(i, None)

        if  flask.request.method == 'POST':
            ip = ip_check()
            input_key = flask.request.form.get('key', '')
            user_agent = flask.request.headers.get('User-Agent', '')

            if flask.session['c_key'] == input_key:
                curs.execute(db_change('delete from user_set where name = "email" and id = ?'), [ip])
                curs.execute(db_change('insert into user_set (name, id, data) values ("email", ?, ?)'), [ip, flask.session['c_email']])

            for i in re_set_list:
                flask.session.pop(i, None)

            return redirect(conn, '/change')
        else:
            curs.execute(db_change('select data from other where name = "check_key_text"'))
            sql_d = curs.fetchall()
            b_text = (sql_d[0][0] + '<hr class="main_hr">') if sql_d and sql_d[0][0] != '' else ''

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'check_key'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        ''' + b_text + '''
                        <input placeholder="''' + get_lang(conn, 'key') + '''" name="key" type="text">
                        <hr class="main_hr">
                        <button type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                ''',
                menu = [['user', get_lang(conn, 'return')]]
            ))