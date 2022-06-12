from .tool.func import *

def view_read(name = 'Test', doc_rev = 0, doc_from = '', do_type = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        sub = 0
        menu = []

        user_doc = ''
        category_doc = ''
        file_data = ''

        ip = ip_check()
            
        uppage = re.sub(r"/([^/]+)$", '', name)
        uppage = 0 if uppage == name else uppage
        num = str(doc_rev)        

        curs.execute(db_change("select sub from rd where title = ? and not stop = 'O' order by date desc"), [name])
        topic = 1 if curs.fetchall() else 0

        curs.execute(db_change("select title from data where title like ?"), [name + '/%'])
        down = 1 if curs.fetchall() else 0

        if re.search(r'^category:', name):
            category_doc = ''
            category_sub = ''

            curs.execute(db_change("select link from back where title = ? and type = 'cat' order by link asc"), [name])
            category_sql = curs.fetchall()
            for data in category_sql:
                if data[0].startswith('category:'):
                    category_sub += '<li><a href="/w/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a></li>'
                else:
                    category_doc += '<li><a href="/w/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a> <a id="inside" href="/xref/' + url_pas(data[0]) + '">(' + load_lang('backlink') + ')</a></li>'

            if category_doc != '':
                category_doc = '<h2 id="cate_normal">' + load_lang('category_title') + '</h2><ul class="inside_ul">' + category_doc + '</ul>'

            if category_sub != '':
                category_doc += '<h2 id="cate_under">' + load_lang('under_category') + '</h2><ul class="inside_ul">' + category_sub + '</ul>'
        elif re.search(r"^user:([^/]*)", name):
            match = re.search(r"^user:([^/]*)", name)
            user_name = html.escape(match.group(1))
            user_doc = '''
                <div id="get_user_info"></div>
                <script>load_user_info("''' + user_name + '''");</script>
                <hr class="main_hr">
            '''
            if name == 'user:' + user_name:
                menu += [['w/' + url_pas(name) + '/' + url_pas(get_time().split()[0]), load_lang('today_doc')]]
        elif re.search(r"^file:", name):
            mime_type = re.search(r'([^.]+)$', name)
            if mime_type:
                mime_type = mime_type.group(1).lower()
            else:
                mime_type = 'jpg'

            file_name = re.sub(r'\.([^.]+)$', '', name)
            file_name = re.sub(r'^file:', '', file_name)

            file_all_name = sha224_replace(file_name) + '.' + mime_type
            file_path_name = os.path.join(load_image_url(), file_all_name)
            if os.path.exists(file_path_name):
                file_size = str(round(os.path.getsize(file_path_name) / 1000, 1))
                file_data = '''
                    <img src="/image/''' + url_pas(file_all_name) + '''">
                    <h2>DATA</h2>
                    <table>
                        <tr><td>URL</td><td><a href="/image/''' + url_pas(file_all_name) + '''">LINK</a></td></tr>
                        <tr><td>VOLUME</td><td>''' + file_size + '''KB</td></tr>
                    </table>
                    <h2>CONTENT</h2>
                '''

                menu += [['delete_file/' + url_pas(name), load_lang('file_delete')]]
            else:
                file_data = ''

        if num != '0':
            curs.execute(db_change("select title from history where title = ? and id = ? and hide = 'O'"), [name, num])
            if curs.fetchall() and admin_check(6) != 1:
                return redirect('/history/' + url_pas(name))

            curs.execute(db_change("select data from history where title = ? and id = ?"), [name, num])
        else:
            curs.execute(db_change("select data from data where title = ?"), [name])

        data = curs.fetchall()
        end_data = render_set(
            doc_name = name,
            doc_data = data[0][0] if data else None
        )

        if end_data == 'HTTP Request 401.3':
            response_data = 401

            curs.execute(db_change('select data from other where name = "error_401"'))
            sql_d = curs.fetchall()
            if sql_d and sql_d[0][0] != '':
                end_data = '<h2>' + load_lang('error') + '</h2><ul class="inside_ul"><li>' + sql_d[0][0] + '</li></ul>'
            else:
                end_data = '<h2>' + load_lang('error') + '</h2><ul class="inside_ul"><li>' + load_lang('authority_error') + '</li></ul>'
        elif end_data == 'HTTP Request 404':
            response_data = 404

            curs.execute(db_change('select data from other where name = "error_404"'))
            sql_d = curs.fetchall()
            if sql_d and sql_d[0][0] != '':
                end_data = '<h2>' + load_lang('error') + '</h2><ul class="inside_ul"><li>' + sql_d[0][0] + '</li></ul>'
            else:
                end_data = '<h2>' + load_lang('error') + '</h2><ul class="inside_ul"><li>' + load_lang('decument_404_error') + '</li></ul>'

            curs.execute(db_change('' + \
                'select ip, date, leng, send, id from history ' + \
                'where title = ? and hide != "O" order by id + 0 desc limit 3' + \
            ''), [name])
            sql_d = curs.fetchall()
            if sql_d:
                end_data += '<h2>' + load_lang('history') + '</h2><ul class="inside_ul">'
                for i in sql_d:
                    if re.search(r"\+", i[2]):
                        leng = '<span style="color:green;">(' + i[2] + ')</span>'
                    elif re.search(r"\-", i[2]):
                        leng = '<span style="color:red;">(' + i[2] + ')</span>'
                    else:
                        leng = '<span style="color:gray;">(' + i[2] + ')</span>'

                    end_data += '<li>' + i[1] + ' | r' + i[4] + ' | ' + ip_pas(i[0]) + ' | ' + leng + (' | ' + i[3] if i[3] != '' else '') + '</li>'

                end_data += '<li><a href="/history/' + url_pas(name) + '">(...)</a></li></ul>'
        else:
            response_data = 200

        if num != '0':
            menu += [['history/' + url_pas(name), load_lang('history')]]
            sub = ' (r' + str(num) + ')'
            acl = 0
            r_date = 0
        else:
            curs.execute(db_change("select title from acl where title = ?"), [name])
            acl = 1 if curs.fetchall() else 0
            menu_acl = 1 if acl_check(name) == 1 else 0
            if response_data == 404:
                menu += [['edit/' + url_pas(name), load_lang('create'), menu_acl]] 
            else:
                menu += [['edit/' + url_pas(name), load_lang('edit'), menu_acl]]

            menu += [
                ['topic/' + url_pas(name), load_lang('discussion'), topic], 
                ['history/' + url_pas(name), load_lang('history')], 
                ['xref/' + url_pas(name), load_lang('backlink')], 
                ['acl/' + url_pas(name), load_lang('acl'), acl],
            ]

            if do_type == 'from':
                menu += [['w/' + url_pas(name), load_lang('pass')]]
                if flask.session and 'lastest_document' in flask.session:
                    end_data = '''
                        <div id="redirect">
                            <a href="/w_from/''' + url_pas(flask.session['lastest_document']) + '''">''' + flask.session['lastest_document'] + '''</a> ⇨ <b>''' + name + '''</b>
                        </div>
                        <br>
                    ''' + end_data
                    
                flask.session['lastest_document'] = name
            else:
                flask.session['lastest_document'] = name

            if uppage != 0:
                menu += [['w/' + url_pas(uppage), load_lang('upper')]]

            if down:
                menu += [['down/' + url_pas(name), load_lang('sub')]]

            curs.execute(db_change("select date from history where title = ? order by date desc limit 1"), [name])
            r_date = curs.fetchall()
            r_date = r_date[0][0] if r_date else 0

        div = file_data + user_doc + end_data + category_doc

        curs.execute(db_change("select data from other where name = 'body'"))
        body = curs.fetchall()
        div = (body[0][0] + div) if body else div

        curs.execute(db_change("select data from other where name = 'bottom_body'"))
        body = curs.fetchall()
        div += body[0][0] if body else ''

        if ip_or_user(ip) == 0:
            curs.execute(db_change("select title from scan where user = ? and title = ?"), [ip, name])
            watch_list = 2 if curs.fetchall() else 1
        else:
            watch_list = 0

        return easy_minify(flask.render_template(skin_check(),
            imp = [name, wiki_set(), wiki_custom(), wiki_css([sub, r_date, watch_list])],
            data = div,
            menu = menu
        )), response_data