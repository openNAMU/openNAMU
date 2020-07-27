from .tool.func import *

def api_topic_sub_2(conn, topic_num):
    curs = conn.cursor()

    topic_num = str(topic_num)
    get_acl = acl_check('', 'render')

    if flask.request.args.get('num', None):
        curs.execute(db_change("select id, data, date, ip, block, top from topic where code = ? and id + 0 = ? + 0 order by id + 0 asc"), [
            topic_num,
            flask.request.args.get('num', '')
        ])
    elif flask.request.args.get('top', None):
        curs.execute(db_change("select id, data, date, ip, block, top from topic where code = ? and top = 'O' order by id + 0 asc"), [topic_num])
    else:
        curs.execute(db_change("select id, data, date, ip, block, top from topic where code = ? order by id + 0 asc"), [topic_num])

    data = curs.fetchall()
    if data:
        json_data = {}
        admin = admin_check(3)

        if flask.request.args.get('render', None):
            all_ip = ip_pas([i[3] for i in data])
            for i in data:
                ip = all_ip[i[3]]

                if i[4] != 'O':
                    t_data_f = i[1]
                    b_color = 'toron_color'
                else:
                    t_data_f = ''
                    b_color = 'toron_color_not'

                    ip += ' <a href="/admin_log?search=blind%20(code%20' + topic_num + '#' + i[0] + '">(B)</a>'

                    if admin == 1:
                        ip += ' <a href="/thread/' + topic_num + '/raw/' + i[0] + '">(R)</a>'

                if i[0] == '1':
                    s_user = i[3]
                else:
                    if flask.request.args.get('num', None):
                        curs.execute(db_change("select ip from topic where code = ? order by id + 0 asc limit 1"), [topic_num])
                        g_data = curs.fetchall()
                        if g_data:
                            s_user = g_data[0][0]
                        else:
                            s_user = ''

                if flask.request.args.get('top', None):
                    t_color = 'toron_color_red'
                elif i[3] == s_user and i[5] != '1':
                    t_color = 'toron_color_green'
                elif i[5] == '1':
                    t_color = 'toron_color_blue'
                else:
                    t_color = 'toron_color'

                if admin == 1 or b_color != 'toron_color_not':
                    ip += ' <a href="/thread/' + topic_num + '/admin/' + i[0] + '">(' + load_lang('discussion_tool') + ')</a>'

                if t_data_f == '':
                    t_data_f = '[br]'

                t_data_f = render_set(data = t_data_f, num = 2, include = 'topic_' + i[0], acl = get_acl)
                t_plus_data = t_data_f[1]
                t_data_f = t_data_f[0]

                t_data_f = re.sub(
                    r'&lt;topic_a&gt;((?:(?!&lt;\/topic_a&gt;).)+)&lt;\/topic_a&gt;', 
                    '<a href="\g<1>">\g<1></a>', 
                    t_data_f
                )
                t_data_f = re.sub(
                    r'&lt;topic_call&gt;@((?:(?!&lt;\/topic_call&gt;).)+)&lt;\/topic_call&gt;', 
                    '<a href="/w/user:\g<1>">@\g<1></a>', 
                    t_data_f
                )

                all_data = '' + \
                    '<table id="toron">' + \
                        '<tbody>' + \
                            '<tr>' + \
                               '<td id="' + t_color + '">' + \
                                    '<a href="javascript:void(0);" id="' + i[0] + '">#' + i[0] + '</a> ' + ip + ' <span style="float: right;">' + i[2] + '</span>' + \
                                '</td>' + \
                            '</tr>' + \
                            '<tr>' + \
                                '<td id="' + b_color + '">' + \
                                    '<div id="topic_scroll">' + t_data_f + '</div>' + \
                                '</td>' + \
                            '</tr>' + \
                        '</tbody>' + \
                    '</table>' + \
                    '<hr class="main_hr">' + \
                ''

                json_data[i[0]] = {
                    "data" : all_data,
                    "plus_data" : t_plus_data
                }
        else:
            for i in data:
                if i[4] != 'O' or (i[4] == 'O' and admin == 1):
                    t_data_f = i[1]
                else:
                    t_data_f = '(B)'

                json_data[i[0]] = {
                    "data" : t_data_f,
                    "date" : i[2],
                    "ip" : ip_pas(i[3], 1),
                    "block" : i[4],
                }

        return flask.jsonify(json_data)
    else:
        return flask.jsonify({})