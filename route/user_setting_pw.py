from .tool.func import *

def user_setting_pw():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if ban_check() == 1:
            return re_error('/ban')

        ip = ip_check()
        if ip_or_user(ip) != 0:
            return redirect('/login')

        if flask.request.method == 'POST':
            now_pw = flask.request.form.get('pw4', None)
            new_pw = flask.request.form.get('pw2', None)
            re_pw = flask.request.form.get('pw3', None)
            if now_pw and new_pw and re_pw:
                if new_pw != re_pw:
                    return re_error('/error/20')

                curs.execute(db_change("" + \
                    "select name, data from user_set " + \
                    "where id = ? and (name = 'encode' or name = 'pw')" + \
                ""), [
                    flask.session['id']
                ])
                sql_data = curs.fetchall()
                if not sql_data:
                    return re_error('/error/2')
                else:
                    user = {}
                    for i in sql_data:
                        user[i[0]] = i[1]

                if pw_check(
                    now_pw,
                    user['pw'], 
                    user['encode'], 
                    ip
                ) != 1:
                    return re_error('/error/10')

                curs.execute(db_change(
                    "update user_set set data = ? where id = ? and name = 'pw'"
                ), [
                    pw_encode(new_pw), 
                    ip
                ])

            return redirect('/user')
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('password_change'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <input placeholder="''' + load_lang('now_password') + '''" name="pw4" type="password">
                        <hr class="main_hr">
                        <input placeholder="''' + load_lang('new_password') + '''" name="pw2" type="password">
                        <hr class="main_hr">
                        <input placeholder="''' + load_lang('password_confirm') + '''" name="pw3" type="password">
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['change', load_lang('return')]]
            ))