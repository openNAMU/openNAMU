from .tool.func import *

async def user_info(name = ''):
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
                tool_menu += '<li><a class="opennamu_not_exist_link" href="/alarm">' + await get_lang('alarm') + ' (' + str(count[0][0]) + ')</a></li>'
            else:
                tool_menu += '<li><a href="/alarm">' + await get_lang('alarm') + '</a></li>'
    
            if ip_or_user(ip) == 0:
                login_menu += '''
                    <li><a href="/logout">''' + await get_lang('logout') + '''</a></li>
                    <li><a href="/change">''' + await get_lang('user_setting') + '''</a></li>
                '''
    
                tool_menu += '<li><a href="/watch_list">' + await get_lang('watchlist') + '</a></li>'
                tool_menu += '<li><a href="/star_doc">' + await get_lang('star_doc') + '</a></li>'
                tool_menu += '<li><a href="/challenge">' + await get_lang('challenge_and_level_manage') + '</a></li>'
                tool_menu += '<li><a href="/acl/user:' + url_pas(ip) + '">' + await get_lang('user_document_acl') + '</a></li>'
            else:
                login_menu += '''
                    <li><a href="/login">''' + await get_lang('login') + '''</a></li>
                    <li><a href="/register">''' + await get_lang('register') + '''</a></li>
                    <li><a href="/change">''' + await get_lang('user_setting') + '''</a></li>
                    <li><a href="/login/find">''' + await get_lang('password_search') + '''</a></li>
                '''
                
            login_menu = '<h2>' + await get_lang('login') + '</h2><ul>' + login_menu + '</ul>'
            tool_menu = '<h2>' + await get_lang('tool') + '</h2><ul>' + tool_menu + '</ul>'
    
        if await acl_check(tool = 'ban_auth') != 1:
            curs.execute(db_change("select block from rb where block = ? and ongoing = '1'"), [ip])
            ban_name = await get_lang('release') if curs.fetchall() else await get_lang('ban')
            
            admin_menu = '''
                <h2>''' + await get_lang('admin') + '''</h2>
                <ul>
                    <li><a href="/auth/ban/''' + url_pas(ip) + '''">''' + ban_name + '''</a></li>
                    <li><a href="/list/user/check_submit/''' + url_pas(ip) + '''">''' + await get_lang('check') + '''</a></li>
                </ul>
            '''
        else:
            admin_menu = ''
                
        return await render_template(
            await get_lang('user_tool'),
            '''
                <h2>''' + await get_lang('state') + '''</h2>
                <div id="opennamu_get_user_info">''' + html.escape(ip) + '''</div>
                ''' + login_menu + '''
                ''' + tool_menu + '''
                <h2>''' + await get_lang('other') + '''</h2>
                <ul>
                    <li><a href="/record/''' + url_pas(ip) + '''">''' + await get_lang('edit_record') + '''</a></li>
                    <li><a href="/record/topic/''' + url_pas(ip) + '''">''' + await get_lang('discussion_record') + '''</a></li>
                    <li><a href="/record/bbs/''' + url_pas(ip) + '''">''' + await get_lang('bbs_record') + '''</a></li>
                    <li><a href="/record/bbs_comment/''' + url_pas(ip) + '''">''' + await get_lang('bbs_comment_record') + '''</a></li>
                    <li><a href="/topic/user:''' + url_pas(ip) + '''">''' + await get_lang('user_discussion') + '''</a></li>
                    <li><a href="/count/''' + url_pas(ip) + '''">''' + await get_lang('count') + '''</a></li>
                </ul>
                ''' + admin_menu + '''
            ''',
            0,
            [['other', await get_lang('other_tool')]]
        )
