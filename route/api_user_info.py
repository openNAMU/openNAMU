from .tool.func import *

def api_user_info(name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if flask.request.method == 'POST':
            try:
                data_list = json.loads(flask.request.form.get('title_list', ''))
                data_list = list(set(title_list))
            except:
                data_list = [name]

            data_result = {}
            for user_name in data_list:
                data_result[user_name] = {}
                
                # auth part
                curs.execute(db_change("select data from user_set where id = ? and name = 'acl'"), [user_name])
                db_data = curs.fetchall()
                if db_data:
                    if db_data[0][0] != 'user':
                        curs.execute(db_change("select name from alist where name = ?"), [db_data[0][0]])
                        if curs.fetchall():
                            data_result[user_name]['auth'] = db_data[0][0]
                        else:
                            data_result[user_name]['auth'] = 1
                    else:
                        data_result[user_name]['auth'] = 1
                else:
                    data_result[user_name]['auth'] = 0
                
                # user document part
                curs.execute(db_change("select title from data where title = ?"), ['user:' + user_name])
                if curs.fetchall():
                    data_result[user_name]['document'] = 1
                else:
                    data_result[user_name]['document'] = 0

                # user title part
                curs.execute(db_change('select data from user_set where name = "user_title" and id = ?'), [user_name])
                db_data = curs.fetchall()
                if db_data:
                    data_result[user_name]['user_title'] = db_data[0][0]
                else:
                    data_result[user_name]['user_title'] = ''
                    
            return flask.jsonify(data_result)
        else:
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