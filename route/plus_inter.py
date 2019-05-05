from .tool.func import *

def plus_inter_2(conn, tools, name):
    curs = conn.cursor()
    
    if flask.request.method == 'POST':
        if tools == 'plus_inter_wiki':
            curs.execute('insert into inter (title, link) values (?, ?)', [flask.request.form.get('title', None), flask.request.form.get('link', None)])
            admin_check(None, 'inter_wiki_plus')
        elif tools == 'plus_edit_filter':
            if admin_check(1, 'edit_filter edit') != 1:
                return re_error('/error/3')

            if flask.request.form.get('limitless', '') != '':
                end = 'X'
            else:
                end = flask.request.form.get('second', 'X')

            try:
                re.compile(flask.request.form.get('content', 'test'))

                curs.execute("select name from filter where name = ?", [name])
                if curs.fetchall():
                    curs.execute("update filter set regex = ?, sub = ? where name = ?", [flask.request.form.get('content', 'test'), end, name])
                else:
                    curs.execute("insert into filter (name, regex, sub) values (?, ?, ?)", [name, flask.request.form.get('content', 'test'), end])
            except:
                return re_error('/error/23')                
        else:
            if tools == 'plus_name_filter':
                admin_check(None, 'name_filter edit')
                type_d = 'name'
            else:
                admin_check(None, 'email_filter edit')
                type_d = 'email'
            
            curs.execute('insert into html_filter (html, kind) values (?, ?)', [flask.request.form.get('title', 'test'), type_d])
        
        conn.commit()
    
        return redirect('/' + re.sub('^plus_', '', tools))
    else:
        if admin_check(1) != 1:
            stat = 'disabled'
        else:
            stat = ''

        if tools == 'plus_inter_wiki':
            title = load_lang('interwiki_add')
            form_data = '''
                        <input placeholder="''' + load_lang('name') + '''" type="text" name="title">
                        <hr class=\"main_hr\">
                        <input placeholder="link" type="text" name="link">
                        '''
        elif tools == 'plus_edit_filter':
            curs.execute("select regex, sub from filter where name = ?", [name])
            exist = curs.fetchall()
            if exist:
                textarea = exist[0][0]
                
                if exist[0][1] == 'X':
                    time_check = 'checked="checked"'
                    time_data = ''
                else:
                    time_check = ''
                    time_data = exist[0][1]
            else:
                textarea = ''
                time_check = ''
                time_data = ''

            title = load_lang('edit_filter_add')
            form_data = '''
                        <input placeholder="''' + load_lang('second') + '''" name="second" type="text" value="''' + html.escape(time_data) + '''">
                        <hr class=\"main_hr\">
                        <input ''' + stat + ''' type="checkbox" ''' + time_check + ''' name="limitless"> ''' + load_lang('limitless') + '''
                        <hr class=\"main_hr\">
                        <input ''' + stat + ''' placeholder="''' + load_lang('regex') + '''" name="content" value="''' + html.escape(textarea) + '''" type="text">
                        '''
        elif tools == 'plus_name_filter':
            title = load_lang('id_filter_add')
            form_data = '<input placeholder="' + load_lang('id') + ' ' + load_lang('regex') + '" type="text" name="title">'
        else:
            title = load_lang('email_filter_add')
            form_data = '<input placeholder="email" type="text" name="title">'

        return easy_minify(flask.render_template(skin_check(), 
            imp = [title, wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        ''' + form_data + '''
                        <hr class=\"main_hr\">
                        <button ''' + stat + ''' type="submit">''' + load_lang('add') + '''</button>
                    </form>
                    ''',
            menu = [[re.sub('^plus_', '', tools), load_lang('return')]]
        ))