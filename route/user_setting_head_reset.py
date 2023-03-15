from .tool.func import *

def user_setting_head_reset():
    with get_db_connect() as conn:
        curs = conn.cursor()

        skin_name = skin_check(1)
        ip = ip_check()

        if flask.request.method == 'POST':
            get_data = ''
            if ip_or_user(ip) == 0:
                curs.execute(db_change("select id from user_set where id = ? and name = ?"), [ip, 'custom_css'])
                if curs.fetchall():
                    curs.execute(db_change("update user_set set data = ? where id = ? and name = ?"), [get_data, ip, 'custom_css'])
                else:
                    curs.execute(db_change("insert into user_set (id, name, data) values (?, ?, ?)"), [ip, 'custom_css', get_data])

                curs.execute(db_change("select id from user_set where id = ? and name = ?"), [ip, 'custom_css_' + skin_name])
                if curs.fetchall():
                    curs.execute(db_change("update user_set set data = ? where id = ? and name = ?"), [get_data, ip, 'custom_css_' + skin_name])
                else:
                    curs.execute(db_change("insert into user_set (id, name, data) values (?, ?, ?)"), [ip, 'custom_css_' + skin_name, get_data])

                conn.commit()

            flask.session['head'] = ''
            flask.session['head_' + skin_name] = ''

            return redirect('/change/head')
        else:
            if ip_or_user(ip) == 0:
                curs.execute(db_change("select data from user_set where id = ? and name = ?"), [ip, 'custom_css'])
                head_data = curs.fetchall()
                data = head_data[0][0] if head_data else ''

                curs.execute(db_change("select data from user_set where id = ? and name = ?"), [ip, 'custom_css_' + skin_name])
                head_data = curs.fetchall()
                data_skin = head_data[0][0] if head_data else ''
            else:
                data = flask.session['head'] if 'head' in flask.session else ''
                data_skin = flask.session['head_' + skin_name] if 'head_' + skin_name in flask.session else ''
            
            return '''
                <form method="post">
                    <style>.main_hr { border: none; }</style>
                    ''' + load_lang('all') + '''
                    <hr class="main_hr">
                    <pre>''' + html.escape(data) + '''</pre>
                    <hr class="main_hr">
                    ''' + skin_name + '''
                    <hr class="main_hr">
                    <pre>''' + html.escape(data_skin) + '''</pre>
                    <hr class="main_hr">
                    <button type="submit">''' + load_lang('reset') + '''</button>
                </form>
            '''