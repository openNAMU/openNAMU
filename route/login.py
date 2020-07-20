from .tool.func import *

def login_2(conn):
    curs = conn.cursor()

    ip = ip_check()
    if ip_or_user(ip) == 0:
        return redirect('/user')

    if ban_check(tool = 'login') == 1:
        return re_error('/ban')

    if flask.request.method == 'POST':
        if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        agent = flask.request.headers.get('User-Agent')

        curs.execute(db_change("select pw, encode from user where id = ?"), [flask.request.form.get('id', None)])
        user = curs.fetchall()
        if not user:
            return re_error('/error/2')

        pw_check_d = pw_check(
            flask.request.form.get('pw', ''),
            user[0][0],
            user[0][1],
            flask.request.form.get('id', None)
        )
        if pw_check_d != 1:
            return re_error('/error/10')

        flask.session['state'] = 1
        flask.session['id'] = flask.request.form.get('id', None)

        curs.execute(db_change("select css from custom where user = ?"), [flask.request.form.get('id', None)])
        css_data = curs.fetchall()
        if css_data:
            flask.session['head'] = css_data[0][0]
        else:
            flask.session['head'] = ''

        curs.execute(db_change("insert into ua_d (name, ip, ua, today, sub) values (?, ?, ?, ?, '')"), [flask.request.form.get('id', None), ip_check(1), agent, get_time()])

        conn.commit()

        return redirect('/user')
    else:
        http_warring = '<hr class=\"main_hr\"><span>' + load_lang('http_warring') + '</span>'

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('login'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        <input placeholder="''' + load_lang('id') + '''" name="id" type="text">
                        <hr class=\"main_hr\">
                        <input placeholder="''' + load_lang('password') + '''" name="pw" type="password">
                        <hr class=\"main_hr\">
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('login') + '''</button>
                        ''' + http_warring + '''
                    </form>
                    ''',
            menu = [['user', load_lang('return')]]
        ))