from .tool.func import *

async def user_setting():
    with get_db_connect() as conn:
        curs = conn.cursor()

        support_language = ['default'] + get_init_set_list()['language']['list']
        
        ip = ip_check()

        if ip_or_user(ip) == 0:
            if flask.request.method == 'POST':
                auto_list = [
                    ['skin', flask.request.form.get('skin', '')], 
                    ['lang', flask.request.form.get('lang', '')],
                    ['user_title', flask.request.form.get('user_title', '')],
                    ['sub_user_name' , flask.request.form.get('sub_user_name', '')]
                ]
                if not auto_list[2][1] in await get_user_title_list(conn, ip):
                    auto_list[2][1] = ''

                twofa_on = flask.request.form.get('2fa', '')
                if twofa_on != '':
                    twofa_pw = flask.request.form.get('2fa_pw', '')
                    if twofa_pw != '':
                        twofa_pw = pw_encode(conn, twofa_pw)

                        curs.execute(db_change("select data from user_set where id = ? and name = 'encode'"), [ip])
                        twofa_encode = curs.fetchall()[0][0]
                        
                        auto_list += [['2fa', 'on'], ['2fa_pw', twofa_pw], ['2fa_pw_encode', twofa_encode]]
                    else:
                        auto_list += [['2fa', 'on']]
                else:
                    auto_list += [['2fa', '']]

                for auto_data in auto_list:
                    curs.execute(db_change('select data from user_set where name = ? and id = ?'), [auto_data[0], ip])
                    if curs.fetchall():
                        curs.execute(db_change("update user_set set data = ? where name = ? and id = ?"), [auto_data[1], auto_data[0], ip])
                    else:
                        curs.execute(db_change("insert into user_set (name, id, data) values (?, ?, ?)"), [auto_data[0], ip, auto_data[1]])

                return redirect(conn, '/change')
            else:
                curs.execute(db_change('select data from user_set where name = "email" and id = ?'), [ip])
                data = curs.fetchall()
                email = data[0][0] if data and data[0][0] != '' else '-'

                curs.execute(db_change('select data from user_set where name = "random_key" and id = ?'), [ip])
                data = curs.fetchall()
                ramdom_key = data[0][0] if data and data[0][0] != '' else '-'

                curs.execute(db_change('select data from user_set where name = "skin" and id = ?'), [ip])
                data = curs.fetchall()
                div2 = load_skin(conn, data[0][0] if data else '', 0, 1)

                curs.execute(db_change('select data from user_set where name = "lang" and id = ?'), [ip])
                data = curs.fetchall()
                data = [['default']] if not data else data
                div3 = ''
                for lang_data in support_language:
                    see_data = lang_data if lang_data != 'default' else get_lang(conn, 'default')

                    if data and data[0][0] == lang_data:
                        div3 = '<option value="' + lang_data + '">' + see_data + '</option>' + div3
                    else:
                        div3 += '<option value="' + lang_data + '">' + see_data + '</option>'

                curs.execute(db_change('select data from user_set where name = "user_title" and id = ?'), [ip])
                data = curs.fetchall()
                data = [['']] if not data else data
                user_title_list = await get_user_title_list(conn, ip)
                div4 = ''
                for user_title in user_title_list:                
                    if data and data[0][0] == user_title:
                        div4 = '<option value="' + user_title + '">' + user_title_list[user_title] + '</option>' + div4
                    else:
                        div4 += '<option value="' + user_title + '">' + user_title_list[user_title] + '</option>'

                curs.execute(db_change('select data from user_set where name = "2fa" and id = ?'), [ip])
                fa_data = curs.fetchall()
                fa_data = fa_data[0][0] if fa_data and fa_data[0][0] != '' else ''
                fa_data_select = ''
                fa_data_sp_list = [[get_lang(conn, 'off'), ''], [get_lang(conn, 'password'), 'on']]
                for fa_data_get in fa_data_sp_list:
                    fa_data_selected = ''
                    if fa_data == fa_data_get[1]:
                        fa_data_selected = 'selected'

                    fa_data_select += '<option ' + fa_data_selected + ' value="' + fa_data_get[1] + '">' + fa_data_get[0] + '</option>'

                curs.execute(db_change('select data from user_set where name = "2fa_pw" and id = ?'), [ip])
                fa_data_pw = curs.fetchall()
                fa_data_pw = get_lang(conn, '2fa_password_change') if fa_data_pw else get_lang(conn, '2fa_password')

                curs.execute(db_change('select data from user_set where name = "user_name" and id = ?'), [ip])
                db_data = curs.fetchall()
                user_name = db_data[0][0] if db_data else ip

                curs.execute(db_change('select data from user_set where name = "sub_user_name" and id = ?'), [ip])
                db_data = curs.fetchall()
                sub_user_name = db_data[0][0] if db_data else ''

                return easy_minify(conn, flask.render_template(skin_check(conn),
                    imp = [get_lang(conn, 'user_setting'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                    data = '''
                        <form method="post">
                            <div id="opennamu_get_user_info">''' + html.escape(ip) + '''</div>
                            <hr class="main_hr">
                            <a href="/change/pw">(''' + get_lang(conn, 'password_change') + ''')</a>
                            <hr class="main_hr">
                            <span>''' + get_lang(conn, 'email') + ''' : ''' + email + '''</span> <a href="/change/email">(''' + get_lang(conn, 'email_change') + ''')</a> <a href="/change/email/delete">(''' + get_lang(conn, 'email_delete') + ''')</a>
                            <hr class="main_hr">
                            <span>''' + get_lang(conn, 'password_instead_key') + ''' : ''' + ramdom_key + ''' <a href="/change/key">(''' + get_lang(conn, 'key_change') + ''')</a> <a href="/change/key/delete">(''' + get_lang(conn, 'key_delete') + ''')</a></span>
                            <h2>''' + get_lang(conn, 'main') + '''</h2>
                            <a href="/change/head">(''' + get_lang(conn, 'user_head') + ''')</a> <a href="/change/top_menu">(''' + get_lang(conn, 'user_added_menu') + ''')</a>
                            <hr class="main_hr">
                            <span>''' + get_lang(conn, 'skin') + '''</span>
                            <hr class="main_hr">
                            <select name="skin">''' + div2 + '''</select>
                            <hr class="main_hr">
                            <a href="/change/skin_set">(''' + get_lang(conn, 'skin_set') + ''')</a> <a href="/change/skin_set/main">(''' + get_lang(conn, 'main_skin_set') + ''')</a>
                            <hr class="main_hr">
                            <span>''' + get_lang(conn, 'language') + '''</span>
                            <hr class="main_hr">
                            <select name="lang">''' + div3 + '''</select>
                            <hr class="main_hr">
                            <span>''' + get_lang(conn, 'user_title') + '''</span>
                            <hr class="main_hr">
                            <select name="user_title">''' + div4 + '''</select>
                            <h2>''' + get_lang(conn, '2fa') + '''</h2>
                            <select name="2fa" id="twofa_check_input">''' + fa_data_select + '''</select>
                            <hr class="main_hr">
                            <input type="password" name="2fa_pw" placeholder="''' + fa_data_pw + '''">
                            <h2>''' + get_lang(conn, 'main_user_name') + '''</h2>
                            <a href="/change/user_name">(''' + get_lang(conn, 'change_user_name') + ''')</a>
                            <hr class="main_hr">
                            ''' + get_lang(conn, 'user_name') + ''' : ''' + html.escape(user_name) + '''
                            <h2>''' + get_lang(conn, 'sub_user_name') + '''</h2>
                            <input name="sub_user_name" value="''' + html.escape(sub_user_name) + '''" placeholder="''' + get_lang(conn, 'sub_user_name') + '''">
                            <hr class="main_hr">
                            <button type="submit">''' + get_lang(conn, 'save') + '''</button>
                            ''' + http_warning(conn) + '''
                        </form>
                    ''',
                    menu = [['user', get_lang(conn, 'return')]]
                ))
        else:
            if flask.request.method == 'POST':
                flask.session['skin'] = flask.request.form.get('skin', '')
                flask.session['lang'] = flask.request.form.get('lang', '')

                return redirect(conn, '/change')
            else:
                div2 = load_skin(conn, 
                    ('' if not 'skin' in flask.session else flask.session['skin']), 
                    0, 
                    1
                )

                data = [['default']] if not 'lang' in flask.session else [[flask.session['lang']]]
                div3 = ''
                for lang_data in support_language:
                    see_data = lang_data if lang_data != 'default' else get_lang(conn, 'default')

                    if data and data[0][0] == lang_data:
                        div3 = '<option value="' + lang_data + '">' + see_data + '</option>' + div3
                    else:
                        div3 += '<option value="' + lang_data + '">' + see_data + '</option>'

                return easy_minify(conn, flask.render_template(skin_check(conn),
                    imp = [get_lang(conn, 'user_setting'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                    data = '''
                        <form method="post">
                            <div id="opennamu_get_user_info">''' + html.escape(ip) + '''</div>
                            <hr class="main_hr">
                            <h2>''' + get_lang(conn, 'main') + '''</h2>
                            <span>''' + get_lang(conn, 'skin') + '''</span>
                            <hr class="main_hr">
                            <select name="skin">''' + div2 + '''</select>
                            <hr class="main_hr">
                            <a href="/change/skin_set">(''' + get_lang(conn, 'skin_set') + ''')</a> <a href="/change/skin_set/main">(''' + get_lang(conn, 'main_skin_set') + ''')</a>
                            <hr class="main_hr">
                            <span>''' + get_lang(conn, 'language') + '''</span>
                            <hr class="main_hr">
                            <select name="lang">''' + div3 + '''</select>
                            <hr class="main_hr">
                            <button type="submit">''' + get_lang(conn, 'save') + '''</button>
                            ''' + http_warning(conn) + '''
                            <hr class="main_hr">
                            <span>''' + get_lang(conn, 'user_head_warning') + '''</span>
                        </form>
                    ''',
                    menu = [['user', get_lang(conn, 'return')]]
                ))
