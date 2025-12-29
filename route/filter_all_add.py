from .tool.func import *

async def filter_all_add(tool, name = None):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if not name and tool == 'edit_filter':
            return redirect(conn, '/manager/9')

        if flask.request.method == 'POST':
            if await acl_check('', 'owner_auth', '', '') == 1:
                return await re_error(conn, 3)

            title = flask.request.form.get('title', 'test')
            if tool in ('inter_wiki', 'outer_link'):
                link = flask.request.form.get('link', 'test')
                icon = flask.request.form.get('icon', '')
                inter_type = flask.request.form.get('inter_type', '')

                curs.execute(db_change("delete from html_filter where html = ? and kind = ?"), [title, tool])
                curs.execute(db_change('insert into html_filter (html, plus, plus_t, kind) values (?, ?, ?, ?)'), [title, link, icon, tool])
                if tool == 'inter_wiki':
                    curs.execute(db_change("delete from html_filter where html = ? and kind = 'inter_wiki_sub'"), [title])
                    curs.execute(db_change('insert into html_filter (html, plus, plus_t, kind) values (?, "inter_wiki_type", ?, "inter_wiki_sub")'), [title, inter_type])
                
                await acl_check(tool = 'owner_auth', memo = tool + ' edit')
            elif tool == 'edit_filter':
                day = flask.request.form.get('day', '0')
                end = 'X' if day == '0' else day
                if end != 'X':
                    end = re.sub(r'[^0-9]', '', end)
                    end = str(int(number_check(end)) * 24 * 60 * 60)

                content = flask.request.form.get('content', 'test')
                try:
                    re.compile(content)
                except:
                    return await re_error(conn, 23)
                
                curs.execute(db_change("delete from html_filter where html = ? and kind = 'regex_filter'"), [name])
                curs.execute(db_change("insert into html_filter (html, plus, plus_t, kind) values (?, ?, ?, 'regex_filter')"), [name, content, end])
                await acl_check(tool = 'owner_auth', memo = 'edit_filter edit')
            elif tool == 'document':
                post_name = flask.request.form.get('name', '')
                if post_name == '':
                    return redirect(conn, '/filter/document')
            
                post_acl = flask.request.form.get('acl', '')
                post_regex = flask.request.form.get('regex', '')
                try:
                    re.compile(post_regex)
                except:
                    return await re_error(conn, 23)
                
                curs.execute(db_change('insert into html_filter (html, kind, plus, plus_t) values (?, "document", ?, ?)'), [post_name, post_regex, post_acl])
                await acl_check(tool = 'owner_auth', memo = 'document_filter edit')
            else:
                plus_d = ''
                if tool == 'name_filter':
                    try:
                        re.compile(title)
                    except:
                        return await re_error(conn, 23)

                    await acl_check(tool = 'owner_auth', memo = 'name_filter edit')
                    type_d = 'name'
                elif tool == 'file_filter':
                    try:
                        re.compile(title)
                    except:
                        return await re_error(conn, 23)

                    await acl_check(tool = 'owner_auth', memo = 'file_filter edit')
                    type_d = 'file'
                elif tool == 'email_filter':
                    await acl_check(tool = 'owner_auth', memo = 'email_filter edit')
                    type_d = 'email'
                elif tool == 'image_license':
                    await acl_check(tool = 'owner_auth', memo = 'image_license edit')
                    type_d = 'image_license'
                elif tool == 'extension_filter':
                    await acl_check(tool = 'owner_auth', memo = 'extension_filter edit')
                    type_d = 'extension'
                    plus_d = flask.request.form.get('max_file_size', '')
                    if plus_d != '':
                        plus_d = number_check(plus_d)
                elif tool == 'template':
                    await acl_check(tool = 'owner_auth', memo = 'template_document edit')
                    type_d = 'template'
                    plus_d = flask.request.form.get('exp', 'test')
                else:
                    await acl_check(tool = 'owner_auth', memo = 'edit_top edit')
                    type_d = 'edit_top'
                    plus_d = flask.request.form.get('markup', 'test')

                if name:
                    curs.execute(db_change("delete from html_filter where html = ? and kind = ?"), [name, type_d])

                curs.execute(db_change('insert into html_filter (html, kind, plus, plus_t) values (?, ?, ?, ?)'), [title, type_d, plus_d, ''])

            return redirect(conn, '/filter/' + tool)
        else:
            get_sub = 0
            stat = 'disabled' if await acl_check('', 'owner_auth', '', '') == 1 else ''
            name = name if name else ''

            if tool in ('inter_wiki', 'outer_link'):
                value = ['', '', '']
                if name != '':
                    curs.execute(db_change("select html, plus, plus_t from html_filter where html = ? and kind = ?"), [name, tool])
                    exist = curs.fetchall()
                    value = exist[0] if exist else value

                select = ''
                if tool == 'inter_wiki':
                    ex = 'https://namu.wiki/w/'

                    select = ['', '']
                    curs.execute(db_change("select plus_t from html_filter where kind = 'inter_wiki_sub' and html = ?"), [name])
                    db_data = curs.fetchall()
                    if db_data and db_data[0][0] == 'under_bar':
                        select = ['', 'selected']

                    select = '''
                        <hr class="main_hr">
                        ''' + await get_lang('inter_wiki_space_change') + '''
                        <hr class="main_hr">
                        <select name="inter_type">
                            <option ''' + select[0] + ''' value="url_encode">%20</option>
                            <option ''' + select[1] + ''' value="under_bar">_</option>
                        </select>
                    '''
                else:
                    ex = 'youtube.com'

                title = await get_lang('interwiki_add') if tool == 'inter_wiki' else await get_lang('outer_link_add')
                form_data = '''
                    ''' + await get_lang('name') + '''
                    <hr class="main_hr">
                    <input class="__ON_INPUT__" value="''' + html.escape(value[0]) + '''" type="text" name="title">
                    <hr class="main_hr">
                    ''' + await get_lang('link') + ''' (EX : ''' + ex + ''')
                    <hr class="main_hr">
                    <input class="__ON_INPUT__" value="''' + html.escape(value[1]) + '''" type="text" name="link">
                    <hr class="main_hr">
                    ''' + await get_lang('icon') + ''' (''' + ('HTML' if tool == 'inter_wiki' else await get_lang('html_or_link')) + ''') (''' + await get_lang('link') + ' - EX' + ''' : /image/Test.svg)
                    <hr class="main_hr">
                    <input class="__ON_INPUT__" value="''' + html.escape(value[2]) + '''" type="text" name="icon">
                    ''' + select + '''
                '''
            elif tool == 'edit_filter':            
                curs.execute(db_change("select plus, plus_t from html_filter where html = ? and kind = 'regex_filter'"), [name])
                exist = curs.fetchall()
                if exist:
                    textarea = exist[0][0]
                    time_data = '' if exist[0][1] == 'X' else exist[0][1]
                    if time_data != '':
                        time_data = re.sub(r'[^0-9]', '', time_data)
                        time_data = str(int(int(number_check(time_data)) / (24 * 60 * 60)))
                else:
                    textarea = ''
                    time_data = ''

                title = await get_lang('edit_filter_add')
                form_data = '''
                    <hr class="main_hr">
                    <input class="__ON_INPUT__" placeholder="''' + await get_lang('day') + '''" name="day" type="text" value="''' + html.escape(time_data) + '''">
                    <hr class="main_hr">
                    <input class="__ON_INPUT__" placeholder="''' + await get_lang('regex') + '''" name="content" value="''' + html.escape(textarea) + '''" type="text">
                '''
            elif tool == 'name_filter':
                title = await get_lang('id_filter_add')
                form_data = '' + \
                    await get_lang('regex') + \
                    '<hr class="main_hr">' + \
                    '<input class="__ON_INPUT__" value="' + html.escape(name) + '" type="text" name="title">' + \
                ''
            elif tool == 'file_filter':
                title = await get_lang('file_filter_add')
                form_data = '' + \
                    await get_lang('regex') + \
                    '<hr class="main_hr">' + \
                    '<input class="__ON_INPUT__" value="' + html.escape(name) + '" type="text" name="title">' + \
                ''
            elif tool == 'email_filter':
                title = await get_lang('email_filter_add')
                form_data = '' + \
                    await get_lang('email') + \
                    '<hr class="main_hr">' + \
                    '<input class="__ON_INPUT__" value="' + html.escape(name) + '" type="text" name="title">' + \
                ''
            elif tool == 'image_license':
                title = await get_lang('image_license_add')
                form_data = '' + \
                    await get_lang('license') + \
                    '<hr class="main_hr">' + \
                    '<input class="__ON_INPUT__" value="' + html.escape(name) + '" type="text" name="title">' + \
                ''
            elif tool == 'extension_filter':
                title = await get_lang('extension_filter_add')
                form_data = '' + \
                    await get_lang('extension') + \
                    '<hr class="main_hr">' + \
                    '<input class="__ON_INPUT__" value="' + html.escape(name) + '" type="text" name="title">' + \
                    '<hr class="main_hr">' + \
                    await get_lang('max_file_size') + ' (MB) (' + await get_lang('default') + ' : ' + await get_lang('empty') + ')' + \
                    '<hr class="main_hr">' + \
                    '<input class="__ON_INPUT__" value="" type="text" name="max_file_size">' + \
                ''
            elif tool == 'document':
                acl_list = await get_acl_list()
                
                curs.execute(db_change("select plus, plus_t from html_filter where html = ? and kind = 'document'"), [name])
                db_data = curs.fetchall()
                acl_list = [['selected' if db_data and db_data[0][1] == for_a else '', for_a] for for_a in acl_list]

                title = await get_lang('document_filter_add')
                form_data = '''
                    ''' + await get_lang('name') + '''
                    <hr class="main_hr">
                    <input class="__ON_INPUT__" value="''' + html.escape(name) + '''" type="text" name="name">
                    <hr class="main_hr">
                    ''' + await get_lang('regex') + '''
                    <hr class="main_hr">
                    <input class="__ON_INPUT__" value="''' + (html.escape(db_data[0][0]) if db_data else '') + '''" type="text" name="regex">
                    <hr class="main_hr">
                    <a href="/acl/Test#exp">''' + await get_lang('acl') + '''</a>
                    <hr class="main_hr">
                    <select name="acl">
                        ''' + ''.join(['<option ' + for_a[0] + ' value=' + for_a[1] + '>' + ('normal' if for_a[1] == '' else for_a[1]) + '</option>' for for_a in acl_list]) + '''
                    </select>
                '''
            elif tool == 'template':
                title = await get_lang('template_document_add')

                value = ''
                if name:
                    curs.execute(db_change("select plus from html_filter where html = ? and kind = 'template'"), [name])
                    exist = curs.fetchall()
                    value = exist[0][0] if exist else '' 

                form_data = '' + \
                    await get_lang('template') + \
                    '<hr class="main_hr">' + \
                    '<input class="__ON_INPUT__" value="' + html.escape(name) + '" type="text" name="title">' + \
                    '<hr class="main_hr">' + \
                    await get_lang('explanation') + \
                    '<hr class="main_hr">' + \
                    '<input class="__ON_INPUT__" value="' + html.escape(value) + '" type="text" name="exp">' + \
                    '<hr class="main_hr">' + \
                ''
            else:
                title = await get_lang('edit_tool_add')
                
                value = ''
                if name:
                    curs.execute(db_change("select plus from html_filter where html = ? and kind = 'edit_top'"), [name])
                    exist = curs.fetchall()
                    value = exist[0][0] if exist else ''    

                form_data = '''
                    ''' + await get_lang('title') + '''
                    <hr class="main_hr">
                    <input class="__ON_INPUT__" value="''' + html.escape(name) + '''" type="text" name="title">
                    <hr class="main_hr">
                    ''' + await get_lang('markup') + '''
                    <hr class="main_hr">
                    <input class="__ON_INPUT__" value="''' + html.escape(value) + '''" type="text" name="markup">
                '''

            return await render_template(
                title,
                '''
                        <form method="post">
                            ''' + form_data + '''
                            <hr class="main_hr">
                            <button class="__ON_BUTTON__" ''' + stat + ''' type="submit">''' + await get_lang('add') + '''</button>
                        </form>
                        ''',
                get_sub,
                [['filter/' + tool, await get_lang('return')]]
            )
