from .tool.func import *

def api_w_2(conn, name):
    curs = conn.cursor()

    data_arg_v = flask.request.args.get('v', '')
    if flask.request.method == 'POST':
        print(data_arg_v)
        if data_arg_v == '' or data_arg_v == 'preview':
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
        elif data_arg_v == 'include':
            name_org = flask.request.args.get('name_org', '')
            name_org = name if name_org == '' else name_org
            
            include_data = flask.request.args.get('include', '')
                
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
                        get_all_change_1 += [['@' + i[1] + '@', i[2]]]

                        json_data = json_data.replace(i[0], '@' + i[1] + '@', 1)
                    else:
                        json_data = json_data.replace(i[0], '@' + i[3] + '@', 1)

                get_all_change_2 = include_list + get_all_change_1
                print(get_all_change_2)
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
        elif data_arg_v == 'exist':
            try:
                title_list = json.loads(flask.request.form.get('title_list', ''))
                title_list = list(set(title_list))
            except:
                title_list = []
            
            data_exist = {}
            for i in title_list:
                curs.execute(db_change("select title from data where title = ?"), [i])
                if curs.fetchall():
                    data_exist[i] = '1'
                    
            return flask.jsonify(data_exist)
        else:
            return flask.jsonify({})
    else:
        data_arg_exist = flask.request.args.get('exist', '')
        if data_arg_v == 'exist' or data_arg_exist != '':
            curs.execute(db_change("select title from data where title = ?"), [name])
            if curs.fetchall():
                return flask.jsonify({ "exist" : "1" })
            else:
                return flask.jsonify({})
        else:
            data_arg_include = flask.request.args.get('include', '')
            if acl_check(name, 'render') == 1:
                return flask.jsonify({})
            else:
                data_arg_rev = flask.request.args.get('num', '')
                if data_arg_rev != '':
                    curs.execute(db_change("select data from history where title = ? and id = ?"), [name, rev])
                else:
                    curs.execute(db_change("select data from data where title = ?"), [name])

                sql_data = curs.fetchall()
                if not sql_data:
                    return flask.jsonify({})
                else:
                    data_pas = render_set(
                        doc_name = name, 
                        doc_data = sql_data[0][0], 
                        data_type = 'api_view'
                    )

                    return flask.jsonify({
                        "data" : data_pas[0], 
                        "js_data" : data_pas[1]
                    })