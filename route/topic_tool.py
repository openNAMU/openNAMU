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
                t_state = load_lang('topic_stop')
            elif close_data[0][0] == 'O':
                t_state = load_lang('topic_close')
            else:
                t_state = load_lang('topic_normal')
                
            if close_data[0][1] == 'O':
                t_state += ' (' + load_lang('topic_agree') + ')'
        else:
            t_state = load_lang('topic_normal')

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

        if admin_check(3) == 1:
            data = '''
                <h2>''' + load_lang('admin_tool') + '''</h2>
                <ul class="opennamu_ul">
                    <li><a href="/thread/''' + topic_num + '/setting">' + load_lang('topic_setting') + '''</a></li>
                    <li><a href="/thread/''' + topic_num + '/acl">' + load_lang('topic_acl_setting') + '''</a></li>
                </ul>
            '''
        data += '''
            <h2>''' + load_lang('tool') + '''</h2>
            <ul class="opennamu_ul">
                <li>''' + load_lang('topic_state') + ''' : ''' + t_state + '''</li>
                <li>''' + load_lang('topic_acl') + ''' : <a href="/acl/TEST#exp">''' + acl_state + '''</a></li>
                <li>''' + load_lang('topic_view_acl') + ''' : <a href="/acl/TEST#exp">''' + acl_view_state + '''</a></li>
            </ul>
        '''

        if admin_check(None) == 1:
            data += '''
                <h2>''' + load_lang('owner') + '''</h2>
                <ul class="opennamu_ul">
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