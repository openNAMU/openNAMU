from .tool.func import *

def user_setting_user_name(user_name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if user_name != '':
            if admin_check() != 1:
                return re_error('/error/3')
            else:
                ip = user_name
    
        if ip_or_user(ip) == 0:
            if flask.request.method == 'POST':
                auto_data = ['user_name', flask.request.form.get('new_user_name', '')]
                if do_user_name_check(auto_data[1]) == 1:
                    return re_error('/error/8')

                curs.execute(db_change('select data from user_set where name = ? and id = ?'), [auto_data[0], ip])
                if curs.fetchall():
                    curs.execute(db_change("update user_set set data = ? where name = ? and id = ?"), [auto_data[1], auto_data[0], ip])
                else:
                    curs.execute(db_change("insert into user_set (name, id, data) values (?, ?, ?)"), [auto_data[0], ip, auto_data[1]])

                if user_name != '':
                    return redirect('/change/user_name/' + url_pas(user_name))
                else:
                    return redirect('/change/user_name')
            else:
                user_name = ip

                curs.execute(db_change("select data from user_set where id = ? and name = 'user_name'"), [ip])
                db_data = curs.fetchall()
                if db_data and db_data[0][0] != '':
                    user_name = db_data[0][0]

                return easy_minify(flask.render_template(skin_check(),
                    imp = [load_lang('change_user_name'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                    data = '''
                        <form method="post">
                            <input name="new_user_name" placeholder="''' + load_lang('user_name') + '''" value="''' + html.escape(user_name) + '''">
                            <hr class="main_hr">
                            <button id="opennamu_save_button" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                    ''',
                    menu = [['change', load_lang('return')]]
                ))
        else:
            return redirect('/login')