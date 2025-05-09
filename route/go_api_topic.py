from .tool.func import *

def api_topic_thread_make(user_id, date, data, code, color = '', blind = '', add_style = '', admin_check = 1, topic_num = ''):
    if blind == 'O':
        if data == '':
            color_b = 'opennamu_comment_blind'
        else:
            color_b = 'opennamu_comment_blind_admin'

        class_b = 'opennamu_comment_blind_js opennamu_list_hidden'
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
                    <td class="''' + color_b + ''' opennamu_comment_data_main" id="thread_''' + code + '''">
                        ''' + data + '''
                        <div id="opennamu_topic_req_''' + code + '''"></div>
                    </td>
                </tr>
            </table>
            <hr class="main_hr">
        </span>
    '''

async def api_topic_thread_pre_render(conn, data, num, ip, topic_num = '', name = '', sub = '', do_type = 'thread'):
    curs = conn.cursor()

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
                    await add_alarm(ip_data[0][0], ip, '<a href="/thread/' + topic_num + '#' + num + '">' + html.escape(name) + ' - ' + html.escape(sub) + '#' + num + '</a>')
                else:
                    set_id = topic_num.split('-')
                    set_id = ['', ''] if len(set_id) < 2 else set_id

                    await add_alarm(ip_data[0][0], ip, 'BBS <a href="/bbs/w/' + set_id[0] + '/' + set_id[1] + '#' + num + '">' + html.escape(name) + ' - ' + html.escape(sub) + '#' + num + '</a>')

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
                    await add_alarm(ip_data[0][0], ip, '<a href="/thread/' + topic_num + '#' + num + '">' + html.escape(name) + ' - ' + html.escape(sub) + '#' + num + '</a>')
                else:
                    set_id = topic_num.split('-')
                    await add_alarm(ip_data[0][0], ip, 'BBS <a href="/bbs/w/' + set_id[0] + '/' + set_id[1] + '#' + num + '">' + html.escape(name) + ' - ' + html.escape(sub) + '#' + num + '</a>')

            data = re.sub(call_user_regex, rd_data[0] + '<topic_call>@' + rd_data[1] + '</topic_call>' + rd_data[2], data, 1)

        call_user_count -= 1

    return data

async def api_topic(topic_num = 1, tool = 'normal', s_num = '', e_num = ''):
    with get_db_connect() as conn:
        topic_num = str(topic_num)

        if await acl_check('', 'topic_view', topic_num) != 1:
            other_set = {}
            other_set["topic_num"] = topic_num
            other_set["tool"] = tool
            other_set["s_num"] = str(s_num)
            other_set["e_num"] = str(e_num)

            return flask.jsonify(await python_to_golang(sys._getframe().f_code.co_name, other_set))
        else:
            return flask.jsonify({})