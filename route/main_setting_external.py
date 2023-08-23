from .tool.func import *

def main_setting_external():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/ban')
        
        i_list = [
            'recaptcha',
            'sec_re',
            'smtp_server',
            'smtp_port',
            'smtp_security',
            'smtp_email',
            'smtp_pass',
            'recaptcha_ver',
            'oauth_client_id',
            'email_have'
        ]

        if flask.request.method == 'POST':
            for data in i_list:
                into_data = flask.request.form.get(data, '')

                curs.execute(db_change("update other set data = ? where name = ?"), [into_data, data])

            conn.commit()

            admin_check(None, 'edit_set (external)')

            return redirect('/setting/external')
        else:
            d_list = []

            x = 0

            for i in i_list:
                curs.execute(db_change('select data from other where name = ?'), [i])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list += [sql_d[0][0]]
                else:
                    curs.execute(db_change('insert into other (name, data, coverage) values (?, ?, "")'), [i, ''])
                    d_list += ['']

                x += 1

            conn.commit()

            security_radios = ''
            for i in ['tls', 'starttls', 'plain']:
                if d_list[4] == i:
                    security_radios = '<option value="' + i + '">' + i + '</option>' + security_radios
                else:
                    security_radios += '<option value="' + i + '">' + i + '</option>'

            re_ver_list = {
                '' : 'reCAPTCHA v2',
                'v3' : 'reCAPTCHA v3',
                'h' : 'hCAPTCHA',
                'cf' : 'Turnstile'
            }
            re_ver = ''
            for i in re_ver_list:
                if d_list[7] == i:
                    re_ver = '<option value="' + i + '">' + re_ver_list[i] + '</option>' + re_ver
                else:
                    re_ver += '<option value="' + i + '">' + re_ver_list[i] + '</option>'

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('ext_api_req_set'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = render_simple_set('''
                    <form method="post">
                        <h2>''' + load_lang('captcha') + '''</h2>
                        <a href="https://www.google.com/recaptcha/">(''' + load_lang('recaptcha') + ''')</a> <a href="https://www.hcaptcha.com/">(''' + load_lang('hcaptcha') + ''')</a>
                        <hr class="main_hr">

                        <span>''' + load_lang('public_key') + '''</span>
                        <hr class="main_hr">
                        <input name="recaptcha" value="''' + html.escape(d_list[0]) + '''">
                        <hr class="main_hr">

                        <span>''' + load_lang('secret_key') + '''</span>
                        <hr class="main_hr">
                        <input name="sec_re" value="''' + html.escape(d_list[1]) + '''">
                        <hr class="main_hr">

                        <span>''' + load_lang('version') + '''</span>
                        <hr class="main_hr">
                        <select name="recaptcha_ver">
                            ''' + re_ver + '''
                        </select>

                        <h2>''' + load_lang('email_setting') + '''</h2>
                        <a href="/setting/phrase#s-6">(''' + load_lang('text_setting') + ''')</a>
                        <hr class="main_hr">

                        <input type="checkbox" name="email_have" ''' + ('checked' if d_list[9] != '' else '')  + '''> ''' + load_lang('email_required') + '''

                        <h3>''' + load_lang('smtp_setting') + '''</h3>
                        <a href="https://support.google.com/mail/answer/7126229">(Google)</a>
                        <hr class="main_hr">

                        <span>''' + load_lang('smtp_server') + '''</span>
                        <hr class="main_hr">
                        <input name="smtp_server" value="''' + html.escape(d_list[2]) + '''">
                        <hr class="main_hr">

                        <span>''' + load_lang('smtp_port') + '''</span>
                        <hr class="main_hr">
                        <input name="smtp_port" value="''' + html.escape(d_list[3]) + '''">
                        <hr class="main_hr">

                        <span>''' + load_lang('smtp_security') + '''</span>
                        <hr class="main_hr">
                        <select name="smtp_security">
                            ''' + security_radios + '''
                        </select>
                        <hr class="main_hr">

                        <span>''' + load_lang('smtp_username') + '''</span>
                        <hr class="main_hr">
                        <input name="smtp_email" value="''' + html.escape(d_list[5]) + '''">
                        <hr class="main_hr">

                        <span>''' + load_lang('smtp_password') + '''</span>
                        <hr class="main_hr">
                        <input type="password" name="smtp_pass" value="''' + html.escape(d_list[6]) + '''">

                        <h2>''' + load_lang('oauth') + ''' (''' + load_lang('not_working') + ''')</h2>
                        <a href="https://developers.google.com/identity/protocols/oauth2">(Google)</a>
                        <hr class="main_hr">

                        <span>''' + load_lang('oauth_client_id') + '''</span>
                        <hr class="main_hr">
                        <input name="oauth_client_id" value="''' + html.escape(d_list[8]) + '''">
                        <hr class="main_hr">

                        <hr class="main_hr">
                        <button id="opennamu_save_button" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                '''),
                menu = [['setting', load_lang('return')]]
            ))