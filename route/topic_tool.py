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
                t_state = await get_lang('topic_stop')
            elif close_data[0][0] == 'O':
                t_state = await get_lang('topic_close')
            else:
                t_state = await get_lang('topic_normal')
                
            if close_data[0][1] == 'O':
                t_state += ' (' + await get_lang('topic_agree') + ')'
        else:
            t_state = await get_lang('topic_normal')

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
                <h2>''' + await get_lang('admin_tool') + '''</h2>
                <ul>
                    <li><a href="/thread/''' + topic_num + '/setting">' + await get_lang('topic_setting') + '''</a></li>
                    <li><a href="/thread/''' + topic_num + '/acl">' + await get_lang('topic_acl_setting') + '''</a></li>
                </ul>
            '''
        data += '''
            <h2>''' + await get_lang('tool') + '''</h2>
            <ul>
                <li>''' + await get_lang('topic_state') + ''' : ''' + t_state + '''</li>
                <li>''' + await get_lang('topic_acl') + ''' : <a href="/acl/TEST#exp">''' + acl_state + '''</a></li>
                <li>''' + await get_lang('topic_view_acl') + ''' : <a href="/acl/TEST#exp">''' + acl_view_state + '''</a></li>
            </ul>
        '''

        if await acl_check(tool = 'owner_auth') != 1:
            data += '''
                <h2>''' + await get_lang('owner') + '''</h2>
                <ul>
                    <li>
                        <a href="/thread/''' + topic_num + '''/delete">
                            ''' + await get_lang('topic_delete') + '''
                        </a>
                    </li>
                    <li>
                        <a href="/thread/''' + topic_num + '''/change">
                            ''' + await get_lang('topic_name_change') + '''
                        </a>
                    </li>
                </ul>
            '''

        return await render_template(
            await get_lang('topic_tool'),
            data,
            0,
            [['thread/' + topic_num, await get_lang('return')]]
        )
