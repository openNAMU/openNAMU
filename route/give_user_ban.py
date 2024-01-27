from .tool.func import *

def give_user_ban(name = None, ban_type = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        
        if ban_check(ip = ip, tool = 'login') == 1:
            if ip_or_user(ip) == 1 or admin_check('all', None, ip) == 0:
                return re_error('/ban')
        else:
            if admin_check(1, None, ip) != 1:
                return re_error('/error/3')

        if flask.request.method == 'POST':
            time_limit = flask.request.form.get('date', '')
            if re.search(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$', time_limit):
                end = time_limit + ' 00:00:00'
            else:
                end = '0'
            
            regex_get = flask.request.form.get('regex', None)
            login = flask.request.form.get('login', '')
            why = flask.request.form.get('why', '')

            release = flask.request.form.get('release', '')

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
                        return re_error('/error/23')
                else:
                    type_d = None

                if type_d:
                    if admin_check(None, 'ban' + (' ' + type_d if type_d else '') + ' (' + name + ')') != 1:
                        return re_error('/error/3')
                else:
                    if name == ip:
                        if admin_check('all', 'ban (' + name + ')') != 1:
                            return re_error('/error/3')
                    else:
                        if admin_check(1, 'ban (' + name + ')') != 1:
                            return re_error('/error/3')

                ban_insert(
                    name,
                    end,
                    why,
                    login,
                    ip_check(),
                    type_d,
                    1 if release != '' else 0
                )

            return redirect('/block_log')
        else:
            if ban_type == 'multiple':
                main_name = load_lang('multiple_ban')
                n_name = '<textarea class="opennamu_textarea_500" placeholder="' + load_lang('name_or_ip_or_regex_multiple') + '" name="name"></textarea><hr class="main_hr">'
            else:
                main_name = load_lang('ban')
                n_name = '<input placeholder="' + load_lang('name_or_ip_or_regex') + '" value="' + (name if name else '') + '" name="name"><hr class="main_hr">'

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

            return easy_minify(flask.render_template(skin_check(),
                imp = [main_name, wiki_set(), wiki_custom(), wiki_css([now, 0])],
                data = info_data + '''
                    <form method="post" ''' + action + '''>
                        ''' + n_name + '''
                        <input type="checkbox" name="regex" ''' + ('checked' if ban_type == 'regex' else '') + '> ' + load_lang('regex') + '''
                        <hr class="main_hr">
                        <input type="date" value="''' + date_value + '''" name="date" pattern="\\d{4}-\\d{2}-\\d{2}">
                        <hr class="main_hr">
                        <input placeholder="''' + load_lang('why') + '''" name="why" type="text">
                        <hr class="main_hr">
                        <input type="checkbox" name="login"> ''' + load_lang('login_able') + '''
                        <hr class="main_hr">
                        <input type="checkbox" name="release"> ''' + load_lang('release') + '''
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['manager', load_lang('return')]]
            ))   