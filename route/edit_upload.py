from .tool.func import *

def edit_upload():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if acl_check(None, 'upload') == 1:
            return re_error('/ban')
        
        curs.execute(db_change('select data from other where name = "upload"'))
        db_data = curs.fetchall()
        file_max = number_check(db_data[0][0]) if db_data and db_data[0][0] != '' else '2'
        file_max = int(file_max)

        if flask.request.method == 'POST':
            if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return re_error('/error/13')
            else:
                captcha_post('', 0)

            file_data = flask.request.files.getlist("f_data[]", None)
            file_len = len(file_data)

            if (file_max * 1000 * 1000 * file_len) < flask.request.content_length:
                return re_error('/error/17')

            if file_len == 1:
                file_num = None
            else:
                if acl_check(None, 'many_upload') == 1:
                    return re_error('/ban')

                file_num = 1

            for data in file_data:
                if data.filename == '':
                    return re_error('/error/9')
                
                value = os.path.splitext(data.filename)[1]

                curs.execute(db_change("select html from html_filter where kind = 'extension'"))
                extension = [i[0].lower() for i in curs.fetchall()]
                if not re.sub(r'^\.', '', value).lower() in extension:
                    return re_error('/error/14')

                if flask.request.form.get('f_name', None):
                    name = flask.request.form.get('f_name', None) + (' ' + str(file_num) if file_num else '') + value
                else:
                    name = data.filename

                piece = os.path.splitext(name)
                if re.search(r'\.', piece[0]):
                    return re_error('/error/22')

                e_data = sha224_replace(piece[0]) + piece[1]

                curs.execute(db_change("select title from data where title = ?"), ['file:' + name])
                if curs.fetchall():
                    return re_error('/error/16')

                curs.execute(db_change("select html from html_filter where kind = 'file'"))
                db_data = curs.fetchall()
                for i in db_data:
                    t_re = re.compile(i[0])
                    if t_re.search(name):
                        return redirect('/filter/file_filter')

                data_url_image = load_image_url()
                if os.path.exists(os.path.join(data_url_image, e_data)):
                    return re_error('/error/16')
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

                render_set(
                    doc_name = 'file:' + name,
                    doc_data = file_d,
                    data_type = 'backlink'
                )

                history_plus(
                    'file:' + name,
                    file_d,
                    get_time(),
                    ip,
                    '',
                    '0',
                    t_check = 'upload',
                    mode = 'upload'
                )

                if file_num:
                    file_num += 1

                conn.commit()

            return redirect('/w/file:' + name)
        else:
            license_list = '<option value="direct_input">' + load_lang('direct_input') + '</option>'
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
            
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('upload'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <a href="/filter/file_filter">(''' + load_lang('file_filter_list') + ''')</a> <a href="/filter/extension_filter">(''' + load_lang('extension_filter_list') + ''')</a>
                    ''' + upload_help + '''
                    <hr class="main_hr">
                    ''' + load_lang('max_file_size') + ''' : ''' + str(file_max) + '''MB
                    <hr class="main_hr">
                    <form method="post" enctype="multipart/form-data" accept-charset="utf8">
                        <input multiple="multiple" type="file" name="f_data[]">
                        <hr class="main_hr">
                        <input placeholder="''' + load_lang('file_name') + '''" name="f_name" value="''' + file_name + '''">
                        <hr class="main_hr">
                        <select name="f_lice_sel">
                            ''' + license_list + '''
                        </select>
                        <hr class="main_hr">
                        <textarea class="opennamu_textarea_100" placeholder="''' + load_lang('other') + '''" name="f_lice">''' + upload_default + '''</textarea>
                        <hr class="main_hr">
                        ''' + captcha_get() + '''
                        <button id="opennamu_save_button" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['other', load_lang('return')]]
            ))