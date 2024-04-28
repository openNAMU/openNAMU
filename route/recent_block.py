from .tool.func import *

def recent_block(name = 'Test', tool = 'all', num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        for data in data_list:
            why = '<br>' if data[0] == '' else html.escape(data[0])
            if why == 'edit filter':
                why = '<a href="/edit_filter/' + url_pas(data[1]) + '">edit filter</a>'

            if data[5] == 'regex':
                ip = data[1]
                if data[6] == '1':
                    ip = '<s>' + ip + '</s> <a href="/auth/give/ban_regex/' + url_pas(data[1]) + '">(' + get_lang(conn, 'release') + ')</a>'
                else:
                    ip += ' <a href="/auth/give/ban_regex/' + url_pas(data[1]) + '">(' + get_lang(conn, 'ban') + ')</a>'

                ip += ' (' + get_lang(conn, 'regex') + ')'
            elif data[5] == 'cidr':
                ip = data[1]
                if data[6] == '1':
                    ip = '<s>' + ip + '</s> <a href="/auth/give/ban_cidr/' + url_pas(data[1]) + '">(' + get_lang(conn, 'release') + ')</a>'
                else:
                    ip += ' <a href="/auth/give/ban_cidr/' + url_pas(data[1]) + '">(' + get_lang(conn, 'ban') + ')</a>'

                ip += ' (' + get_lang(conn, 'cidr') + ')'
            else:
                ip = all_ip[data[1]]

            if data[3] == '':
                end = get_lang(conn, 'limitless')
            elif data[3] == 'release':
                end = get_lang(conn, 'release')
            else:
                end = data[3]

            if data[2] == '':
                admin = ''
            else:
                admin = all_ip[data[2]]

            start = get_lang(conn, 'start') + ' : ' + (data[4] if data[4] != '' else '0')
            div += '''
                <tr>
                    <td>''' + ip + '''</td>
                    <td>''' + admin + '''</td>
                    <td>
                        ''' + start + '''
                        <br>
                        ''' + get_lang(conn, 'end') + ' : ' + end + '''
                    </td>
                </tr>
                <tr>
                    <td colspan="3">''' + why + '''</td>
                </tr>
            '''

        div += '</tbody>'
        div += '</table>'

        if tool == 'all':
            div += next_fix(conn, '/block_log/', num, data_list)
        else:
            div += next_fix(conn, '/block_log/' + url_pas(tool) + '/' + url_pas(name) + '/', num, data_list)

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'recent_ban'), wiki_set(conn), wiki_custom(conn), wiki_css([sub, 0])],
            data = div,
            menu = menu
        ))