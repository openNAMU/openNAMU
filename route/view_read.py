from .tool.func import *

def view_read_2(conn, name):
    curs = conn.cursor()

    sub = ''
    acl = 0
    div = ''
    ip = ip_check()

    num = flask.request.args.get('num', None)
    if num:
        num = int(number_check(num))
    else:
        if not flask.request.args.get('from', None):
            curs.execute(db_change("select title from back where link = ? and type = 'redirect'"), [name])
            r_db = curs.fetchall()
            if r_db:
                r_data = link_fix(r_db[0][0])

                return redirect('/w/' + url_pas(r_data[0]) + '?from=' + name + r_data[1])

    curs.execute(db_change("select sub from rd where title = ? and not stop = 'O' order by date desc"), [name])
    if curs.fetchall():
        topic = 1
    else:
        topic = 0

    curs.execute(db_change("select link from back where title = ? and type = 'cat' order by link asc"), [name])

    curs.execute(db_change("select title from data where title like ?"), ['%' + name + '/%'])
    if curs.fetchall():
        down = 1
    else:
        down = 0

    m = re.search("^(.*)\/(.*)$", name)
    if m:
        uppage = m.groups()[0]
    else:
        uppage = 0

    if re.search('^category:', name):
        curs.execute(db_change("select link from back where title = ? and type = 'cat' order by link asc"), [name])
        back = curs.fetchall()
        if back:
            div = '<br><h2 id="cate_normal">' + load_lang('category') + '</h2><ul>'
            u_div = ''

            for data in back:
                if re.search('^category:', data[0]):
                    u_div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'
                else:
                    curs.execute(db_change("select title from back where title = ? and type = 'include'"), [data[0]])
                    db_data = curs.fetchall()
                    if db_data:
                        div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a> <a id="inside" href="/xref/' + url_pas(data[0]) + '">(' + load_lang('backlink') + ')</a></li>'
                    else:
                        div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'

            div += '</ul>'

            if div == '<br><h2 id="cate_normal">' + load_lang('category') + '</h2><ul></ul>':
                div = ''

            if u_div != '':
                div += '<br><h2 id="cate_under">' + load_lang('under_category') + '</h2><ul>' + u_div + '</ul>'


    cache_data = None
    if num:
        curs.execute(db_change("select title from history where title = ? and id = ? and hide = 'O'"), [name, str(num)])
        if curs.fetchall() and admin_check(6) != 1:
            return redirect('/history/' + url_pas(name))

        curs.execute(db_change("select data from history where title = ? and id = ?"), [name, str(num)])
    else:
        curs.execute(db_change("select id from history where title = ? order by id + 0 desc limit 1"), [name])
        last_history_num = curs.fetchall()
        if last_history_num and not flask.request.args.get('reload', None):
            curs.execute(db_change("select data from cache_data where title = ? and id = ?"), [name, last_history_num[0][0]])
            cache_data = curs.fetchall()
            if not cache_data:
                curs.execute(db_change("select data from data where title = ?"), [name])
        else:
            curs.execute(db_change("select data from data where title = ?"), [name])

    if cache_data:
        end_data = cache_data[0][0]
    else:
        data = curs.fetchall()
        if data:
            else_data = data[0][0]
        else:
            else_data = None

        curs.execute(db_change("select decu from acl where title = ?"), [name])
        data = curs.fetchall()
        if data:
            acl = 1

        if flask.request.args.get('from', None) and else_data:
            else_data = re.sub('^\r\n', '', else_data)
            else_data = re.sub('\r\n$', '', else_data)

        end_data = render_set(
            title = name,
            data = else_data
        )

        if not num:
            curs.execute(db_change("delete from cache_data where title = ?"), [name])
            if last_history_num:
                curs.execute(db_change("insert into cache_data (title, data, id) values (?, ?, ?)"), [name, end_data, last_history_num[0][0]])

    if end_data == 'HTTP Request 401.3':
        response_data = 401

        curs.execute(db_change('select data from other where name = "error_401"'))
        sql_d = curs.fetchall()
        if sql_d and sql_d[0][0] != '':
            end_data = '<h2>' + load_lang('error') + '</h2><ul><li>' + sql_d[0][0] + '</li></ul>'
        else:
            end_data = '<h2>' + load_lang('error') + '</h2><ul><li>' + load_lang('authority_error') + '</li></ul>'
    elif end_data == 'HTTP Request 404':
        response_data = 404

        curs.execute(db_change('select data from other where name = "error_404"'))
        sql_d = curs.fetchall()
        if sql_d and sql_d[0][0] != '':
            end_data = '<h2>' + load_lang('error') + '</h2><ul><li>' + sql_d[0][0] + '</li></ul>'
        else:
            end_data = '<h2>' + load_lang('error') + '</h2><ul><li>' + load_lang('decument_404_error') + '</li></ul>'

        curs.execute(db_change('select ip, date, leng, send from history where title = ? order by id desc limit 3'), [name])
        sql_d = curs.fetchall()
        if sql_d:
            end_data += '<h2>' + load_lang('history') + '</h2><ul>'
            for i in sql_d:
                if re.search("\+", i[2]):
                    leng = '<span style="color:green;">(' + i[2] + ')</span>'
                elif re.search("\-", i[2]):
                    leng = '<span style="color:red;">(' + i[2] + ')</span>'
                else:
                    leng = '<span style="color:gray;">(' + i[2] + ')</span>'

                end_data += '<li>' + i[1] + ' | ' + ip_pas(i[0]) + ' | ' + leng + (' | ' + i[3] if i[3] != '' else '') + '</li>'

            end_data += '<li><a href="/history/' + url_pas(name) + '">(...)</a></li></ul>'
    else:
        response_data = 200

    if num:
        menu = [['history/' + url_pas(name), load_lang('history')]]
        sub = ' (r' + str(num) + ')'
        acl = 0
        r_date = 0
    else:
        if response_data == 404:
            menu = [['edit/' + url_pas(name), load_lang('create')]]
        else:
            menu = [['edit/' + url_pas(name), load_lang('edit')]]

        menu += [
            ['topic/' + url_pas(name), load_lang('discussion'), topic], 
            ['history/' + url_pas(name), load_lang('history')], 
            ['xref/' + url_pas(name), load_lang('backlink')], 
            ['acl/' + url_pas(name), load_lang('acl'), acl],
            ['w/' + url_pas(name) + '?reload=true', load_lang('reload')]
        ]

        if flask.request.args.get('from', None):
            menu += [['w/' + url_pas(name), load_lang('pass')]]
            end_data = '''
                <div id="redirect">
                    <a href="/w/''' + url_pas(flask.request.args.get('from', None)) + '?from=' + url_pas(name) + '">' + flask.request.args.get('from', None) + '</a> ⇨ <b>' + name + '''</b>
                </div>
                <br>
            ''' + end_data

        if uppage != 0:
            menu += [['w/' + url_pas(uppage), load_lang('upper')]]

        if down:
            menu += [['down/' + url_pas(name), load_lang('sub')]]

        curs.execute(db_change("select date from history where title = ? order by date desc limit 1"), [name])
        date = curs.fetchall()
        if date:
            r_date = date[0][0]
        else:
            r_date = 0

    div = end_data + div

    adsense_code = '<div align="center" style="display: block; margin-bottom: 10px;">{}</div>'

    curs.execute(db_change("select data from other where name = 'adsense'"))
    adsense_enabled = curs.fetchall()[0][0]
    if adsense_enabled == 'True':
        curs.execute(db_change("select data from other where name = 'adsense_code'"))
        adsense_code = adsense_code.format(curs.fetchall()[0][0])
    else:
        adsense_code = adsense_code.format('')

    div = adsense_code + '<div>' + div + '</div>'

    match = re.search("^user:([^/]*)", name)
    if match:
        user_name = match.groups()[0]
        div = '''
            <div id="get_user_info"></div>
            <script>load_user_info("''' + user_name + '''");</script>
        ''' + div

    curs.execute(db_change("select data from other where name = 'body'"))
    body = curs.fetchall()
    if body:
        div = body[0][0] + '<hr class=\"main_hr\">' + div

    curs.execute(db_change("select data from other where name = 'bottom_body'"))
    body = curs.fetchall()
    if body:
        div += '<hr class=\"main_hr\">' + body[0][0]

    if ip_or_user(ip) == 0:
        curs.execute(db_change("select title from scan where user = ? and title = ?"), [ip, name])
        if curs.fetchall():
            watch_list = 2
        else:
            watch_list = 1
    else:
        watch_list = 0

    return easy_minify(flask.render_template(skin_check(),
        imp = [flask.request.args.get('show', name), wiki_set(), custom(), other2([sub, r_date, watch_list])],
        data = div,
        menu = menu
    )), response_data