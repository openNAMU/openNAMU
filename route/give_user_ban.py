from .tool.func import *

async def give_user_ban(name = None, ban_type = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        
        if (await ban_check(ip = ip, tool = 'login'))[0] == 1:
            if ip_or_user(ip) == 1 or await acl_check(tool = 'all_admin_auth', ip = ip) != 0:
                return await re_error(conn, 0)
        else:
            if await acl_check(tool = 'ban_auth', ip = ip) == 1:
                return await re_error(conn, 3)

        if flask.request.method == 'POST':
            end = '0'

            date_select = flask.request.form.get('date_type', 'days')
            if date_select == 'date': 
                time_limit = flask.request.form.get('date', '')
                if re.search(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$', time_limit):
                    end = time_limit + ' 00:00:00'
            else:
                time_limit = int(number_check(flask.request.form.get('date_days', '1')))

                time = datetime.datetime.now()
                plus = datetime.timedelta(days = time_limit)
                end = (time + plus).strftime("%Y-%m-%d %H:%M:%S")
            
            regex_get = flask.request.form.get('do_ban_type', '')
            why = flask.request.form.get('why', '')

            release = ''
            login = ''
            
            ban_option = flask.request.form.get('ban_option', '')
            if ban_option == 'login_able_and_regsiter_disable':
                login = 'O'
            elif ban_option == 'login_able':
                login = 'L'
            elif ban_option == 'edit_request_able':
                login = 'E'
            elif ban_option == 'completely_ban':
                login = 'A'
            elif ban_option == 'dont_come_this_site':
                login = 'D'
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
                if regex_get == 'regex':
                    type_d = 'regex'

                    try:
                        re.compile(name)
                    except:
                        return await re_error(conn, 23)
                elif regex_get == 'cidr':
                    type_d = 'cidr'

                    try:
                        ipaddress.IPv4Network(name, False)
                    except:
                        try:
                            ipaddress.IPv6Network(name, False)
                        except:
                            return await re_error(conn, 45)
                elif regex_get == 'private':
                    type_d = 'private'

                    if await acl_check(tool = 'owner_auth', ip = ip) == 1:
                        return await re_error(conn, 0)
                else:
                    type_d = None

                if regex_get != 'private':
                    if name == ip:
                        if await acl_check(tool = 'all_admin_auth', memo = 'ban (' + name + ')') == 1:
                            return await re_error(conn, 3)
                    else:
                        if await acl_check(tool = 'ban_auth', memo = 'ban (' + name + ')') == 1:
                            return await re_error(conn, 3)

                ban_insert(conn, 
                    name,
                    end,
                    why,
                    login,
                    ip_check(),
                    type_d,
                    1 if release != '' else 0
                )

            return redirect(conn, '/recent_block')
        else:
            if ban_type == 'multiple':
                main_name = await get_lang('multiple_ban')
                n_name = '<textarea class="opennamu_textarea_500" placeholder="' + await get_lang('name_or_ip_or_regex_or_cidr_multiple') + '" name="name"></textarea><hr class="main_hr">'
            else:
                main_name = await get_lang('ban')
                n_name = '<input placeholder="' + await get_lang('name_or_ip_or_regex_or_cidr') + '" value="' + (name if name else '') + '" name="name"><hr class="main_hr">'

            now = 0
            
            if ban_type == 'multiple':
                action = 'action="/auth/ban/multiple"'
            else:
                action = 'action="/auth/ban"'
                
            date_value = ''
            info_data = ''
            if name:
                curs.execute(db_change("select end from rb where block = ? and ongoing = '1'"), [name])
                db_data = curs.fetchall()
                if db_data and db_data[0][0] != '':
                    date_value = db_data[0][0].split()[0]

                if ban_type == '':
                    info_data = '<div id="opennamu_get_user_info">' + html.escape(name) + '</div>'

            owner_option = ''
            if await acl_check(tool = 'owner_auth', ip = ip) != 1:
                owner_option = '<option value="private" ' + ('selected' if ban_type == 'private' else '') + '>' + await get_lang('private') + '</option>'

            return easy_minify(flask.render_template(await skin_check(conn),
                imp = [main_name, await wiki_set(), await wiki_custom(conn), wiki_css([now, 0])],
                data = info_data + '''
                    <form method="post" ''' + action + '''>
                        <h2>''' + await get_lang('method') + '''</h2>
                        ''' + n_name + '''
        
                        <select name="do_ban_type">
                            <option value="normal">''' + await get_lang('normal') + '''</option>
                            <option value="regex" ''' + ('selected' if ban_type == 'regex' else '') + '>' + await get_lang('regex') + '''</option>
                            <option value="cidr" ''' + ('selected' if ban_type == 'cidr' else '') + '>' + await get_lang('cidr') + '''</option>
                            ''' + owner_option + '''
                        </select>
                        <hr class="main_hr">
        
                        <select name="ban_option">
                            <option value="">''' + await get_lang('default') + '''</option>
                            <option value="login_able">''' + await get_lang('login_able') + '''</option>
                            <option value="login_able_and_regsiter_disable">''' + await get_lang('login_able_and_regsiter_disable') + '''</option>
                            <option value="edit_request_able">''' + await get_lang('edit_request_able') + '''</option>
                            <option value="completely_ban">''' + await get_lang('completely_ban') + '''</option>
                            <option value="dont_come_this_site">''' + await get_lang('dont_come_this_site') + '''</option>
                            <option value="release">''' + await get_lang('release') + '''</option>
                        </select>
        
                        <h2>''' + await get_lang('date') + '''</h2>
                        <select name="date_type">
                            <option value="date">''' + await get_lang('date') + '''</option>
                            <option value="days">''' + await get_lang('day') + '''</option>
                        </select>
                        <hr class="main_hr">
        
                        <span>''' + await get_lang('day') + '''</span>
                        <hr class="main_hr">
                        <input name="date_days">
                        <hr class="main_hr">

                        <span>''' + await get_lang('date') + '''</span>
                        <hr class="main_hr">
                        <input type="date" value="''' + date_value + '''" name="date" pattern="\\d{4}-\\d{2}-\\d{2}">
        
                        <h2>''' + await get_lang('other') + '''</h2>
                        <input placeholder="''' + await get_lang('why') + '''" name="why" type="text">
                        <hr class="main_hr">
        
                        <button type="submit">''' + await get_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['manager', await get_lang('return')]]
            ))   