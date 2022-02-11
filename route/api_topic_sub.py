from .tool.func import *

def api_topic_sub(topic_num = 1, tool = 'normal', num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        topic_num = str(topic_num)

        if tool == 'normal':
            if num != '':
                curs.execute(db_change(
                    "select id, data, date, ip, block, top from topic where code = ? and id + 0 = ? + 0 order by id + 0 asc"
                ), [
                    topic_num,
                    num
                ])
            else:
                curs.execute(db_change(
                    "select id, data, date, ip, block, top from topic where code = ? order by id + 0 asc"
                ), [
                    topic_num
                ]) 
        else:
            curs.execute(db_change(
                "select id, data, date, ip, block, top from topic where code = ? and top = 'O' order by id + 0 asc"
            ), [
                topic_num
            ])
            
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
                    "data_pas" : render_set(
                        doc_data = data_v, 
                        data_type = 'api_view',
                        data_in = 'topic_' + topic_num + '_' + i[0],
                        doc_acl = 0
                    )
                }

            return flask.jsonify(data_a)
        else:
            return flask.jsonify({})