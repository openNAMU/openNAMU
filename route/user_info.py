from .tool.func import *

def user_info_2(conn):
    curs = conn.cursor()

    ip = ip_check()
    
    curs.execute("select acl from user where id = ?", [ip])
    data = curs.fetchall()
    if ban_check() == 0:
        if data:
            if data[0][0] != 'user':
                acl = data[0][0]
            else:
                acl = load_lang('member')
        else:
            acl = load_lang('normal')
    else:
        acl = load_lang('blocked')

        match = re.search("^([0-9]{1,3}\.[0-9]{1,3})", ip)
        if match:
            match = match.groups()[0]
        else:
            match = '-'

        curs.execute("select end, login, band from ban where block = ? or block = ?", [ip, match])
        block_data = curs.fetchall()
        if block_data:
            if block_data[0][0] != '':
                acl += ' (' + load_lang('period') + ' : ' + block_data[0][0] + ')'
            else:
                acl += ' (' + load_lang('limitless') + ')'        

            if block_data[0][1] != '':
                acl += ' (' + load_lang('login_able') + ')'

            if block_data[0][2] == 'O':
                acl += ' (' + load_lang('band_blocked') + ')'
            
    if custom()[2] != 0:
        ip_user = '<a href="/w/user:' + ip + '">' + ip + '</a>'
        
        plus =  '''
                <li><a href="/logout">''' + load_lang('logout') + '''</a></li>
                <li><a href="/change">''' + load_lang('user_setting') + '''</a></li>
                '''
        
        curs.execute('select name from alarm where name = ? limit 1', [ip_check()])
        if curs.fetchall():
            plus2 = '<li><a href="/alarm">' + load_lang('alarm') + ' (O)</a></li>'
        else:
            plus2 = '<li><a href="/alarm">' + load_lang('alarm') + '</a></li>'

        plus2 += '<li><a href="/watch_list">' + load_lang('watchlist') + '</a></li>'
        plus3 = '<li><a href="/acl/user:' + url_pas(ip) + '">' + load_lang('user_document_acl') + '</a></li>'
    else:
        ip_user = ip
        
        plus =  '''
                <li><a href="/login">''' + load_lang('login') + '''</a></li>
                <li><a href="/register">''' + load_lang('register') + '''</a></li>
                '''
        plus2 = ''
        plus3 = ''

        curs.execute("select data from other where name = 'email_have'")
        test = curs.fetchall()
        if test and test[0][0] != '':
            plus += '<li><a href="/pass_find">' + load_lang('password_search') + '</a></li>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('user') + ' ' + load_lang('tool'), wiki_set(), custom(), other2([0, 0])],
        data =  '''
                <h2>''' + load_lang('state') + '''</h2>
                <ul>
                    <li>''' + ip_user + ''' <a href="/record/''' + url_pas(ip) + '''">(''' + load_lang('record') + ''')</a></li>
                    <li>''' + load_lang('state') + ''' : ''' + acl + '''</li>
                </ul>
                <br>
                <h2>''' + load_lang('login') + '''</h2>
                <ul>
                    ''' + plus + '''
                </ul>
                <br>
                <h2>''' + load_lang('tool') + '''</h2>
                <ul>
                    ''' + plus3 + '''
                    <li><a href="/custom_head">''' + load_lang('user_head') + '''</a></li>
                </ul>
                <br>
                <h2>''' + load_lang('other') + '''</h2>
                <ul>
                ''' + plus2 + '''
                <li>
                    <a href="/count">''' + load_lang('count') + '''</a>
                </li>
                </ul>
                ''',
        menu = 0
    ))