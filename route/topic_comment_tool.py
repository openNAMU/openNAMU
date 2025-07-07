from .tool.func import *

async def topic_comment_tool(topic_num = 1, num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        num = str(num)
        topic_num = str(topic_num)
        
        if await acl_check('', 'topic_view', topic_num) == 1:
            return await re_error(conn, 0)

        curs.execute(db_change("select block, ip, date from topic where code = ? and id = ?"), [topic_num, num])
        data = curs.fetchall()
        if not data:
            return redirect(conn, '/thread/' + topic_num)

        ban = '''
            <h2>''' + await get_lang('state') + '''</h2>
            <ul>
                <li>''' + await get_lang('writer') + ' : ''' + await ip_pas(data[0][1]) + '''</li>
                <li>''' + await get_lang('time') + ' : ' + data[0][2] + '''</li>
            </ul>
            <h2>''' + await get_lang('other_tool') + '''</h2>
            <ul>
                <li>
                    <a href="/thread/''' + topic_num + '/comment/' + num + '''/raw">''' + await get_lang('raw') + '''</a>
                </li>
            </ul>
        '''

        if await acl_check(tool = 'toron_auth') != 1:
            ban += '''
                <h2>''' + await get_lang('admin_tool') + '''</h2>
                <ul>
                    <li>
                        <a href="/auth/ban/''' + url_pas(data[0][1]) + '''">
                            ''' + (await get_lang('ban') + ' | ' + await get_lang('release')) + '''
                        </a>
                    </li>
                    <li>
                        <a href="/thread/''' + topic_num + '''/comment/''' + num + '''/blind">
                            ''' + (await get_lang('hide') + ' | ' + await get_lang('hide_release')) + '''
                        </a>
                    </li>
                    <li>
                        <a href="/thread/''' + topic_num + '''/comment/''' + num + '''/notice">
                            ''' + (await get_lang('pinned') + ' | ' + await get_lang('pinned_release')) + '''
                        </a>
                    </li>
                    <li>
                        <a href="/thread/''' + topic_num + '''/comment/''' + num + '''/delete">
                            ''' + await get_lang('delete') + '''
                        </a>
                </ul>
            '''

        return easy_minify(flask.render_template(await skin_check(conn),
            imp = [await get_lang('discussion_tool'), await wiki_set(), await wiki_custom(conn), wiki_css(['(#' + num + ')', 0])],
            data = ban,
            menu = [['thread/' + topic_num + '#' + num, await get_lang('return')]]
        ))