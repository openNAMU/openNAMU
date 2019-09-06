from .tool.func import *
import pymysql

def user_setting_2(conn, server_init):
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
            auto_list = ['email', 'skin', 'lang']

            for auto_data in auto_list:
                if flask.request.form.get(auto_data, '') != '':
                    curs.execute('select data from user_set where name = %s and id = %s', [auto_data, ip])
                    if curs.fetchall():
                        curs.execute("update user_set set data = %s where name = %s and id = %s", [flask.request.form.get(auto_data, ''), auto_data, ip])
                    else:
                        curs.execute("insert into user_set (name, id, data) values (%s, %s, %s)", [auto_data, ip, flask.request.form.get(auto_data, '')])

            
            
            return redirect('/change')
        else:        
            curs.execute('select data from user_set where name = "email" and id = %s', [ip])
            data = curs.fetchall()
            if data:
                email = data[0][0]
            else:
                email = '-'

            div2 = load_skin()
            div3 = ''

            curs.execute('select data from user_set where name = "lang" and id = %s', [flask.session['id']])
            data = curs.fetchall()
            if not data:
                curs.execute('select data from other where name = "language"')
                data = curs.fetchall()
                if not data:
                    data = [['en-US']]

            for lang_data in support_language:
                if data and data[0][0] == lang_data:
                    div3 = '<option value="' + lang_data + '">' + lang_data + '</option>' + div3
                else:
                    div3 += '<option value="' + lang_data + '">' + lang_data + '</option>'

            oauth_provider = load_oauth('_README')['support']
            oauth_content = '<ul>'
            for i in range(len(oauth_provider)):
                curs.execute('select name, picture from oauth_conn where wiki_id = %s and provider = %s', [flask.session['id'], oauth_provider[i]])
                oauth_data = curs.fetchall()
                if len(oauth_data) == 1:
                    oauth_content += '<li>{} - {}</li>'.format(oauth_provider[i].capitalize(), load_lang('connection') + ' : <img src="{}" width="17px" height="17px">{}'.format(oauth_data[0][1], oauth_data[0][0]))
                else:
                    oauth_content += '<li>{} - {}</li>'.format(oauth_provider[i].capitalize(), load_lang('connection') + ' : <a href="/oauth/{}/init">{}</a>'.format(oauth_provider[i], load_lang('connect')))
            
            oauth_content += '</ul>'

            http_warring = '<hr class=\"main_hr\"><span>' + load_lang('http_warring') + '</span>'

            return easy_minify(flask.render_template(skin_check(),    
                imp = [load_lang('user_setting'), wiki_set(), custom(), other2([0, 0])],
                data = '''
                    <form method="post">
                        <span>''' + load_lang('id') + ''' : ''' + ip + '''</span>
                        <hr class=\"main_hr\">
                        <a href="/pw_change">(''' + load_lang('password_change') + ''')</a>
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('email') + ''' : ''' + email + '''</span> <a href="/email_change">(''' + load_lang('email_change') + ''')</a>
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