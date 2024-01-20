from .tool.func import *

def api_topic_thread_make(user_id, date, data, code, color = '', blind = '', add_style = '', admin_check = 1, topic_num = ''):
    if blind == 'O':
        if data == '':
            color_b = 'opennamu_comment_blind'
        else:
            color_b = 'opennamu_comment_blind_admin'

        class_b = 'opennamu_comment_blind_js'
    else:
        color_b = 'opennamu_comment_blind_not'
        class_b = ''

    admin_check_box = ''
    if admin_check == 1 and topic_num != '':
        admin_check_box = '<input type="checkbox" class="opennamu_blind_button" id="opennamu_blind_' + topic_num + '_' + code + '">'

    return '''
        <span class="''' + class_b + '''">
            <table class="opennamu_comment" style="''' + add_style + '''">
                <tr>
                    <td class="opennamu_comment_color_''' + color + '''">
                        ''' + admin_check_box + '''
                        <a href="#thread_shortcut" id="''' + code + '''">#''' + code + '''</a>
                        ''' + user_id + '''
                        <span style="float: right;">''' + date + '''</span>
                    </td>
                </tr>
                <tr>
                    <td class="''' + color_b + '''" id="opennamu_comment_data_main">
                        ''' + data + '''
                    </td>
                </tr>
            </table>
            <hr class="main_hr">
        </span>
    '''

def api_topic_thread_pre_render(curs, data, num, ip, topic_num = '', name = '', sub = '', do_type = 'thread'):
    # 이거 좀 엉성해서 언젠간 손 보고 싶음

    call_thread_regex = r"( |\n|^)(?:#([0-9]+)(?:-([0-9]+))?)( |\n|$)"
    call_thread_count = len(re.findall(call_thread_regex, data)) * 3
    while 1:
        rd_data = re.search(call_thread_regex, data)
        if call_thread_count < 0:
            break
        elif not rd_data:
            break
        else:
            rd_data = rd_data.groups()

            view_data = rd_data[1]
            send_topic_num = topic_num
            if rd_data[2]:
                view_data += '-' + rd_data[2]
                if do_type == 'thread':
                    send_topic_num = rd_data[2]
                else:
                    set_id = topic_num.split('-')

                    send_topic_num = set_id[0] + '-' + rd_data[2]
                    view_data += '-' + set_id[0]

            if do_type == 'thread':
                curs.execute(db_change("select ip from topic where code = ? and id = ?"), [send_topic_num, rd_data[1]])
            else:
                if rd_data[1] == '0':
                    set_id = send_topic_num.split('-')
                    set_id = ['', ''] if len(set_id) < 2 else set_id

                    curs.execute(db_change('select set_data from bbs_data where set_name = "user_id" and set_id = ? and set_code = ?'), [set_id[0], set_id[1]])
                else:
                    curs.execute(db_change('select set_data from bbs_data where set_name = "comment_user_id" and set_id = ? and set_code = ?'), [send_topic_num, rd_data[1]])

            ip_data = curs.fetchall()
            if ip_data and ip_or_user(ip_data[0][0]) == 0:
                if do_type == 'thread':
                    add_alarm(ip_data[0][0], ip, '<a href="/thread/' + topic_num + '#' + num + '">' + html.escape(name) + ' - ' + html.escape(sub) + '#' + num + '</a>')
                else:
                    set_id = topic_num.split('-')
                    set_id = ['', ''] if len(set_id) < 2 else set_id

                    add_alarm(ip_data[0][0], ip, 'BBS <a href="/bbs/w/' + set_id[0] + '/' + set_id[1] + '#' + num + '">' + html.escape(name) + ' - ' + html.escape(sub) + '#' + num + '</a>')

            data = re.sub(call_thread_regex, rd_data[0] + '<topic_a_' + do_type + '>#' + view_data + '</topic_a_' + do_type + '>' + rd_data[3], data, 1)

        call_thread_count -= 1

    call_user_regex = r"( |\n|^)(?:@([^ \n]+))( |\n|$)"
    call_user_count = len(re.findall(call_user_regex, data)) * 3
    while 1:
        rd_data = re.search(call_user_regex, data)
        if call_user_count < 0:
            break
        elif not rd_data:
            break
        else:
            rd_data = rd_data.groups()

            curs.execute(db_change("select ip from history where ip = ? limit 1"), [rd_data[1]])
            ip_data = curs.fetchall()
            if not ip_data:
                curs.execute(db_change("select ip from topic where ip = ? limit 1"), [rd_data[1]])
                ip_data = curs.fetchall()

            if ip_data and ip_or_user(ip_data[0][0]) == 0:
                if do_type == 'thread':
                    add_alarm(ip_data[0][0], ip, '<a href="/thread/' + topic_num + '#' + num + '">' + html.escape(name) + ' - ' + html.escape(sub) + '#' + num + '</a>')
                else:
                    set_id = topic_num.split('-')
                    add_alarm(ip_data[0][0], ip, 'BBS <a href="/bbs/w/' + set_id[0] + '/' + set_id[1] + '#' + num + '">' + html.escape(name) + ' - ' + html.escape(sub) + '#' + num + '</a>')

            data = re.sub(call_user_regex, rd_data[0] + '<topic_call>@' + rd_data[1] + '</topic_call>' + rd_data[2], data, 1)

        call_user_count -= 1

    return data

