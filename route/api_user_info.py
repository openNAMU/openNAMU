from .tool.func import *

def api_user_info(name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if flask.request.args.get('render', None):
            plus_d = ''
            plus_t = []

            curs.execute(db_change("update rb set ongoing = '' where end < ? and end != '' and ongoing = '1'"), [get_time()])
            conn.commit()

            plus_d = '''
                <table class="user_info_table">
                    <tbody>
                        <tr>
                            <td>''' + load_lang('user_name') + '''</td>
                            <td>{}</td>
                        </tr>
                        <tr>
                            <td>''' + load_lang('authority') + '''</td>
                            <td>{}</td>
                        </tr>
                        <tr>
                            <td>''' + load_lang('state') + '''</td>
                            <td>{}</td>
                        </tr>
                    </tbody>
                </table>
            '''

            curs.execute(db_change("select data from user_set where id = ? and name = 'acl'"), [name])
            data = curs.fetchall()
            if data:
                if data[0][0] != 'user':
                    curs.execute(db_change("select name from alist where name = ?"), [data[0][0]])
                    if curs.fetchall():
                        plus_t += [data[0][0]]
                    else:
                        plus_t += [load_lang('member')]
                else:
                    plus_t += [load_lang('member')]
            else:
                plus_t += [load_lang('normal')]

            if ban_check(name) == 0:
                plus_t += [load_lang('normal')]
            else:
                plus_t += [load_lang('blocked') + '<br>']
                regex_ban = 0

                curs.execute(db_change("select login, block, end, why from rb where band = 'regex' and ongoing = '1'"))
                for test_r in curs.fetchall():
                    if re.compile(test_r[1]).search(name):
                        plus_t[1] += load_lang('type') + ' : ' + load_lang('regex')
                        plus_t[1] += '<br>' + load_lang('period') + ' : ' + (test_r[2] if test_r[2] != '' else load_lang('limitless'))
                        plus_t[1] += ('<br>' + load_lang('login_able') if test_r[0] == 'O' else '')
                        plus_t[1] += ('<br>' + load_lang('why') + ' : ' + test_r[3] if test_r[3] != '' else '')
                        regex_ban = 1

                if regex_ban == 0:
                    curs.execute(db_change("select end, login, band, why from rb where block = ? and ongoing = '1'"), [name])
                    block_data = curs.fetchall()
                    if block_data:
                        plus_t[1] += load_lang('type') + ' : ' + (load_lang('band_blocked') if block_data[0][2] == 'O' else load_lang('normal'))
                        plus_t[1] += (' (' + load_lang('login_able') + ')' if block_data[0][1] != '' else '')
                        plus_t[1] += '<br>' + load_lang('period') + ' : ' + (block_data[0][0] if block_data[0][0] != '' else load_lang('limitless'))
                        plus_t[1] += ('<br>' + load_lang('band_blocked') if block_data[0][2] == 'O' else '')
                        plus_t[1] += ('<br>' + load_lang('why') + ' : ' + block_data[0][3] if block_data[0][3] != '' else '')

            plus_d = plus_d.format(ip_pas(name), plus_t[0], plus_t[1])

            return flask.jsonify({ "data" : plus_d })
        else:
            return flask.jsonify({})