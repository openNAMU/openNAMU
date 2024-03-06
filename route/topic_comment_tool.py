from .tool.func import *

def topic_comment_tool(topic_num = 1, num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        num = str(num)
        topic_num = str(topic_num)
        
        if acl_check('', 'topic_view', topic_num) == 1:
            return re_error('/ban')

        curs.execute(db_change("select block, ip, date from topic where code = ? and id = ?"), [topic_num, num])
        data = curs.fetchall()
        if not data:
            return redirect('/thread/' + topic_num)

        ban = '''
            <h2>''' + load_lang('state') + '''</h2>
            <ul class="opennamu_ul">
                <li>''' + load_lang('writer') + ' : ''' + ip_pas(data[0][1]) + '''</li>
                <li>''' + load_lang('time') + ' : ' + data[0][2] + '''</li>
            </ul>
            <h2>''' + load_lang('other_tool') + '''</h2>
            <ul class="opennamu_ul">
                <li>
                    <a href="/thread/''' + topic_num + '/comment/' + num + '''/raw">''' + load_lang('raw') + '''</a>
                </li>
            </ul>
        '''

        if admin_check(3) == 1:
            ban += '''
                <h2>''' + load_lang('admin_tool') + '''</h2>
                <ul class="opennamu_ul">
                    <li>
                        <a href="/auth/give/ban/''' + url_pas(data[0][1]) + '''">
                            ''' + (load_lang('ban') + ' | ' + load_lang('release')) + '''
                        </a>
                    </li>
                    <li>
                        <a href="/thread/''' + topic_num + '''/comment/''' + num + '''/blind">
                            ''' + (load_lang('hide') + ' | ' + load_lang('hide_release')) + '''
                        </a>
                    </li>
                    <li>
                        <a href="/thread/''' + topic_num + '''/comment/''' + num + '''/notice">
                            ''' + (load_lang('pinned') + ' | ' + load_lang('pinned_release')) + '''
                        </a>
                    </li>
                    <li>
                        <a href="/thread/''' + topic_num + '''/comment/''' + num + '''/delete">
                            ''' + load_lang('delete') + '''
                        </a>
                </ul>
            '''

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('discussion_tool'), wiki_set(), wiki_custom(), wiki_css(['(#' + num + ')', 0])],
            data = ban,
            menu = [['thread/' + topic_num + '#' + num, load_lang('return')]]
        ))