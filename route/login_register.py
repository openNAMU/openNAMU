from .tool.func import *

def login_register_2():
    

    if ban_check() == 1:
        return re_error('/ban')

    if custom()[2] != 0:
        return redirect('/user')

    if not admin_check() == 1:
        sqlQuery('select data from other where name = "reg"')
        set_d = sqlQuery("fetchall")
        if set_d and set_d[0][0] == 'on':
            return re_error('/ban')
    
    if flask.request.method == 'POST': 
        if captcha_post(flask.request.form.get('g-recaptcha-response', '')) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        if flask.request.form.get('pw', None) != flask.request.form.get('pw2', None):
            return re_error('/error/20')

        if re.search('(?:[^A-Za-zㄱ-힣0-9 ])', flask.request.form.get('id', None)):
            return re_error('/error/8')
            
        sqlQuery('select html from html_filter where kind = "name"')
        set_d = sqlQuery("fetchall")
        for i in set_d:
            check_r = re.compile(i[0], re.I)
            if check_r.search(flask.request.form.get('id', None)):
                return re_error('/error/8')

        if len(flask.request.form.get('id', None)) > 32:
            return re_error('/error/7')

        sqlQuery("select id from user where id = ?", [flask.request.form.get('id', None)])
        if sqlQuery("fetchall"):
            return re_error('/error/6')

        hashed = pw_encode(flask.request.form.get('pw', None))
        
        sqlQuery('select data from other where name = "email_have"')
        sql_data = sqlQuery("fetchall")
        if sql_data and sql_data[0][0] != '':
            flask.session['c_id'] = flask.request.form.get('id', None)
            flask.session['c_pw'] = hashed
            flask.session['c_key'] = ''.join(random.choice("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") for i in range(16))

            return redirect('/need_email')
        else:
            sqlQuery('select data from other where name = "encode"')
            db_data = sqlQuery("fetchall")

            sqlQuery("select id from user limit 1")
            if not sqlQuery("fetchall"):
                sqlQuery("insert into user (id, pw, acl, date, encode) values (?, ?, 'owner', ?, ?)", [flask.request.form.get('id', None), hashed, get_time(), db_data[0][0]])

                first = 1
            else:
                sqlQuery("insert into user (id, pw, acl, date, encode) values (?, ?, 'user', ?, ?)", [flask.request.form.get('id', None), hashed, get_time(), db_data[0][0]])

                first = 0

            ip = ip_check()
            agent = flask.request.headers.get('User-Agent')

            sqlQuery("insert into ua_d (name, ip, ua, today, sub) values (?, ?, ?, ?, '')", [flask.request.form.get('id', None), ip, agent, get_time()])  

            flask.session['state'] = 1
            flask.session['id'] = flask.request.form.get('id', None)
            flask.session['head'] = ''
                  
            sqlQuery("commit")
            
            if first == 0:
                return redirect('/change')
            else:
                return redirect('/setting')
    else:        
        contract = ''
        
        sqlQuery('select data from other where name = "contract"')
        data = sqlQuery("fetchall")
        if data and data[0][0] != '':
            contract = data[0][0] + '<hr class=\"main_hr\">'
        
        http_warring = '<hr class=\"main_hr\"><span>' + load_lang('http_warring') + '</span>'

        return easy_minify(flask.render_template(skin_check(),    
            imp = [load_lang('register'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        ''' + contract + '''
                        <input placeholder="''' + load_lang('id') + '''" name="id" type="text">
                        <hr class=\"main_hr\">
                        <input placeholder="''' + load_lang('password') + '''" name="pw" type="password">
                        <hr class=\"main_hr\">
                        <input placeholder="''' + load_lang('password_confirm') + '''" name="pw2" type="password">
                        <hr class=\"main_hr\">
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('save') + '''</button>
                        ''' + http_warring + '''
                    </form>
                    ''',
            menu = [['user', load_lang('return')]]
        ))