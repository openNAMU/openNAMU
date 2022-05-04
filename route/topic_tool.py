from .tool.func import *

def topic_tool(topic_num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        data = ''
        topic_num = str(topic_num)

        curs.execute(db_change("select stop, agree from rd where code = ?"), [topic_num])
        close_data = curs.fetchall()
        if close_data:
            if close_data[0][0] == 'S':
                t_state = 'Stop'
            elif close_data[0][0] == 'O':
                t_state = 'Close'
            else:
                t_state = 'Normal'
        else:
            t_state = 'Normal'

        curs.execute(db_change("select acl from rd where code = ?"), [topic_num])
        topic_acl_get = curs.fetchall()

        if admin_check(3) == 1:
            data = '''
                <h2>''' + load_lang('admin_tool') + '''</h2>
                <ul class="inside_ul">
                    <li><a href="/thread/''' + topic_num + '/setting">' + load_lang('topic_setting') + '''</a></li>
                    <li><a href="/thread/''' + topic_num + '/acl">' + load_lang('topic_acl_setting') + '''</a></li>
                </ul>
            '''
        data += '''
            <h2>''' + load_lang('tool') + '''</h2>
            <ul class="inside_ul">
                <li>''' + load_lang('topic_state') + ''' : ''' + t_state + '' + (' (Agree)' if close_data and (close_data[0][1] == 'O') else '') + '''</li>
                <li>''' + load_lang('topic_acl') + ''' : <a href="/acl/TEST#exp">''' + ('Normal' if not topic_acl_get or (topic_acl_get[0][0] == '') else topic_acl_get[0][0]) + '''</a></li>
            </ul>
        '''

        if admin_check(None) == 1:
            data += '''
                <h2>''' + load_lang('owner') + '''</h2>
                <ul class="inside_ul">
                    <li>
                        <a href="/thread/''' + topic_num + '''/delete">
                            ''' + load_lang('topic_delete') + '''
                        </a>
                    </li>
                    <li>
                        <a href="/thread/''' + topic_num + '''/change">
                            ''' + load_lang('topic_name_change') + '''
                        </a>
                    </li>
                </ul>
            '''

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('topic_tool'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = data,
            menu = [['thread/' + topic_num, load_lang('return')]]
        ))