from .tool.func import *

def api_w_2(conn, name):
    curs = conn.cursor()

    if flask.request.args.get('exist', None):
        curs.execute(db_change("select title from data where title = ?"), [name])
        if curs.fetchall():
            return flask.jsonify({ "exist" : "1" })
        else:
            return flask.jsonify({})
    else:
        if acl_check(name, 'render') != 1:
            if flask.request.method == 'POST':
                g_data = render_set(title = name, data = flask.request.form.get('data', ''), num = 2)

                return flask.jsonify({ "title" : name, "data" : g_data[0], "js_data" : g_data[1] })
            else:
                curs.execute(db_change("select data from data where title = ?"), [name])
                data = curs.fetchall()
                if data:
                    if flask.request.args.get('include', 'include_1'):
                        include_re = re.compile(r'\[include\(((?:(?!\)\]).)+)\)\]', re.I)
                        category_re = re.compile(r'\[\[(?:(?:category|분류):(?:(?!\[\[|\]\]).)+)\]\]', re.I)

                        json_data = include_re.sub('', data[0][0])
                        json_data = category_re.sub('', json_data)

                        get_all_change_1 = [('@' + i[0] + '@', i[1]) for i in re.findall(r'@([^=]+)=([^@]+)@', json_data)]
                        json_data = re.sub(r'@(?P<in>[^=]+)=([^@]+)@', '@\g<in>@', json_data)

                        get_all_change_2 = re.findall(r'(@(?:[^@]*)@),([^,]*),', flask.request.args.get('change', '')) + get_all_change_1
                        for i in get_all_change_2:
                            json_data = json_data.replace(
                                i[0].replace('<amp>', '&'), 
                                i[1].replace('<amp>', '&').replace('<comma>', ',')
                            )

                        g_data = render_set(title = name, data = json_data, num = 2, include = flask.request.args.get('include', 'include_1'))
                    else:
                        json_data = g_data[0]

                        g_data = render_set(title = name, data = json_data, num = 2)

                    return flask.jsonify({ "title" : name, "data" : g_data[0], "js_data" : g_data[1] })
                else:
                    return flask.jsonify({})
        else:
            return flask.jsonify({})