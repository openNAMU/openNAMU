from .tool.func import *

def recent_block(name = 'Test', tool = 'all', num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        sql_num = (num * 50 - 50) if num * 50 > 0 else 0

        div = '''
            <table id="main_table_set">
                <tbody>
                    <tr id="main_table_top_tr">
                        <td id="main_table_width">''' + load_lang('blocked') + '''</td>
                        <td id="main_table_width">''' + load_lang('admin') + '''</td>
                        <td id="main_table_width">''' + load_lang('period') + '''</td>
                    </tr>
        '''

        div = '' + \
            '<a href="/block_log">(' + load_lang('all') + ')</a> ' + \
            '<a href="/manager/11">(' + load_lang('blocked') + ')</a> ' + \
            '<a href="/manager/12">(' + load_lang('admin') + ')</a> ' + \
            '<a href="/block_log/ongoing">(' + load_lang('in_progress') + ')</a> ' + \
            '<a href="/block_log/regex">(' + load_lang('regex') + ')</a>' + \
            '<hr class="main_hr">' + \
        '' + div

        if tool == 'all':
            sub = 0
            menu = [['other', load_lang('return')]]

            curs.execute(db_change("select why, block, blocker, end, today, band, ongoing from rb order by today desc limit ?, 50"), [sql_num])
        elif tool == 'ongoing':
            sub = '(' + load_lang('in_progress') + ')'
            menu = [['other', load_lang('return')]]

            curs.execute(db_change("select why, block, blocker, end, today, band, ongoing from rb where ongoing = '1' order by end desc limit ?, 50"), [sql_num])
        elif tool == 'regex':
            sub = '(' + load_lang('regex') + ')'
            menu = [['other', load_lang('return')]]

            curs.execute(db_change("select why, block, blocker, end, today, band, ongoing from rb where band = 'regex' order by today desc limit ?, 50"), [sql_num])
        elif tool == 'user':
            sub = '(' + load_lang('blocked') + ')'
            menu = [['other', load_lang('return')]]

            curs.execute(db_change("select why, block, blocker, end, today, band, ongoing from rb where block = ? order by today desc limit ?, 50"), [name, sql_num])
        else:
            sub = '(' + load_lang('admin') + ')'
            menu = [['other', load_lang('return')]]

            curs.execute(db_change("select why, block, blocker, end, today, band, ongoing from rb where blocker = ? order by today desc limit ?, 50"), [name, sql_num])

        data_list = curs.fetchall()
        all_ip = ip_pas([i[1] for i in data_list] + [i[2] for i in data_list])
        for data in data_list:
            why = '<br>' if data[0] == '' else html.escape(data[0])
            if why == 'edit filter':
                why = '<a href="/edit_filter/' + url_pas(data[1]) + '">edit filter</a>'

            if data[5] == 'regex':
                ip = data[1]
                if data[6] == '1':
                    ip = '<s>' + ip + '</s> <a href="/auth/give/ban_regex/' + url_pas(data[1]) + '">(' + load_lang('release') + ')</a>'
                else:
                    ip += ' <a href="/auth/give/ban_regex/' + url_pas(data[1]) + '">(' + load_lang('ban') + ')</a>'

                ip += ' (' + load_lang('regex') + ')'
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

        div += '</tbody>'
        div += '</table>'

        if tool == 'all':
            div += next_fix('/block_log/', num, data_list)
        else:
            div += next_fix('/block_log/' + url_pas(tool) + '/' + url_pas(name) + '/', num, data_list)

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('recent_ban'), wiki_set(), wiki_custom(), wiki_css([sub, 0])],
            data = div,
            menu = menu
        ))