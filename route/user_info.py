from .tool.func import *

def user_info_2(conn):
    curs = conn.cursor()

    ip = ip_check()

    curs.execute(db_change("select count(name) from alarm where name = ?"), [ip])
    count = curs.fetchall()
    if count:
        plus2 = '<li><a href="/alarm">' + load_lang('alarm') + ' (' + str(count[0][0]) + ')</a></li>'
    else:
        plus2 = '<li><a href="/alarm">' + load_lang('alarm') + ' (0)</a></li>'

    if ip_or_user(ip) == 0:  
        plus = '''
            <li><a href="/logout">''' + load_lang('logout') + '''</a></li>
            <li><a href="/change">''' + load_lang('user_setting') + '''</a></li>
        '''

        plus2 += '<li><a href="/watch_list">' + load_lang('watchlist') + '</a></li>'
        plus3 = '<li><a href="/acl/user:' + url_pas(ip) + '">' + load_lang('user_document_acl') + '</a></li>'
    else:        
        plus = '''
            <li><a href="/login">''' + load_lang('login') + '''</a></li>
            <li><a href="/register">''' + load_lang('register') + '''</a></li>
        '''
        plus3 = ''

        curs.execute(db_change("select data from other where name = 'email_have'"))
        test = curs.fetchall()
        if test and test[0][0] != '':
            plus += '<li><a href="/pass_find">' + load_lang('password_search') + '</a></li>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('user_tool'), wiki_set(), custom(), other2([0, 0])],
        data = '''
            <h2>''' + load_lang('state') + '''</h2>
            <div id="get_user_info"></div>
            <script>load_user_info("''' + ip + '''");</script>
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