from .tool.func import *

async def main_setting_external():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 0)
        
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

            await acl_check(tool = 'owner_auth', memo = 'edit_set (external)')

            return redirect(conn, '/setting/external')
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

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'ext_api_req_set'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = render_simple_set(conn, '''
                    <form method="post">
                        <h2>''' + get_lang(conn, 'captcha') + '''</h2>
                        <a href="https://www.google.com/recaptcha/">(''' + get_lang(conn, 'recaptcha') + ''')</a> <a href="https://www.hcaptcha.com/">(''' + get_lang(conn, 'hcaptcha') + ''')</a>
                        <hr class="main_hr">

                        <span>''' + get_lang(conn, 'public_key') + '''</span>
                        <hr class="main_hr">
                        <input name="recaptcha" value="''' + html.escape(d_list[0]) + '''">
                        <hr class="main_hr">

                        <span>''' + get_lang(conn, 'secret_key') + '''</span>
                        <hr class="main_hr">
                        <input name="sec_re" value="''' + html.escape(d_list[1]) + '''">
                        <hr class="main_hr">

                        <span>''' + get_lang(conn, 'version') + '''</span>
                        <hr class="main_hr">
                        <select name="recaptcha_ver">
                            ''' + re_ver + '''
                        </select>

                        <h2>''' + get_lang(conn, 'email_setting') + '''</h2>
                        <a href="/setting/phrase#s-6">(''' + get_lang(conn, 'text_setting') + ''')</a>
                        <hr class="main_hr">

                        <label><input type="checkbox" name="email_have" ''' + ('checked' if d_list[9] != '' else '')  + '''> ''' + get_lang(conn, 'email_required') + '''</label>

                        <h3>''' + get_lang(conn, 'smtp_setting') + '''</h3>
                        <a href="https://support.google.com/mail/answer/7126229">(Google)</a>
                        <hr class="main_hr">
                        <a href="/setting/email_test">(''' + get_lang(conn, 'test') + ''')</a>
                        <hr class="main_hr">

                        <span>''' + get_lang(conn, 'smtp_server') + '''</span>
                        <hr class="main_hr">
                        <input name="smtp_server" value="''' + html.escape(d_list[2]) + '''">
                        <hr class="main_hr">

                        <span>''' + get_lang(conn, 'smtp_port') + '''</span>
                        <hr class="main_hr">
                        <input name="smtp_port" value="''' + html.escape(d_list[3]) + '''">
                        <hr class="main_hr">

                        <span>''' + get_lang(conn, 'smtp_security') + '''</span>
                        <hr class="main_hr">
                        <select name="smtp_security">
                            ''' + security_radios + '''
                        </select>
                        <hr class="main_hr">

                        <span>''' + get_lang(conn, 'smtp_username') + '''</span>
                        <hr class="main_hr">
                        <input name="smtp_email" value="''' + html.escape(d_list[5]) + '''">
                        <hr class="main_hr">

                        <span>''' + get_lang(conn, 'smtp_password') + '''</span>
                        <hr class="main_hr">
                        <input type="password" name="smtp_pass" value="''' + html.escape(d_list[6]) + '''">

                        <h2>''' + get_lang(conn, 'oauth') + ''' (''' + get_lang(conn, 'not_working') + ''')</h2>
                        <a href="https://developers.google.com/identity/protocols/oauth2">(Google)</a>
                        <hr class="main_hr">

                        <span>''' + get_lang(conn, 'oauth_client_id') + '''</span>
                        <hr class="main_hr">
                        <input name="oauth_client_id" value="''' + html.escape(d_list[8]) + '''">
                        <hr class="main_hr">

                        <hr class="main_hr">
                        <button id="opennamu_save_button" type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                '''),
                menu = [['setting', get_lang(conn, 'return')]]
            ))