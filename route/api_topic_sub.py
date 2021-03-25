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
        data_a = {}
        admin = admin_check(3)

        curs.execute(db_change("select ip from topic where code = ? order by id + 0 asc limit 1"), [topic_num])
        data_f = curs.fetchall()
        data_f = data_f[0][0] if data_f else ''
        data_a['data_main'] = {
            "ip_first" : ip_pas(data_f, 1),
            "admin" : str(admin)
        }
            
        ip_a = ip_pas([i[3] for i in data])
        ip_a_2 = ip_pas([i[3] for i in data], 1)
        for i in data:
            data_v = i[1] if i[4] != 'O' or admin == 1 else ''
            data_a[i[0]] = {
                "data" : data_v,
                "date" : i[2],
                "ip" : ip_a_2[i[3]],
                "blind" : i[4],
                
                "ip_pas" : ip_a[i[3]],
                "data_pas" : render_set(data = data_v, num = 2, include = 'topic_' + i[0], acl = get_acl)
            }

        return flask.jsonify(data_a)
    else:
        return flask.jsonify({})