from .tool.func import *

def recent_block_2(conn, name, tool):
    curs = conn.cursor()

    num = int(number_check(flask.request.args.get('num', '1')))
    sql_num = (num * 50 - 50) if num * 50 > 0 else 0

    div = '''
        <table id="main_table_set">
            <tbody>
                <tr>
                    <td id="main_table_width">''' + load_lang('blocked') + '''</td>
                    <td id="main_table_width">''' + load_lang('admin') + '''</td>
                    <td id="main_table_width">''' + load_lang('period') + '''</td>
                </tr>
    '''

    curs.execute(db_change("update rb set ongoing = '' where end < ? and end != '' and ongoing = '1'"), [get_time()])
    conn.commit()

    if not name:
        if flask.request.args.get('type', '') == 'ongoing':
            sub = ' (' + load_lang('in_progress') + ')'
            menu = [['block_log', load_lang('normal')]]

            curs.execute(db_change("" + \
                "select why, block, blocker, end, today, band from rb " + \
                "where ((end > ? and end like '2%') or end = '') and ongoing = '1' " + \
                "order by end desc limit ?, 50"
            ), [
                get_time(),
                sql_num
            ])
        else:
            sub = 0
            menu = 0

            div = '' + \
                '<a href="/manager/11">(' + load_lang('blocked') + ')</a> ' + \
                '<a href="/manager/12">(' + load_lang('admin') + ')</a> ' + \
                '<a href="?type=ongoing">(' + load_lang('in_progress') + ')</a>' + \
                '<hr class="main_hr">' + \
            '' + div

            curs.execute(db_change("select why, block, blocker, end, today, band from rb order by today desc limit ?, 50"), [sql_num])
    else:
        menu = [['block_log', load_lang('normal')]]

        if tool == 'block_user':
            sub = ' (' + load_lang('blocked') + ')'

            curs.execute(db_change("select why, block, blocker, end, today, band from rb where block = ? order by today desc limit ?, 50"), [
                name, 
                sql_num
            ])
        else:
            sub = ' (' + load_lang('admin') + ')'

            curs.execute(db_change("select why, block, blocker, end, today, band from rb where blocker = ? order by today desc limit ?, 50"), [
                name, 
                sql_num
            ])

    data_list = curs.fetchall()
    all_ip = ip_pas([i[1] for i in data_list] + [i[2] for i in data_list])
    for data in data_list:
        why = '<br>' if data[0] == '' else html.escape(data[0])

        if data[5] == 'O':
            ip = data[1] + ' (' + load_lang('range') + ')'
        elif data[5] == 'regex':
        else:
            ip = all_ip[data[1]]

        if data[3] == '':
            end = load_lang('limitless')
        elif data[3] == 'release':
            end = load_lang('release')
        else:
            end = data[3]

        if data[2] == '':
            admin = ''
        elif re.search(r'^tool:', data[2]):
            admin = data[2]
        else:
            admin = all_ip[data[2]]

        start = load_lang('start') + ' : ' + (data[4] if data[4] != '' else '0')
        div += '''
            <tr>
                <td>''' + ip + '''</td>
                <td>''' + admin + '''</td>
                <td>
                    ''' + start + '''
                    <br>
                    ''' + load_lang('end') + ' : ' + end + '''
                </td>
            </tr>
            <tr>
                <td colspan="3">''' + why + '''</td>
            </tr>
        '''

    div += '</tbody></table>'
    div += next_fix('/block_log?num=', num, data_list) if not name else next_fix('/' + url_pas(tool) + '/' + url_pas(name) + '?num=', num, data_list)

    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('recent_ban'), wiki_set(), custom(), other2([sub, 0])],
        data = div,
        menu = menu
    ))