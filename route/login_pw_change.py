from .tool.func import *

def login_pw_change_2(conn):
    curs = conn.cursor()
    
    if ban_check() == 1:
        return re_error('/ban')

    ip = ip_check()
    if ip_or_user(ip) != 0:
        return redirect('/login')

    if flask.request.method == 'POST':
        if flask.request.form.get('pw4', None) and flask.request.form.get('pw2', None):
            if flask.request.form.get('pw2', None) != flask.request.form.get('pw3', None):
                return re_error('/error/20')

            curs.execute(db_change("select pw, encode from user where id = ?"), [flask.session['id']])
            user = curs.fetchall()
            if not user:
                return re_error('/error/2')
               
            pw_check_d = pw_check(
                flask.request.form.get('pw4', ''), 
                user[0][0],
                user[0][1],
                ip
            )
            if pw_check_d != 1:
                return re_error('/error/10')

            hashed = pw_encode(flask.request.form.get('pw2', None))
                
            curs.execute(db_change("update user set pw = ? where id = ?"), [hashed, ip])

            return redirect('/user')
    else:
        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('password_change'), wiki_set(), custom(), other2([0, 0])],
            data = '''
                <form method="post">
                    <input placeholder="''' + load_lang('now_password') + '''" name="pw4" type="password">
                    <hr class=\"main_hr\">
                    <input placeholder="''' + load_lang('new_password') + '''" name="pw2" type="password">
                    <hr class=\"main_hr\">
                    <input placeholder="''' + load_lang('password_confirm') + '''" name="pw3" type="password">
                    <hr class=\"main_hr\">
                    <button type="submit">''' + load_lang('save') + '''</button>
                </form>
            ''',
            menu = [['change', load_lang('return')]]
        ))