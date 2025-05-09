from .tool.func import *

def recent_change_send_render(data):
    def send_render_href_replace(match):
        match = match.group(1)
        data_unescape = html.unescape(match)

        return '<a href="/w/' + url_pas(data_unescape) + '">' + match + '</a>'
    
    def send_render_link(match):
        link_main = match[2]
        link_main = link_main.replace('"', '&quot;')

        return match[1] + '<a href="' + link_main + '">' + link_main + '</a>'

    if data == '&lt;br&gt;' or data == '' or re.search(r'^ +$', data):
        data = '<br>'
    else:
        data = data.replace('javascript:', '')

        data = re.sub(r'( |^)(https?:\/\/(?:[^ ]+))', send_render_link, data)
        data = re.sub(r'&lt;a(?:(?:(?!&gt;).)*)&gt;((?:(?!&lt;\/a&gt;).)+)&lt;\/a&gt;', send_render_href_replace, data)

    return data

async def recent_change(name = '', tool = '', num = 1, set_type = 'normal'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        
        all_admin = await acl_check(tool = 'all_admin_auth', ip = ip)
        all_admin = 1 if all_admin == 0 else 0

        owner = await acl_check(tool = 'owner_auth', ip = ip)
        owner = 1 if owner == 0 else 0

        option_list = [
            ['normal', get_lang(conn, 'normal')],
            ['edit', get_lang(conn, 'edit')],
            ['move', get_lang(conn, 'move')],
            ['delete', get_lang(conn, 'delete')],
            ['revert', get_lang(conn, 'revert')],
            ['r1', get_lang(conn, 'new_doc')],
            ['edit_request', get_lang(conn, 'edit_request')],
            ['file', get_lang(conn, 'file')],
            ['category', get_lang(conn, 'category')]
        ]
        if tool == 'history':
            option_list += [['setting', get_lang(conn, 'setting')]]

        if flask.request.method == 'POST':
            return redirect(conn, '/diff/' + flask.request.form.get('b', '1') + '/' + flask.request.form.get('a', '1') + '/' + url_pas(name))
        else:
            ban = ''
            select = ''
            sub = ''
            admin = owner
            div = '''
                <table id="main_table_set">
                    <tbody>
                        <tr id="main_table_top_tr">
            '''

            sql_num = (num * 50 - 50) if num * 50 > 0 else 0

            if tool == 'history':
                div += '''
                    <td id="main_table_width">''' + get_lang(conn, 'version') + '''</td>
                    <td id="main_table_width">''' + get_lang(conn, 'editor') + '''</td>
                    <td id="main_table_width">''' + get_lang(conn, 'time') + '''</td>
                '''
                sub = '(' + get_lang(conn, 'history') + ')'

                set_type = '' if set_type == 'edit' else set_type
                if set_type != 'normal':
                    curs.execute(db_change('select id, title, date, ip, send, leng, hide, type from history where title = ? and type = ? order by id + 0 desc limit ?, 50'), [name, set_type, sql_num])
                else:
                    curs.execute(db_change('select id, title, date, ip, send, leng, hide, type from history where title = ? order by id + 0 desc limit ?, 50'), [name, sql_num])

                data_list = curs.fetchall()
            elif tool == 'record':
                div +=  '''
                    <td id="main_table_width">''' + get_lang(conn, 'document_name') + '''</td>
                    <td id="main_table_width">''' + get_lang(conn, 'editor') + '''</td>
                    <td id="main_table_width">''' + get_lang(conn, 'time') + '''</td>
                '''
                sub = '(' + get_lang(conn, 'edit_record') + ')'
                set_type = '' if set_type == 'edit' else set_type

                if set_type != 'normal':
                    curs.execute(db_change('select id, title, date, ip, send, leng, hide, type from history where ip = ? and type = ? order by date desc limit ?, 50'), [name, set_type, sql_num])
                else:
                    curs.execute(db_change('select id, title, date, ip, send, leng, hide, type from history where ip = ? order by date desc limit ?, 50'), [name, sql_num])
                
                data_list = curs.fetchall()
            else:
                div +=  '''
                    <td id="main_table_width">''' + get_lang(conn, 'document_name') + '''</td>
                    <td id="main_table_width">''' + get_lang(conn, 'editor') + '''</td>
                    <td id="main_table_width">''' + get_lang(conn, 'time') + '''</td>
                '''
                sub = ''
                set_type = '' if set_type == 'edit' else set_type

                data_list = []

                if num == 1 or all_admin != 1:
                    curs.execute(db_change('select title, id from rc where type = ? order by date desc limit 50'), [set_type])
                    for for_a in curs.fetchall():
                        curs.execute(db_change('select id, title, date, ip, send, leng, hide, type from history where title = ? and id = ?'), for_a)
                        data_list += curs.fetchall()
                else:
                    if set_type != 'normal':
                        curs.execute(db_change('select id, title, date, ip, send, leng, hide, type from history where type = ? order by date desc limit ?, 50'), [set_type, sql_num])
                    else:
                        curs.execute(db_change('select id, title, date, ip, send, leng, hide, type from history order by date desc limit ?, 50'), [sql_num])

                    data_list = curs.fetchall()

            div += '</tr>'

            all_ip = await ip_pas([i[3] for i in data_list])
            for data in data_list:
                select += '<option value="' + data[0] + '">' + data[0] + '</option>'
                send = data[4]

                if re.search(r"\+", data[5]):
                    leng = '<span style="color:green;">(' + data[5] + ')</span>'
                elif re.search(r"\-", data[5]):
                    leng = '<span style="color:red;">(' + data[5] + ')</span>'
                else:
                    leng = '<span style="color:gray;">(' + data[5] + ')</span>'

                ip = all_ip[data[3]]
                m_tool = '<a href="/history_tool/' + data[0] + '/' + url_pas(data[1]) + '">(' + get_lang(conn, 'tool') + ')</a>'

                style = ['', '']
                date = data[2]

                type_data = ''
                if data[7] != '':
                    if data[7] == 'r1':
                        type_data = ' (' + data[7] + ')'
                    else:
                        type_data = ' (' + get_lang(conn, data[7]) + ')'

                send += type_data

                if data[6] == 'O':
                    if admin == 1:
                        style[0] = 'class="opennamu_history_blind"'
                        style[1] = 'class="opennamu_history_blind"'
                    else:
                        ip = ''
                        ban = ''
                        date = ''
                        send = ''

                        style[0] = 'style="display: none;"'
                        style[1] = 'class="opennamu_history_blind"'

                if tool == 'history':
                    if int(data[0]) < 2:
                        title = '<a href="/raw_rev/' + data[0] + '/' + url_pas(name) + '">r' + data[0] + '</a> '
                    else:
                        title = '<a href="/diff/' + str(int(data[0]) - 1) + '/' + data[0] + '/' + url_pas(name) + '">r' + data[0] + '</a> '
                else:
                    title = '<a href="/w/' + url_pas(data[1]) + '">' + html.escape(data[1]) + '</a> '
                    if int(data[0]) < 2:
                        title += '<a href="/history/' + url_pas(data[1]) + '">(r' + data[0] + ')</a> '
                    else:
                        title += '<a href="/diff/' + str(int(data[0]) - 1) + '/' + data[0] + '/' + url_pas(data[1]) + '">(r' + data[0] + ')</a> '

                div += '''
                    <tr ''' + style[0] + '''>
                        <td>''' + title + m_tool + ' ' + leng + '''</td>
                        <td>''' + ip + ban + '''</td>
                        <td>''' + date + '''</td>
                    </tr>
                    <tr ''' + style[1] + '''>
                        <td colspan="3">''' + recent_change_send_render(html.escape(send)) + '''</td>
                    </tr>
                '''

            div += '''
                    </tbody>
                </table>
            '''

            set_type = 'edit' if set_type == '' else set_type
            if tool == 'history':
                div = '' + \
                    ' '.join(['<a href="/history_page/1/' + for_a[0] + '/' + url_pas(name) + '">(' + for_a[1] + ')</a> ' for for_a in option_list]) + \
                    '<hr class="main_hr">' + div + \
                ''
                menu = [['w/' + url_pas(name), get_lang(conn, 'return')]]

                if set_type == 'normal':
                    div = '''
                        <form method="post">
                            <select name="a">''' + select + '''</select> <select name="b">''' + select + '''</select>
                            <button type="submit">''' + get_lang(conn, 'compare') + '''</button>
                        </form>
                        <hr class="main_hr">
                    ''' + div

                if admin == 1:
                    menu += [
                        ['history_add/' + url_pas(name), get_lang(conn, 'history_add')],
                        ['history_reset/' + url_pas(name), get_lang(conn, 'history_reset')]
                    ]

                title = name
                div += get_next_page_bottom(conn, '/history_page/{}/' + set_type + '/' + url_pas(name), num, data_list)
            elif tool == 'record':
                div = '' + \
                    ' '.join(['<a href="/record/1/' + for_a[0] + '/' + url_pas(name) + '">(' + for_a[1] + ')</a> ' for for_a in option_list]) + \
                    '<hr class="main_hr">' + div + \
                ''

                title = name
                menu = [['user/' + url_pas(name), get_lang(conn, 'user_tool')]]
                if admin == 1:
                    menu += [['record/reset/' + url_pas(name), get_lang(conn, 'record_reset')]]

                div += get_next_page_bottom(conn, '/record/{}/' + url_pas(set_type) + '/' + url_pas(name), num, data_list)
            else:
                div = '' + \
                    ' '.join(['<a href="/recent_change/1/' + for_a[0] + '">(' + for_a[1] + ')</a> ' for for_a in option_list]) + \
                    '<a href="/recent_change/1/user">(' + get_lang(conn, 'user_document') + ')</a> ' + \
                    '<hr class="main_hr">' + div + \
                ''

                menu = [['other', get_lang(conn, 'return')], ['recent_edit_request', get_lang(conn, 'edit_request')]]
                title = get_lang(conn, 'recent_change')

                if all_admin == 1:
                    div += get_next_page_bottom(conn, '/recent_change/{}/' + set_type, num, data_list)

            if sub == '':
                sub = 0

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [title, await wiki_set(), await wiki_custom(conn), wiki_css([sub, 0])],
                data = div,
                menu = menu
            ))