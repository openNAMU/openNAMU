from .tool.func import *

def topic_admin_2(conn, name, sub, num):
    curs = conn.cursor()

    curs.execute("select block, ip, date from topic where title = ? and sub = ? and id = ?", [name, sub, str(num)])
    data = curs.fetchall()
    if not data:
        return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub))

    ban = ''

    if admin_check(3) == 1:
        ban +=  '''
            </ul>
            <br>
            <h2>''' + load_lang('admin_tool') + '''</h2>
            <ul>
        '''
        is_ban = '<li><a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/b/' + str(num) + '">'

        if data[0][0] == 'O':
            is_ban += load_lang('hide_release')
        else:
            is_ban += load_lang('hide')
        
        is_ban +=   '''
                        </a>
                    </li>
                    <li>
                        <a href="/topic/''' + url_pas(name) + '/sub/' + url_pas(sub) + '/notice/' + str(num) + '''">
                    '''

        curs.execute("select id from topic where title = ? and sub = ? and id = ? and top = 'O'", [name, sub, str(num)])
        if curs.fetchall():
            is_ban += load_lang('notice_release')
        else:
            is_ban += load_lang('notice') + ''
        
        is_ban += '</a></li></ul>'
        ban += '<li><a href="/ban/' + url_pas(data[0][1]) + '">'

        curs.execute("select end from ban where block = ?", [data[0][1]])
        if curs.fetchall():
            ban += load_lang('ban_release')
        else:
            ban += load_lang('ban')
        
        ban += '</a></li>' + is_ban

    ban +=  '''
            </ul>
            <br>
            <h2>''' + load_lang('other_tool') + '''</h2>
            <ul>
                <li>
                    <a href="/topic/''' + url_pas(name) + '/sub/' + url_pas(sub) + '/raw/' + str(num) + '''">''' + load_lang('raw') + '''</a>
                </li>
            '''
    ban = '<li>' + load_lang('time') + ' : ' + data[0][2] + '</li>' + ban
    
    if ip_or_user(data[0][1]) == 1:
        ban = '<li>' + load_lang('writer') + ' : ' + data[0][1] + ' <a href="/record/' + url_pas(data[0][1]) + '">(' + load_lang('record') + ')</a></li>' + ban
    else:
        ban =   '''
                <li>
                    ''' + load_lang('writer') + ' : <a href="/w/user:' + data[0][1] + '">' + data[0][1] + '</a> <a href="/record/' + url_pas(data[0][1]) + '">(' + load_lang('record') + ''')</a>
                </li>
                ''' + ban

    ban = '<h2>' + load_lang('state') + '</h2><ul>' + ban

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('discussion_tool'), wiki_set(), custom(), other2([' (' + str(num) + ')', 0])],
        data = ban,
        menu = [['topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num), load_lang('return')]]
    ))
