from .tool.func import *

def api_w_2(conn, name):
    curs = conn.cursor()

    if flask.request.args.get('exist', None):
        curs.execute(db_change("select title from data where title = ?"), [name])
        if curs.fetchall():
            return flask.jsonify({ "exist" : "1" })
    else:
        if acl_check(name, 'render') != 1:
            if flask.request.method == 'POST':
                g_data = flask.request.form.get('data', '')
                g_data = render_set(title = name, data = g_data, num = 2)

                return flask.jsonify({ "title" : name, "data" : g_data[0], "js_data" : g_data[1] })
            else:
                rev = flask.request.args.get('num', '')
                if rev != '':
                    curs.execute(db_change("select data from history where title = ? and id = ?"), [name, rev])
                else:
                    curs.execute(db_change("select data from data where title = ?"), [name])
                data = curs.fetchall()
                if data:
                    json_data = data[0][0]
                    include_data = flask.request.args.get('include', None)
                    if include_data:
                        get_all_change_1 = []
                        find_replace_moment = re.findall(r'(@([^=@]+)=([^=@]+)@|@([^=@]+)@)', json_data)
                        for i in find_replace_moment:
                            if i[1] != '':
                                get_all_change_1 += [['@' + i[1] + '@', i[2]]]

                                json_data = json_data.replace(i[0], '@' + i[1] + '@', 1)
                            else:
                                json_data = json_data.replace(i[0], '@' + i[3] + '@', 1)

                        get_all_change_2 = re.findall(r'(@(?:[^@]*)@),([^,]*),', flask.request.args.get('change', '')) + get_all_change_1
                        for i in get_all_change_2:
                            json_data = json_data.replace(
                                i[0].replace('<amp>', '&'), 
                                i[1].replace('<amp>', '&').replace('<comma>', ','), 
                                1
                            )
                        
                    g_data = render_set(title = name, data = json_data, num = 2, include = include_data)

                    return flask.jsonify({ "title" : name, "data" : g_data[0], "js_data" : g_data[1] })

    return flask.jsonify({})
