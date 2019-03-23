from .tool.func import *

def need_email_2(conn, tool):
    curs = conn.cursor()

    if flask.request.method == 'POST':
        if tool == 'need_email':
            if 'c_id' in flask.session:
                main_email = ['naver.com', 'gmail.com', 'daum.net', 'hanmail.net', 'hanmail2.net']
                data = re.search('@([^@]+)$', flask.request.form.get('email', ''))
                if data:
                    data = data.groups()[0]

                    curs.execute("select html from html_filter where html = ? and kind = 'email'", [data])
                    if curs.fetchall() or (data in main_email):
                        curs.execute('select id from user_set where name = "email" and data = ?', [flask.request.form.get('email', '')])
                        if curs.fetchall():
                            flask.session.pop('c_id', None)
                            flask.session.pop('c_pw', None)
                            flask.session.pop('c_key', None)

                            return redirect('/register')
                        else:
                            send_email(flask.request.form.get('email', ''), wiki_set()[0] + ' key', 'key : ' + flask.session['c_key'])
                            flask.session['c_email'] = flask.request.form.get('email', '')

                            return redirect('/check_key')

            return redirect('/register')
        else:
            curs.execute("select id from user where id = ? and email = ?", [flask.request.form.get('id', ''), flask.request.form.get('email', '')])
            if curs.fetchall():
                flask.session['c_key'] = ''.join(random.choice("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") for i in range(16))
                flask.session['c_id'] = flask.request.form.get('id', '')

                send_email(flask.request.form.get('email', ''), wiki_set()[0] + '\'s key', 'Key : ' + flask.session['c_key'])

                return redirect('/check_pass_key')
    else:
        if tool == 'need_email':
            return easy_minify(flask.render_template(skin_check(),    
                imp = [load_lang('email'), wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <a href="/email_filter">(''' + load_lang('email_filter_list') + ''')</a>
                        <hr class=\"main_hr\">
                        <form method="post">
                            <input placeholder="''' + load_lang('email') + '''" name="email" type="text">
                            <hr class=\"main_hr\">
                            <button type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        ''',
                menu = [['user', load_lang('return')]]
            ))
        else:
            return easy_minify(flask.render_template(skin_check(),    
                imp = [load_lang('password_search'), wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <form method="post">
                            <input placeholder="''' + load_lang('id') + '''" name="id" type="text">
                            <hr class=\"main_hr\">
                            <input placeholder="''' + load_lang('email') + '''" name="email" type="text">
                            <hr class=\"main_hr\">
                            <button type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        ''',
                menu = [['user', load_lang('return')]]
            ))