from .tool.func import *
import pymysql

def login_2(conn):
    curs = conn.cursor()

    if custom()[2] != 0:
        return redirect('/user')
    
    if ban_check(tool = 'login') == 1:
        return re_error('/ban')
        
    if flask.request.method == 'POST':        
        if captcha_post(flask.request.form.get('g-recaptcha-response', '')) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)
            
        agent = flask.request.headers.get('User-Agent')

        curs.execute("select pw, encode from user where id = %s", [flask.request.form.get('id', None)])
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
        
        curs.execute("select css from custom where user = %s", [flask.request.form.get('id', None)])
        css_data = curs.fetchall()
        if css_data:
            flask.session['head'] = css_data[0][0]
        else:
            flask.session['head'] = ''

        curs.execute("insert into ua_d (name, ip, ua, today, sub) values (%s, %s, %s, %s, '')", [flask.request.form.get('id', None), ip_check(1), agent, get_time()])

        conn.commit()
        
        return redirect('/user')  
    else:
        oauth_check = 0
        oauth_content = '<hr class=\"main_hr\"><div class="oauth-wrapper"><ul class="oauth-list">'
        oauth_supported = load_oauth('_README')['support']
        for i in range(len(oauth_supported)):
            oauth_data = load_oauth(oauth_supported[i])
            if oauth_data['client_id'] != '' and oauth_data['client_secret'] != '':
                oauth_content += '''
                    <li>
                        <a href="/oauth/{}/init">
                            <div class="oauth-btn oauth-btn-{}">
                                <div class="oauth-btn-logo oauth-btn-{}"></div>
                                {}
                            </div>
                        </a>
                    </li>
                '''.format(
                    oauth_supported[i], 
                    oauth_supported[i], 
                    oauth_supported[i], 
                    load_lang('oauth_signin_' + oauth_supported[i])
                )

                oauth_check = 1
        
        oauth_content += '</ul></div>'

        if oauth_check == 0:
            oauth_content = ''

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
                        ''' + oauth_content + '''
                        ''' + http_warring + '''
                    </form>
                    ''',
            menu = [['user', load_lang('return')]]
        ))