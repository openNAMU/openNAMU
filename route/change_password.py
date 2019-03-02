from .tool.func import *

def change_password_2(conn, server_init):
    curs = conn.cursor()

    support_language = server_init.server_set_var['language']['list']

    if ban_check() == 1:
        return re_error('/ban')

    if custom()[2] == 0:
        return redirect('/login')

    ip = ip_check()
    user_state = flask.request.args.get('user', 'ip')
    
    if user_state == 'ip':
        if flask.request.method == 'POST':    
            if flask.request.form.get('pw4', None) and flask.request.form.get('pw2', None):
                if flask.request.form.get('pw2', None) != flask.request.form.get('pw3', None):
                    return re_error('/error/20')

                curs.execute("select pw, encode from user where id = ?", [flask.session['id']])
                user = curs.fetchall()
                if not user:
                    return re_error('/error/2')
                
                pw_check_d = pw_check(
                    flask.request.form.get('pw4', ''), 
                    user[0][0],
                    user[0][1],
                    flask.request.form.get('id', None)
                )
                if pw_check_d != 1:
                    return re_error('/error/10')

                hashed = pw_encode(flask.request.form.get('pw2', None))
                
                curs.execute("update user set pw = ? where id = ?", [hashed, flask.session['id']])

            auto_list = ['email', 'skin', 'lang']

            for auto_data in auto_list:
                curs.execute('select data from user_set where name = ? and id = ?', [auto_data, ip])
                if curs.fetchall():
                    curs.execute("update user_set set data = ? where name = ? and id = ?", [flask.request.form.get(auto_data, ''), auto_data, ip])
                else:
                    curs.execute("insert into user_set (name, id, data) values (?, ?, ?)", [auto_data, ip, flask.request.form.get(auto_data, '')])

            conn.commit()
            
            return redirect('/change')
        else:        
            curs.execute('select data from user_set where name = "email" and id = ?', [ip])
            data = curs.fetchall()
            if data:
                email = data[0][0]
            else:
                email = ''

            div2 = load_skin()
            
            div3 = ''
            var_div3 = ''

            curs.execute('select data from user_set where name = "lang" and id = ?', [flask.session['id']])
            data = curs.fetchall()

            for lang_data in support_language:
                if data and data[0][0] == lang_data:
                    div3 = '<option value="' + lang_data + '">' + lang_data + '</option>'
                else:
                    var_div3 += '<option value="' + lang_data + '">' + lang_data + '</option>'

            div3 += var_div3

            oauth_provider = load_oauth('_README')['support']
            oauth_content = '<ul>'
            for i in range(len(oauth_provider)):
                curs.execute('select name, picture from oauth_conn where wiki_id = ? and provider = ?', [flask.session['id'], oauth_provider[i]])
                oauth_data = curs.fetchall()
                if len(oauth_data) == 1:
                    oauth_content += '<li>{} - {}</li>'.format(oauth_provider[i], load_lang('connection') + ' : <img src="{}" width="17px" height="17px">{}'.format(oauth_data[0][1], oauth_data[0][0]))
                else:
                    oauth_content += '<li>{} - {}</li>'.format(oauth_provider[i], load_lang('connection') + ' : <a href="/oauth/{}/init">{}</a>'.format(oauth_provider[i], load_lang('connect')))
            
            oauth_content += '</ul>'

            forwarded_protocol = request.headers.get('X-Forwarded-Proto', None)
            if forwarded_protocol == 'http':
                http_warring = '<hr class=\"main_hr\"><span>' + load_lang('http_warring') + '</span>'
            else:
                http_warring = ''

            return easy_minify(flask.render_template(skin_check(),    
                imp = [load_lang('user_setting'), wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <form method="post">
                            <span>id : ''' + ip + '''</span>
                            <hr class=\"main_hr\">
                            <input placeholder="''' + load_lang('now_password') + '''" name="pw4" type="password">
                            <hr class=\"main_hr\">
                            <input placeholder="''' + load_lang('new_password') + '''" name="pw2" type="password">
                            <hr class=\"main_hr\">
                            <input placeholder="''' + load_lang('password_confirm') + '''" name="pw3" type="password">
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('skin') + '''</span>
                            <hr class=\"main_hr\">
                            <select name="skin">''' + div2 + '''</select>
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('language') + '''</span>
                            <hr class=\"main_hr\">
                            <select name="lang">''' + div3 + '''</select>
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('oauth_connection') + '''</span>
                            ''' + oauth_content + '''
                            <hr class=\"main_hr\">
                            <button type="submit">''' + load_lang('save') + '''</button>
                            ''' + http_warring + '''
                        </form>
                        ''',
                menu = [['user', load_lang('return')]]
            ))
    else:
        pass