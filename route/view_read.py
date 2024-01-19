from .tool.func import *

def view_read(name = 'Test', doc_rev = '', doc_from = '', do_type = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        sub = 0
        menu = []

        user_doc = ''
        category_total = ''
        file_data = ''

        doc_type = ''
        now_time = get_time()

        ip = ip_check()
            
        uppage = re.sub(r"/([^/]+)$", '', name)
        uppage = 0 if uppage == name else uppage
        num = str(doc_rev)

        curs.execute(db_change("select sub from rd where title = ? and not stop = 'O' order by date desc"), [name])
        topic = 1 if curs.fetchall() else 0

        curs.execute(db_change("select title from data where title like ?"), [name + '/%'])
        down = 1 if curs.fetchall() else 0

        if re.search(r'^category:', name):
            name_view = name
            doc_type = 'category'

            category_doc = ''
            category_sub = ''

            count_sub_category = 0
            count_category = 0

            curs.execute(db_change("select distinct link from back where title = ? and type = 'cat' order by link asc"), [name])
            category_sql = curs.fetchall()
            for data in category_sql:
                link_view = data[0]
                if get_main_skin_set(curs, flask.session, 'main_css_category_change_title', ip) != 'off':
                    curs.execute(db_change("select data from back where title = ? and link = ? and type = 'cat_view' limit 1"), [name, data[0]])
                    db_data = curs.fetchall()
                    if db_data and db_data[0][0] != '':
                        link_view = db_data[0][0]
                        
                link_blur = ''
                curs.execute(db_change("select data from back where title = ? and link = ? and type = 'cat_blur' limit 1"), [name, data[0]])
                db_data = curs.fetchall()
                if db_data:
                    link_blur = 'opennamu_category_blur'

                if data[0].startswith('category:'):
                    category_sub += '<li><a class="' + link_blur + '" href="/w/' + url_pas(data[0]) + '">' + html.escape(link_view) + '</a></li>'
                    count_sub_category += 1
                else:
                    category_doc += '' + \
                        '<li>' + \
                            '<a class="' + link_blur + '" href="/w/' + url_pas(data[0]) + '">' + html.escape(link_view) + '</a> ' + \
                            '<a class="opennamu_link_inter" href="/xref/' + url_pas(data[0]) + '">(' + load_lang('backlink') + ')</a>' + \
                        '</li>' + \
                    ''
                    count_category += 1

            if category_sub != '':
                category_total += '' + \
                    '<h2 id="cate_under">' + load_lang('under_category') + '</h2>' + \
                    '<ul class="opennamu_ul">' + \
                        '<li>' + load_lang('all') + ' : ' + str(count_sub_category) + '</li>' + \
                        category_sub + \
                    '</ul>' + \
                ''

            if category_doc != '':
                category_total += '' + \
                    '<h2 id="cate_normal">' + load_lang('category_title') + '</h2>' + \
                    '<ul class="opennamu_ul">' + \
                        '<li>' + load_lang('all') + ' : ' + str(count_category) + '</li>' + \
                        category_doc + \
                    '</ul>' + \
                ''
        elif re.search(r"^user:([^/]*)", name):
            name_view = name
            doc_type = 'user'

            match = re.search(r"^user:([^/]*)", name)
            
            user_name = html.escape(match.group(1))
            user_doc = ''
            
            # S admin or owner 특수 틀 추가
            if admin_check('all', None, user_name) == 1:
                if admin_check(None, None, user_name) == 1:
                    curs.execute(db_change('select data from other where name = "phrase_user_page_owner"'))
                    db_data = curs.fetchall()
                    if db_data and db_data[0][0] != '':
                        user_doc += db_data[0][0] + '<br>'
                    else:
                        curs.execute(db_change('select data from other where name = "phrase_user_page_admin"'))
                        db_data = curs.fetchall()
                        if db_data and db_data[0][0] != '':
                            user_doc += db_data[0][0] + '<br>'
                else:
                    curs.execute(db_change('select data from other where name = "phrase_user_page_admin"'))
                    db_data = curs.fetchall()
                    if db_data and db_data[0][0] != '':
                        user_doc += db_data[0][0] + '<br>'
            # E
            
            user_doc += '''
                <div id="opennamu_get_user_info">''' + html.escape(user_name) + '''</div>
                <hr class="main_hr">
            '''
            if name == 'user:' + user_name:
                menu += [['w/' + url_pas(name) + '/' + url_pas(now_time.split()[0]), load_lang('today_doc')]]
        elif re.search(r"^file:", name):
            curs.execute(db_change('select id from history where title = ? order by date desc limit 1'), [name])
            db_data = curs.fetchall()
            rev = db_data[0][0] if db_data else '1' 

            name_view = name
            doc_type = 'file'

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
                try:
                    img = Image.open(file_path_name)
                    width, height = img.size
                    file_res = str(width) + 'x' + str(height)
                except:
                    file_res = 'Vector'
                
                file_size = str(round(os.path.getsize(file_path_name) / 1000, 1))
                
                file_data = '''
                    <img src="/image/''' + url_pas(file_all_name) + '''.cache_v''' + rev + '''">
                    <h2>''' + load_lang('data') + '''</h2>
                    <table>
                        <tr><td>URL</td><td><a href="/image/''' + url_pas(file_all_name) + '''">''' + load_lang('link') + '''</a></td></tr>
                        <tr><td>''' + load_lang('volume') + '''</td><td>''' + file_size + '''KB</td></tr>
                        <tr><td>''' + load_lang('resolution') + '''</td><td>''' + file_res + '''</td></tr>
                    </table>
                    <h2>''' + load_lang('content') + '''</h2>
                '''

                menu += [['delete_file/' + url_pas(name), load_lang('file_delete')]]
            else:
                file_data = ''
        else:
            name_view = name

        if num != '':
            curs.execute(db_change("select title from history where title = ? and id = ? and hide = 'O'"), [name, num])
            if curs.fetchall() and admin_check(6) != 1:
                return redirect('/history/' + url_pas(name))

            curs.execute(db_change("select data from history where title = ? and id = ?"), [name, num])
        else:
            curs.execute(db_change("select data from data where title = ?"), [name])

        data = curs.fetchall()
        end_data = render_set(
            doc_name = name,
            doc_data = data[0][0] if data else None,
            data_type = 'from' if do_type == 'from' else 'view'
        )

        if end_data == 'HTTP Request 401.3':
            response_data = 401

            curs.execute(db_change('select data from other where name = "error_401"'))
            sql_d = curs.fetchall()
            if sql_d and sql_d[0][0] != '':
                end_data = '<h2>' + load_lang('error') + '</h2><ul class="opennamu_ul"><li>' + sql_d[0][0] + '</li></ul>'
            else:
                end_data = '<h2>' + load_lang('error') + '</h2><ul class="opennamu_ul"><li>' + load_lang('authority_error') + '</li></ul>'
        elif end_data == 'HTTP Request 404':
            response_data = 404

            curs.execute(db_change('select data from other where name = "error_404"'))
            sql_d = curs.fetchall()
            if sql_d and sql_d[0][0] != '':
                end_data = '<h2>' + load_lang('error') + '</h2><ul class="opennamu_ul"><li>' + sql_d[0][0] + '</li></ul>'
            else:
                end_data = '<h2>' + load_lang('error') + '</h2><ul class="opennamu_ul"><li>' + load_lang('decument_404_error') + '</li></ul>'

            curs.execute(db_change('' + \
                'select ip, date, leng, send, id from history ' + \
                'where title = ? and hide != "O" order by id + 0 desc limit 3' + \
            ''), [name])
            sql_d = curs.fetchall()
            if sql_d:
                end_data += '<h2>' + load_lang('history') + '</h2><ul class="opennamu_ul">'
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

        if num != '':
            menu += [['history/' + url_pas(name), load_lang('return')]]
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
                ['acl/' + url_pas(name), load_lang('setting'), acl],
            ]

            if flask.session and 'lastest_document' in flask.session:
                if type(flask.session['lastest_document']) != type([]):
                    flask.session['lastest_document'] = []
            else:
                flask.session['lastest_document'] = []

            if do_type == 'from':
                menu += [['w/' + url_pas(name), load_lang('pass')]]
                
                last_page = ''
                for for_a in reversed(range(0, len(flask.session['lastest_document']))):
                    last_page = flask.session['lastest_document'][for_a]

                    curs.execute(db_change("select link from back where (title = ? or link = ?) and type = 'redirect' limit 1"), [last_page, last_page])
                    if curs.fetchall():
                        break

                redirect_text = '{0} ➤ {1}'

                curs.execute(db_change('select data from other where name = "redirect_text"'))
                db_data = curs.fetchall()
                if db_data and db_data[0][0] != '':
                    redirect_text = db_data[0][0]

                try:
                    redirect_text = redirect_text.format('<a href="/w_from/' + url_pas(last_page) + '">' + html.escape(last_page) + '</a>', '<b>' + html.escape(name) + '</b>')
                except:
                    redirect_text = '{0} ➤ {1}'
                    redirect_text = redirect_text.format('<a href="/w_from/' + url_pas(last_page) + '">' + html.escape(last_page) + '</a>', '<b>' + html.escape(name) + '</b>')

                end_data = '''
                    <div id="redirect">
                        ''' + redirect_text + '''
                    </div>
                    <hr class="main_hr">
                ''' + end_data
                    
            if len(flask.session['lastest_document']) >= 10:
                flask.session['lastest_document'] = flask.session['lastest_document'][-9:] + [name]
            else:
                flask.session['lastest_document'] += [name]
            
            flask.session['lastest_document'] = list(reversed(dict.fromkeys(reversed(flask.session['lastest_document']))))

            view_history_on = get_main_skin_set(curs, flask.session, 'main_css_view_history', ip)
            if view_history_on == 'on':
                end_data = '' + \
                    '<div id="redirect">' + \
                        load_lang('trace') + ' : ' + \
                        ' ➥ '.join(
                            [
                                '<a href="/w/' + url_pas(for_a) + '">' + html.escape(for_a) + '</a>'
                                for for_a in flask.session['lastest_document']
                            ]
                        ) + \
                    '</div>' + \
                    '<hr class="main_hr">' + \
                '' + end_data

            if uppage != 0:
                menu += [['w/' + url_pas(uppage), load_lang('upper')]]

            if down:
                menu += [['down/' + url_pas(name), load_lang('sub')]]

            curs.execute(db_change("select set_data from data_set where doc_name = ? and set_name = 'last_edit'"), [name])
            r_date = curs.fetchall()
            r_date = r_date[0][0] if r_date else 0

        div = file_data + user_doc + end_data + category_total

        if num != '':
            curs.execute(db_change('select data from other where name = "phrase_old_page_warning"'))
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '':
                div = db_data[0][0] + '<hr class="main_hr">' + div

            doc_type = 'rev'
        
        if doc_type == '':
            curs.execute(db_change('select data from other where name = "outdated_doc_warning_date"'))
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '' and r_date != 0:
                time_1 = datetime.datetime.strptime(r_date, '%Y-%m-%d %H:%M:%S') + datetime.timedelta(days = int(number_check(db_data[0][0])))
                time_2 = datetime.datetime.strptime(now_time, '%Y-%m-%d %H:%M:%S')
                if time_2 > time_1:
                    curs.execute(db_change('select data from other where name = "outdated_doc_warning"'))
                    db_data = curs.fetchall()
                    div = (db_data[0][0] if db_data and db_data[0][0] != '' else load_lang('old_page_warning')) + '<hr class="main_hr">' + div

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
            imp = [name_view, wiki_set(), wiki_custom(), wiki_css([sub, r_date, watch_list])],
            data = div,
            menu = menu
        )), response_data