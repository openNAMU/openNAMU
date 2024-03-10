from .tool.func import *

def give_user_ban(name = None, ban_type = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        
        if ban_check(conn, ip = ip, tool = 'login')[0] == 1:
            if ip_or_user(ip) == 1 or admin_check(conn, 'all', None, ip) == 0:
                return re_error(conn, '/ban')
        else:
            if admin_check(conn, 1, None, ip) != 1:
                return re_error(conn, '/error/3')

        if flask.request.method == 'POST':
            time_limit = flask.request.form.get('date', '')
            if re.search(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$', time_limit):
                end = time_limit + ' 00:00:00'
            else:
                end = '0'
            
            regex_get = flask.request.form.get('regex', None)
            why = flask.request.form.get('why', '')

            release = ''
            login = ''
            
            ban_option = flask.request.form.get('ban_option', '')
            if ban_option == 'login_able':
                login = 'O'
            elif ban_option == 'edit_request_able':
                login = 'E'
            elif ban_option == 'release':
                release = '1'

            if ban_type == 'multiple':
                all_user = re.findall(r'([^\n]+)\n', flask.request.form.get('name', 'test').replace('\r', '') + '\n')
            else:
                if name:
                    all_user = [name]
                else:
                    all_user = [flask.request.form.get('name', 'test')]

            for name in all_user:
                if regex_get or ban_type == 'regex':
                    type_d = 'regex' if regex_get else ban_type

                    try:
                        re.compile(name)
                    except:
                        return re_error(conn, '/error/23')
                else:
                    type_d = None

                if type_d:
                    if admin_check(conn, None, 'ban' + (' ' + type_d if type_d else '') + ' (' + name + ')') != 1:
                        return re_error(conn, '/error/3')
                else:
                    if name == ip:
                        if admin_check(conn, 'all', 'ban (' + name + ')') != 1:
                            return re_error(conn, '/error/3')
                    else:
                        if admin_check(conn, 1, 'ban (' + name + ')') != 1:
                            return re_error(conn, '/error/3')

                ban_insert(conn, 
                    name,
                    end,
                    why,
                    login,
                    ip_check(),
                    type_d,
                    1 if release != '' else 0
                )

            return redirect(conn, '/block_log')
        else:
            if ban_type == 'multiple':
                main_name = get_lang(conn, 'multiple_ban')
                n_name = '<textarea class="opennamu_textarea_500" placeholder="' + get_lang(conn, 'name_or_ip_or_regex_multiple') + '" name="name"></textarea><hr class="main_hr">'
            else:
                main_name = get_lang(conn, 'ban')
                n_name = '<input placeholder="' + get_lang(conn, 'name_or_ip_or_regex') + '" value="' + (name if name else '') + '" name="name"><hr class="main_hr">'

            now = 0
            
            if ban_type == 'multiple':
                action = 'action="/auth/give/ban_multiple"'
            else:
                action = 'action="/auth/give/ban"'
                
            date_value = ''
            info_data = ''
            if name:
                curs.execute(db_change("select end from rb where block = ? and ongoing = '1'"), [name])
                db_data = curs.fetchall()
                if db_data and db_data[0][0] != '':
                    date_value = db_data[0][0].split()[0]

                info_data = '''
                    <div id="opennamu_get_user_info">''' + html.escape(name) + '''</div>
                    <hr class="main_hr">
                '''

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [main_name, wiki_set(conn), wiki_custom(conn), wiki_css([now, 0])],
                data = info_data + '''
                    <form method="post" ''' + action + '''>
                        ''' + n_name + '''
                        <input type="checkbox" name="regex" ''' + ('checked' if ban_type == 'regex' else '') + '> ' + get_lang(conn, 'regex') + '''
                        <hr class="main_hr">
                        <input type="date" value="''' + date_value + '''" name="date" pattern="\\d{4}-\\d{2}-\\d{2}">
                        <hr class="main_hr">
                        <input placeholder="''' + get_lang(conn, 'why') + '''" name="why" type="text">
                        <hr class="main_hr">
                        <select name="ban_option">
                            <option value="">''' + get_lang(conn, 'default') + '''</option>
                            <option value="login_able">''' + get_lang(conn, 'login_able') + '''</option>
                            <option value="edit_request_able">''' + get_lang(conn, 'edit_request_able') + '''</option>
                            <option value="release">''' + get_lang(conn, 'release') + '''</option>
                        </select>
                        <hr class="main_hr">
                        <button type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                ''',
                menu = [['manager', get_lang(conn, 'return')]]
            ))   