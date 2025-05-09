from .tool.func import *

async def topic_tool(topic_num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        data = ''
        topic_num = str(topic_num)

        curs.execute(db_change("select stop, agree from rd where code = ?"), [topic_num])
        close_data = curs.fetchall()
        if close_data:
            if close_data[0][0] == 'S':
                t_state = get_lang(conn, 'topic_stop')
            elif close_data[0][0] == 'O':
                t_state = get_lang(conn, 'topic_close')
            else:
                t_state = get_lang(conn, 'topic_normal')
                
            if close_data[0][1] == 'O':
                t_state += ' (' + get_lang(conn, 'topic_agree') + ')'
        else:
            t_state = get_lang(conn, 'topic_normal')

        curs.execute(db_change("select acl from rd where code = ?"), [topic_num])
        db_data = curs.fetchall()
        if db_data:
            if db_data[0][0] == '':
                acl_state = 'normal'
            else:
                acl_state = db_data[0][0]
        else:
            acl_state = 'normal'
        
        curs.execute(db_change("select set_data from topic_set where thread_code = ? and set_name = 'thread_view_acl'"), [topic_num])
        db_data = curs.fetchall()
        if db_data:
            if db_data[0][0] == '':
                acl_view_state = 'normal'
            else:
                acl_view_state = db_data[0][0]
        else:
            acl_view_state = 'normal'

        if await acl_check(tool = 'toron_auth') != 1:
            data = '''
                <h2>''' + get_lang(conn, 'admin_tool') + '''</h2>
                <ul>
                    <li><a href="/thread/''' + topic_num + '/setting">' + get_lang(conn, 'topic_setting') + '''</a></li>
                    <li><a href="/thread/''' + topic_num + '/acl">' + get_lang(conn, 'topic_acl_setting') + '''</a></li>
                </ul>
            '''
        data += '''
            <h2>''' + get_lang(conn, 'tool') + '''</h2>
            <ul>
                <li>''' + get_lang(conn, 'topic_state') + ''' : ''' + t_state + '''</li>
                <li>''' + get_lang(conn, 'topic_acl') + ''' : <a href="/acl/TEST#exp">''' + acl_state + '''</a></li>
                <li>''' + get_lang(conn, 'topic_view_acl') + ''' : <a href="/acl/TEST#exp">''' + acl_view_state + '''</a></li>
            </ul>
        '''

        if await acl_check(tool = 'owner_auth') != 1:
            data += '''
                <h2>''' + get_lang(conn, 'owner') + '''</h2>
                <ul>
                    <li>
                        <a href="/thread/''' + topic_num + '''/delete">
                            ''' + get_lang(conn, 'topic_delete') + '''
                        </a>
                    </li>
                    <li>
                        <a href="/thread/''' + topic_num + '''/change">
                            ''' + get_lang(conn, 'topic_name_change') + '''
                        </a>
                    </li>
                </ul>
            '''

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'topic_tool'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
            data = data,
            menu = [['thread/' + topic_num, get_lang(conn, 'return')]]
        ))