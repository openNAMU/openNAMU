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
            <h2>''' + get_lang(conn, 'state') + '''</h2>
            <ul>
                <li>''' + get_lang(conn, 'writer') + ' : ''' + await ip_pas(data[0][1]) + '''</li>
                <li>''' + get_lang(conn, 'time') + ' : ' + data[0][2] + '''</li>
            </ul>
            <h2>''' + get_lang(conn, 'other_tool') + '''</h2>
            <ul>
                <li>
                    <a href="/thread/''' + topic_num + '/comment/' + num + '''/raw">''' + get_lang(conn, 'raw') + '''</a>
                </li>
            </ul>
        '''

        if await acl_check(tool = 'toron_auth') != 1:
            ban += '''
                <h2>''' + get_lang(conn, 'admin_tool') + '''</h2>
                <ul>
                    <li>
                        <a href="/auth/ban/''' + url_pas(data[0][1]) + '''">
                            ''' + (get_lang(conn, 'ban') + ' | ' + get_lang(conn, 'release')) + '''
                        </a>
                    </li>
                    <li>
                        <a href="/thread/''' + topic_num + '''/comment/''' + num + '''/blind">
                            ''' + (get_lang(conn, 'hide') + ' | ' + get_lang(conn, 'hide_release')) + '''
                        </a>
                    </li>
                    <li>
                        <a href="/thread/''' + topic_num + '''/comment/''' + num + '''/notice">
                            ''' + (get_lang(conn, 'pinned') + ' | ' + get_lang(conn, 'pinned_release')) + '''
                        </a>
                    </li>
                    <li>
                        <a href="/thread/''' + topic_num + '''/comment/''' + num + '''/delete">
                            ''' + get_lang(conn, 'delete') + '''
                        </a>
                </ul>
            '''

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'discussion_tool'), await wiki_set(), await wiki_custom(conn), wiki_css(['(#' + num + ')', 0])],
            data = ban,
            menu = [['thread/' + topic_num + '#' + num, get_lang(conn, 'return')]]
        ))