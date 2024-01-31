from .tool.func import *

def user_info(name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()
    
        if name == '':
            ip = ip_check()
        else:
            ip = name
    
        login_menu = ''
        tool_menu = ''
        
        if name == '':
            curs.execute(db_change("select count(*) from user_notice where name = ? and readme = ''"), [ip])
            count = curs.fetchall()
            if count and count[0][0] != 0:
                tool_menu += '<li><a class="opennamu_not_exist_link" href="/alarm">' + load_lang('alarm') + ' (' + str(count[0][0]) + ')</a></li>'
            else:
                tool_menu += '<li><a href="/alarm">' + load_lang('alarm') + '</a></li>'
    
            if ip_or_user(ip) == 0:
                login_menu += '''
                    <li><a href="/logout">''' + load_lang('logout') + '''</a></li>
                    <li><a href="/change">''' + load_lang('user_setting') + '''</a></li>
                '''
    
                tool_menu += '<li><a href="/watch_list">' + load_lang('watchlist') + '</a></li>'
                tool_menu += '<li><a href="/star_doc">' + load_lang('star_doc') + '</a></li>'
                tool_menu += '<li><a href="/challenge">' + load_lang('challenge_and_level_manage') + '</a></li>'
                tool_menu += '<li><a href="/acl/user:' + url_pas(ip) + '">' + load_lang('user_document_acl') + '</a></li>'
            else:
                login_menu += '''
                    <li><a href="/login">''' + load_lang('login') + '''</a></li>
                    <li><a href="/register">''' + load_lang('register') + '''</a></li>
                    <li><a href="/change">''' + load_lang('user_setting') + '''</a></li>
                    <li><a href="/login/find">''' + load_lang('password_search') + '''</a></li>
                '''
                
            login_menu = '<h2>' + load_lang('login') + '</h2><ul class="opennamu_ul">' + login_menu + '</ul>'
            tool_menu = '<h2>' + load_lang('tool') + '</h2><ul class="opennamu_ul">' + tool_menu + '</ul>'
    
        if admin_check(1) == 1:
            curs.execute(db_change("select block from rb where block = ? and ongoing = '1'"), [ip])
            ban_name = load_lang('release') if curs.fetchall() else load_lang('ban')
            
            admin_menu = '''
                <h2>''' + load_lang('admin') + '''</h2>
                <ul class="opennamu_ul">
                    <li><a href="/auth/give/ban/''' + url_pas(ip) + '''">''' + ban_name + '''</a></li>
                    <li><a href="/list/user/check/''' + url_pas(ip) + '''">''' + load_lang('check') + '''</a></li>
                </ul>
            '''
        else:
            admin_menu = ''
                
        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('user_tool'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = '''
                <h2>''' + load_lang('state') + '''</h2>
                <div id="opennamu_get_user_info">''' + html.escape(ip) + '''</div>
                ''' + login_menu + '''
                ''' + tool_menu + '''
                <h2>''' + load_lang('other') + '''</h2>
                <ul class="opennamu_ul">
                    <li><a href="/record/''' + url_pas(ip) + '''">''' + load_lang('edit_record') + '''</a></li>
                    <li><a href="/record/topic/''' + url_pas(ip) + '''">''' + load_lang('discussion_record') + '''</a></li>
                    <li><a href="/record/bbs/''' + url_pas(ip) + '''">''' + load_lang('bbs_record') + '''</a></li>
                    <li><a href="/record/bbs_comment/''' + url_pas(ip) + '''">''' + load_lang('bbs_comment_record') + '''</a></li>
                    <li><a href="/topic/user:''' + url_pas(ip) + '''">''' + load_lang('user_discussion') + '''</a></li>
                    <li><a href="/count/''' + url_pas(ip) + '''">''' + load_lang('count') + '''</a></li>
                </ul>
                ''' + admin_menu + '''
            ''',
            menu = [['other', load_lang('other_tool')]]
        ))