def api_topic(topic_num = 1, tool = 'normal', num = '', render = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        topic_num = str(topic_num)

        if acl_check('', 'topic_view', topic_num) != 1:
            if tool == 'normal':
                if num != '':
                    curs.execute(db_change("select id, data, date, ip, block, top from topic where code = ? and id + 0 = ? + 0 order by id + 0 asc"), [topic_num, num])
                else:
                    curs.execute(db_change("select id, data, date, ip, block, top from topic where code = ? order by id + 0 asc"), [topic_num])
            elif tool == 'top':
                curs.execute(db_change("select id, data, date, ip, block, top from topic where code = ? and top = 'O' order by id + 0 asc"), [topic_num])
            else:
                # tool == 'length'
                curs.execute(db_change("select id from topic where code = ? order by id + 0 desc limit 1"), [topic_num])
                db_data = curs.fetchall()
                if db_data:
                    return flask.jsonify({ 'length' : db_data[0][0] })
                else:
                    return flask.jsonify({})

            data = curs.fetchall()
            if data:
                data_a = {}
                admin = admin_check(3)

                curs.execute(db_change("select ip from topic where code = ? order by id + 0 asc limit 1"), [topic_num])
                data_f = curs.fetchall()
                data_f = data_f[0][0] if data_f else ''
                data_a['data_main'] = {
                    "ip_first" : ip_pas(data_f, 1),
                    "admin" : str(admin)
                }
                data_a['data'] = []

                ip_a = ip_pas([i[3] for i in data])
                ip_a_2 = ip_pas([i[3] for i in data], 1)
                for i in data:
                    data_v = i[1] if i[4] != 'O' or admin == 1 else ''
                    if data_v != '':
                        data_v = render_set(
                            doc_data = data_v, 
                            data_type = 'api_thread',
                            data_in = 'topic_' + topic_num + '_' + i[0]
                        )
                    else:
                        data_v = ['', '']

                    data_a['data'] += [{
                        "id" : i[0],

                        "data" : data_v,
                        "date" : i[2],
                        "ip" : ip_a_2[i[3]],
                        "blind" : i[4],

                        "ip_pas" : ip_a[i[3]],
                        "data_pas" : data_v
                    }]

                if render == '':
                    return flask.jsonify(data_a)
                else:
                    data_r = ''
                    if 'data' in data_a:
                        for for_a in data_a['data']:
                            if tool == 'top':
                                color = 'red'
                            elif for_a["blind"] == '1':
                                color = 'blue'
                            elif data_a['data_main']["ip_first"] == for_a["ip"]:
                                color = 'green'
                            else:
                                color = 'default'

                            data_r += api_topic_thread_make(
                                for_a["ip_pas"],
                                '<a href="/thread/' + topic_num + '/comment/' + for_a["id"] + '/tool">(' + load_lang('tool') + ')</a> ' + for_a["date"],
                                for_a["data_pas"][0] + '<script>' + for_a["data_pas"][1] + '</script>',
                                for_a["id"],
                                color = color,
                                blind = for_a["blind"],
                                add_style = '',
                                admin_check = admin if tool == 'normal' else 0,
                                topic_num = topic_num
                            )

                    if admin == 1 and tool == 'normal':
                        data_r += '''
                            <a href="javascript:opennamu_thread_blind();">(''' + load_lang('hide') + ''' | ''' + load_lang('hide_release') + ''')</a>
                            <hr class="main_hr">
                        '''

                    return flask.jsonify({ "data" : data_r })
            else:
                return flask.jsonify({})
        else:
            return flask.jsonify({})