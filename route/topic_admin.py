from .tool.func import *

def topic_admin_2(conn, topic_num, num):
    curs = conn.cursor()

    topic_change_data = topic_change(topic_num)
    name = topic_change_data[0]
    sub = topic_change_data[1]

    curs.execute(db_change("select block, ip, date from topic where title = ? and sub = ? and id = ?"), [name, sub, str(num)])
    data = curs.fetchall()
    if not data:
        return redirect('/thread/' + str(topic_num))

    ban = '''
        <h2>''' + load_lang('state') + '''</h2>
        <ul>
            <li>''' + load_lang('writer') + ' : ''' + ip_pas(data[0][1]) + '''</li>
            <li>''' + load_lang('time') + ' : ' + data[0][2] + '''</li>
        </ul>
        <br>
        <h2>''' + load_lang('other_tool') + '''</h2>
        <ul>
            <li>
                <a href="/thread/''' + str(topic_num) + '/raw/' + str(num) + '''">''' + load_lang('raw') + '''</a>
            </li>
        </ul>
    '''

    if admin_check(3) == 1:
        curs.execute(db_change("select id from topic where title = ? and sub = ? and id = ? and top = 'O'"), [name, sub, str(num)])
        top_topic_d = curs.fetchall()

        curs.execute(db_change("select end from ban where block = ?"), [data[0][1]])
        user_ban_d = curs.fetchall()

        ban += '''
            <br>
            <h2>''' + load_lang('admin_tool') + '''</h2>
            <ul>
                <li>
                    <a href="/ban/''' + url_pas(data[0][1]) + '''">
                        ''' + (load_lang('ban_release') if user_ban_d else load_lang('ban')) + '''
                    </a>
                </li>
                <li>
                    <a href="/thread/''' + str(topic_num) + '/b/' + str(num) + '''">
                        ''' + (load_lang('hide_release') if data[0][0] == 'O' else load_lang('hide')) + '''
                    </a>
                </li>
                <li>
                    <a href="/thread/''' + str(topic_num) + '/notice/' + str(num) + '''">
                        ''' + (load_lang('pinned_release') if top_topic_d else load_lang('pinned')) + '''
                    </a>
                </li>
            </ul>
        '''

    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('discussion_tool'), wiki_set(), custom(), other2([' (#' + str(num) + ')', 0])],
        data = ban,
        menu = [['thread/' + str(topic_num) + '#' + str(num), load_lang('return')]]
    ))
