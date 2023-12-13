from .tool.func import *

def api_w(name = 'Test', tool = '', rev = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if flask.request.method == 'POST':
            if tool == '' or tool == 'preview':
                # data_in 말고 data_use_type이랑 data_use_num 추가 예정
                data_org = flask.request.form.get('data', '')
                data_pas = render_set(
                    doc_name = name, 
                    doc_data = data_org, 
                    data_type = 'api_view'
                )

                return flask.jsonify({
                    "data" : data_pas[0], 
                    "js_data" : data_pas[1]
                })
            elif tool == 'include':
                name_org = flask.request.form.get('name_org', '')
                name_org = name if name_org == '' else name_org

                include_data = flask.request.form.get('name_include', '')

                try:
                    include_list = json.loads(flask.request.form.get('include_list', ''))
                except:
                    include_list = []

                curs.execute(db_change("select data from data where title = ?"), [name])
                sql_data = curs.fetchall()
                if not sql_data:
                    return flask.jsonify({})
                else:
                    json_data = sql_data[0][0]

                    get_all_change_1 = []
                    find_replace_moment = re.findall(r'(@([^=@]+)=([^=@]+)@|@([^=@]+)@)', json_data)
                    for i in find_replace_moment:
                        if i[1] != '':
                            get_all_change_1 += [[i[1], i[2]]]

                            json_data = json_data.replace(i[0], '@' + i[1] + '@', 1)
                        else:
                            json_data = json_data.replace(i[0], '@' + i[3] + '@', 1)

                    get_all_change_2 = include_list + get_all_change_1
                    for i in get_all_change_2:
                        json_data = json_data.replace('@' + i[0] + '@', i[1])

                    data_pas = render_set(
                        doc_name = name_org, 
                        doc_data = json_data, 
                        data_type = 'api_view',
                        data_in = include_data
                    )

                    return flask.jsonify({
                        "data" : data_pas[0], 
                        "js_data" : data_pas[1]
                    })
            elif tool == 'exist':
                try:
                    title_list = json.loads(flask.request.form.get('title_list', ''))
                    title_list = list(set(title_list))
                except:
                    title_list = [name]

                data_exist = {}
                for i in title_list:
                    curs.execute(db_change("select title from data where title = ?"), [i])
                    if curs.fetchall():
                        data_exist[i] = '1'

                return flask.jsonify(data_exist)

            return flask.jsonify({})
        else:
            if tool == '' or tool == 'view':
                if acl_check(name, 'render') != 1:
                    if number_check(rev) == '':
                        curs.execute(db_change("select data from data where title = ?"), [name])
                    else:
                        curs.execute(db_change("select data from history where title = ? and id = ?"), [name, rev])

                    sql_data = curs.fetchall()
                    if sql_data:
                        data_pas = render_set(
                            doc_name = name, 
                            doc_data = sql_data[0][0], 
                            data_type = 'api_view'
                        )

                        return flask.jsonify({
                            "data" : data_pas[0], 
                            "js_data" : data_pas[1]
                        })

            return flask.jsonify({})