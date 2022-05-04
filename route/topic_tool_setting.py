from .tool.func import *

def topic_tool_setting(topic_num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check(3) != 1:
            return re_error('/error/3')

        ip = ip_check()
        time = get_time()
        topic_num = str(topic_num)

        curs.execute(db_change("select stop, agree from rd where code = ?"), [topic_num])
        rd_d = curs.fetchall()
        if not rd_d:
            return redirect('/')

        if flask.request.method == 'POST':
            admin_check(3, 'change_topic_set (code ' + topic_num + ')')

            curs.execute(db_change("select id from topic where code = ? order by id + 0 desc limit 1"), [topic_num])
            topic_check = curs.fetchall()
            if topic_check:
                stop_d = flask.request.form.get('stop_d', '')
                why_d = flask.request.form.get('why', '')
                agree_d = flask.request.form.get('agree', '')

                curs.execute(db_change("update rd set stop = ?, agree = ? where code = ?"), [
                    stop_d,
                    agree_d,
                    topic_num
                ])

                if stop_d == 'S':
                    t_state = 'Stop'
                elif stop_d == 'O':
                    t_state = 'Close'
                else:
                    t_state = 'Normal'

                curs.execute(db_change("insert into topic (id, data, date, ip, top, code) values (?, ?, ?, ?, '1', ?)"), [
                    str(int(topic_check[0][0]) + 1),
                    t_state + (' (Agree)' if agree_d != '' else '') + (('[br][br]Why : ' + why_d) if why_d else ''),
                    time,
                    ip,
                    topic_num
                ])

                rd_plus(topic_num, time)

            return redirect('/thread/' + topic_num)
        else:
            stop_d_list = ''
            agree_check = ''
            for_list = [
                ['O', 'Close'],
                ['S', 'Stop'],
                ['', 'Normal']
            ]

            for i in for_list:
                if rd_d and rd_d[0][0] == i[0]:
                    stop_d_list = '<option value="' + i[0] + '">' + i[1] + '</option>' + stop_d_list
                else:
                    stop_d_list += '<option value="' + i[0] + '">' + i[1] + '</option>'

            agree_check = 'checked="checked"' if rd_d[0][1] == 'O' else ''

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('topic_setting'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <select name="stop_d">
                            ''' + stop_d_list + '''
                        </select>
                        <hr class=\"main_hr\">
                        <input type="checkbox" name="agree" value="O" ''' + agree_check + '''> Agree
                        <hr class=\"main_hr\">
                        <input placeholder="''' + load_lang('why') + ''' (''' + load_lang('markup_enabled') + ''')" name="why" type="text">
                        <hr class=\"main_hr\">
                        <button type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['thread/' + topic_num + '/tool', load_lang('return')]]
            ))