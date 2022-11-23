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
            end = flask.request.form.get('second', '0')
            end = end if end else '0'
            
            regex_get = flask.request.form.get('regex', None)
            login = flask.request.form.get('login', '')
            why = flask.request.form.get('why', '')

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
                    type_d
                )

            return redirect('/block_log')
        else:
            curs.execute(db_change("select end, why from rb where block = ? and ongoing = '1' and band = ?"), [name, ban_type])
            end = curs.fetchall()
            if end:
                main_name = name
                b_now = load_lang('release')
                now = '(' + b_now + ')'
                
                if ban_type == 'regex':
                    action = 'action="/auth/give/ban_regex/' + url_pas(name) + '"'
                else:
                    action = 'action="/auth/give/ban/' + url_pas(name) + '"'

                if end[0][0] == '':
                    data = '<ul class="opennamu_ul"><li>' + load_lang('limitless') + '</li>'
                else:
                    data = '<ul class="opennamu_ul"><li>' + load_lang('period') + ' : ' + end[0][0] + '</li>'

                curs.execute(db_change("select block from rb where block = ? and login = 'O' and ongoing = '1'"), [name])
                if curs.fetchall():
                    data += '<li>' + load_lang('login_able') + '</li>'

                if end[0][1] != '':
                    data += '<li>' + load_lang('why') + ' : ' + end[0][1] + '</li></ul><hr class="main_hr">'
                else:
                    data += '</ul><hr class="main_hr">'
            else:
                if ban_type == 'multiple':
                    main_name = load_lang('multiple_ban')
                    n_name = '<textarea class="opennamu_textarea_500" placeholder="' + load_lang('name_or_ip_or_regex_multiple') + '" name="name"></textarea><hr class="main_hr">'
                else:
                    main_name = load_lang('ban')
                    n_name = '<input placeholder="' + load_lang('name_or_ip_or_regex') + '" value="' + (name if name else '') + '" name="name"><hr class="main_hr">'
                
                regex = '<input type="checkbox" name="regex" ' + ('checked' if ban_type == 'regex' else '') + '> ' + load_lang('regex') + '<hr class="main_hr">'
                plus = '<input type="checkbox" name="login"> ' + load_lang('login_able') + '<hr class="main_hr">'
                now = 0
                b_now = load_lang('ban')
                
                if ban_type == 'multiple':
                    action = 'action="/auth/give/ban_multiple"'
                else:
                    action = 'action="/auth/give/ban"'

                time_data = [
                    ['86400', load_lang('1_day')],
                    ['432000', load_lang('5_day')],
                    ['2592000', load_lang('30_day')],
                    ['15552000', load_lang('180_day')],
                    ['31104000', load_lang('360_day')],
                    ['0', load_lang('limitless')]
                ]
                insert_data = ''
                for i in time_data:
                    insert_data += '<a href="javascript:insert_v(\'second\', \'' + i[0] + '\')">(' + i[1] + ')</a> '

                data = n_name + '''
                    ''' + regex + '''
                    <script>function insert_v(name, data) { document.getElementById(name).value = data; }</script>''' + insert_data + '''
                    <hr class="main_hr">
                    <input placeholder="''' + load_lang('ban_period') + ''' (''' + load_lang('second') + ''')" name="second" id="second" type="text">
                    <hr class="main_hr">
                    <input placeholder="''' + load_lang('why') + '''" name="why" type="text">
                    <hr class="main_hr">
                ''' + plus

            return easy_minify(flask.render_template(skin_check(),
                imp = [main_name, wiki_set(), wiki_custom(), wiki_css([now, 0])],
                data = '''
                    <form method="post" ''' + action + '''>
                        ''' + data + '''
                        <button type="submit">''' + b_now + '''</button>
                    </form>
                ''',
                menu = [['manager', load_lang('return')]]
            ))   