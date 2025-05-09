from .tool.func import *

async def edit_upload():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'upload') == 1:
            return await re_error(conn, 0)
        
        curs.execute(db_change('select data from other where name = "upload"'))
        db_data = curs.fetchall()
        file_max = number_check(db_data[0][0]) if db_data and db_data[0][0] != '' else '2'
        file_max = int(file_max)

        if flask.request.method == 'POST':
            if await captcha_post(conn, flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return await re_error(conn, 13)

            file_data = flask.request.files.getlist("f_data[]")
            file_len = len(file_data)

            file_size_all = flask.request.content_length
            if file_size_all == None:
                file_size_all = 0

            if (file_max * 1000 * 1000 * file_len) < file_size_all or file_size_all == 0:
                return await re_error(conn, 17)

            if file_len == 1:
                file_num = None
            else:
                if await acl_check('', 'many_upload') == 1:
                    return await re_error(conn, 0)

                file_num = 1

            for data in file_data:
                file_name = data.filename if data.filename else ''
                if file_name == '':
                    return await re_error(conn, 9)
                
                value_tmp = os.path.splitext(file_name)
                value = ''
                if len(value_tmp) >= 2:
                    value = value_tmp[1]

                curs.execute(db_change("select html from html_filter where kind = 'extension'"))
                extension = [i[0].lower() for i in curs.fetchall()]
                if not re.sub(r'^\.', '', value).lower() in extension:
                    return await re_error(conn, 14)

                name = ''
                if flask.request.form.get('f_name', None):
                    name = flask.request.form.get('f_name', '') + (' ' + str(file_num) if file_num else '') + value
                else:
                    name = file_name

                piece = os.path.splitext(name)
                if re.search(r'\.', piece[0]):
                    return await re_error(conn, 22)

                e_data = sha224_replace(piece[0]) + piece[1]

                curs.execute(db_change("select title from data where title = ?"), ['file:' + name])
                if curs.fetchall():
                    return await re_error(conn, 16)

                curs.execute(db_change("select html from html_filter where kind = 'file'"))
                db_data = curs.fetchall()
                for i in db_data:
                    t_re = re.compile(i[0])
                    if t_re.search(name):
                        return redirect(conn, '/filter/file_filter')

                data_url_image = load_image_url(conn)
                if os.path.exists(os.path.join(data_url_image, e_data)):
                    return await re_error(conn, 16)
                else:
                    data.save(os.path.join(data_url_image, e_data))

                ip = ip_check()
                g_lice = flask.request.form.get('f_lice', '')
                file_size = os.stat(os.path.join(data_url_image, e_data)).st_size
                file_size = str(round(file_size / 1000, 1))

                curs.execute(db_change("select data from other where name = 'markup'"))
                db_data = curs.fetchall()
                if db_data and db_data[0][0] == 'namumark':
                    file_d = '' + \
                        flask.request.form.get('f_lice_sel', 'direct_input') + '\n' + \
                        '[[category:' + re.sub(r'\]', '_', flask.request.form.get('f_lice_sel', '')) + ']]\n' + \
                        (g_lice if g_lice != '' else '') + \
                    ''
                else:
                    file_d = '' + \
                        flask.request.form.get('f_lice_sel', 'direct_input') + '\n' + \
                        (g_lice if g_lice != '' else '') + \
                    ''

                curs.execute(db_change("insert into data (title, data) values (?, ?)"), ['file:' + name, file_d])

                render_set(conn, 
                    doc_name = 'file:' + name,
                    doc_data = file_d,
                    data_type = 'backlink'
                )

                history_plus(conn, 
                    'file:' + name,
                    file_d,
                    get_time(),
                    ip,
                    '',
                    '0',
                    mode = 'upload'
                )

                if file_num:
                    file_num += 1

            return redirect(conn, '/w/file:' + name)
        else:
            license_list = '<option value="direct_input">' + get_lang(conn, 'direct_input') + '</option>'
            file_name = html.escape(flask.request.args.get('name', ''))

            curs.execute(db_change("select html from html_filter where kind = 'image_license'"))
            db_data = curs.fetchall()
            license_list += ''.join(['<option value="' + i[0] + '">' + i[0] + '</option>' for i in db_data])

            curs.execute(db_change("select data from other where name = 'upload_help'"))
            db_data = curs.fetchall()
            upload_help = ('<hr class="main_hr">' + db_data[0][0]) if db_data and db_data[0][0] != '' else ''

            curs.execute(db_change("select data from other where name = 'upload_default'"))
            db_data = curs.fetchall()
            upload_default = html.escape(db_data[0][0]) if db_data and db_data[0][0] != '' else ''
            
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'upload'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <a href="/filter/file_filter">(''' + get_lang(conn, 'file_filter_list') + ''')</a> <a href="/filter/extension_filter">(''' + get_lang(conn, 'extension_filter_list') + ''')</a>
                    ''' + upload_help + '''
                    <hr class="main_hr">
                    ''' + get_lang(conn, 'max_file_size') + ''' : ''' + str(file_max) + '''MB
                    <hr class="main_hr">
                    <form method="post" enctype="multipart/form-data" accept-charset="utf8">
                        <input multiple="multiple" type="file" name="f_data[]">
                        <hr class="main_hr">
                        <input placeholder="''' + get_lang(conn, 'file_name') + '''" name="f_name" value="''' + file_name + '''">
                        <hr class="main_hr">
                        <select name="f_lice_sel">
                            ''' + license_list + '''
                        </select>
                        <hr class="main_hr">
                        <textarea class="opennamu_textarea_100" placeholder="''' + get_lang(conn, 'other') + '''" name="f_lice">''' + upload_default + '''</textarea>
                        <hr class="main_hr">
                        ''' + await captcha_get(conn) + '''
                        <button id="opennamu_save_button" type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                ''',
                menu = [['other', get_lang(conn, 'return')]]
            ))