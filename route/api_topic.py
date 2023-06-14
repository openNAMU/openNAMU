from .tool.func import *
from .bbs_w_post import bbs_w_post_make_thread

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
                            data_type = 'api_view',
                            data_in = 'topic_' + topic_num + '_' + i[0],
                            doc_acl = 0
                        )
                        data_v[0] = re.sub(
                            r'&lt;topic_a&gt;(?P<in>(?:(?!&lt;\/topic_a&gt;).)+)&lt;\/topic_a&gt;',
                            '<a href="\g<in>">\g<in></a>',
                            data_v[0]
                        )
                        data_v[0] = re.sub(
                            r'&lt;topic_call&gt;@(?P<in>(?:(?!&lt;\/topic_call&gt;).)+)&lt;\/topic_call&gt;',
                            '<a href="/w/user:\g<in>">@\g<in></a>',
                            data_v[0]
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
                            else:
                                if data_a['data_main']["ip_first"] == for_a["ip"]:
                                    color = 'green'
                                else:
                                    color = 'default'

                            data_r += bbs_w_post_make_thread(
                                for_a["ip_pas"],
                                '<a href="/thread/' + topic_num + '/comment/' + for_a["id"] + '/tool">(' + load_lang('tool') + ')</a> ' + for_a["date"],
                                for_a["data_pas"][0] + '<script>' + for_a["data_pas"][1] + '</script>',
                                for_a["id"],
                                color = color,
                                blind = for_a["blind"],
                                add_style = ''
                            )
                            data_r += '<hr class="main_hr">'

                    return flask.jsonify({ "data" : data_r })
            else:
                return flask.jsonify({})
        else:
            return flask.jsonify({})