from .tool.func import *

class run_count_section:
    def __init__(self, key, change):
        self.counter = key
        self.change = change

    def __call__(self, match):
        self.counter -= 1

        if self.counter == 0:
            return '\n' + self.change + '\n'
        else:
            return '\n' + match[1]

def edit_2(conn, name):
    curs = conn.cursor()

    ip = ip_check()
    section = flask.request.args.get('section', None)
    if section:
        curs.execute(db_change("select data from other where name = 'markup'"))
        markup = curs.fetchall()
        if markup[0][0] == 'namumark':
            section = int(number_check(section))
        else:
            return redirect('/edit/' + url_pas(name))

    if acl_check(name) == 1:
        return re_error('/ban')
    
    curs.execute(db_change("select id from history where title = ? order by id + 0 desc"), [name])
    doc_ver = curs.fetchall()
    doc_ver = doc_ver[0][0] if doc_ver else '0'

    edit_repeat = 0
    if flask.request.method == 'POST':
        edit_repeat = 1
        if flask.request.form.get('ver', '') != doc_ver:
            edit_repeat = 2
    
    if edit_repeat == 1:
        if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        if slow_edit_check() == 1:
            return re_error('/error/24')

        today = get_time()
        content = flask.request.form.get('content', '').replace('\r\n', '\n')
        
        if edit_filter_do(content) == 1:
            return re_error('/error/21')
            
        curs.execute(db_change('select data from other where name = "copyright_checkbox_text"'))
        copyright_checkbox_text_d = curs.fetchall()
        if copyright_checkbox_text_d and copyright_checkbox_text_d[0][0] != '' and flask.request.form.get('copyright_agreement', '') != 'yes':
            return re_error('/error/29')
        
        curs.execute(db_change("select data from data where title = ?"), [name])
        old = curs.fetchall()
        if old:  
            o_data = old[0][0].replace('\r\n', '\n')

            if section:
                run_count = run_count_section(section, content)

                c_data = html.escape('\n' + o_data)
                c_data = re.sub(r'\n(?P<in>={1,6})', '<br>\g<in>', c_data)
                c_data = re.sub(r'<br>((?:(?:(?!<br>).)*\n*)*)', run_count, c_data)
                c_data = re.sub(r'^\n', '', c_data)
                c_data = html.unescape(c_data)

                content = c_data

            leng = leng_check(len(o_data), len(content))
            
            curs.execute(db_change("update data set data = ? where title = ?"), [content, name])
        else:
            leng = '+' + str(len(content))

            curs.execute(db_change("insert into data (title, data) values (?, ?)"), [name, content])

            curs.execute(db_change('select data from other where name = "count_all_title"'))
            curs.execute(db_change("update other set data = ? where name = 'count_all_title'"), [str(int(curs.fetchall()[0][0]) + 1)])

        curs.execute(db_change("select user from scan where title = ? and type = ''"), [name])
        for scan_user in curs.fetchall():
            curs.execute(db_change("insert into alarm (name, data, date) values (?, ?, ?)"), [
                scan_user[0],
                ip + ' | <a href="/w/' + url_pas(name) + '">' + name + '</a> | Edit', 
                today
            ])
                
        history_plus(
            name,
            content,
            today,
            ip,
            flask.request.form.get('send', ''),
            leng
        )
        
        curs.execute(db_change("delete from back where link = ?"), [name])
        curs.execute(db_change("delete from back where title = ? and type = 'no'"), [name])
        
        render_set(
            title = name,
            data = content,
            num = 1
        )
        
        conn.commit()
        
        return redirect('/w/' + url_pas(name) + (('#edit_load_' + str(section)) if section else ''))
    else:
        editor_top_text = ''
        if edit_repeat != 2:
            load_title = flask.request.args.get('plus', None)
            if load_title:
                curs.execute(db_change("select data from data where title = ?"), [load_title])
                get_data = curs.fetchall()
                data = get_data[0][0] if get_data else ''
            else:
                curs.execute(db_change("select data, id from history where title = ? order by id + 0 desc"), [name])
                old = curs.fetchall()
                old = old if old else [['']]
                if section:
                    data = html.escape('\n' + old[0][0].replace('\r\n', '\n'))
                    data = re.sub(r'\n(?P<in>={1,6})', '<br>\g<in>', data)

                    section_data = re.findall(r'<br>((?:(?:(?!<br>).)*\n*)*)', data)
                    if len(section_data) >= section:
                        data = html.unescape(section_data[section - 1])
                    else:
                        return redirect('/edit/' + url_pas(name))
                else:
                    data = old[0][0].replace('\r\n', '\n')
                    editor_top_text += '<a href="/manager/15?plus=' + url_pas(name) + '">(' + load_lang('load') + ')</a> '
        else:
            data = flask.request.form.get('content', '')
            warring_edit = load_lang('exp_edit_conflict') + ' '

            if flask.request.form.get('ver', '0') == '0':
                warring_edit += '<a href="/raw/' + url_pas(name) + '">(r' + doc_ver + ')</a>'
            else:
                warring_edit += '' + \
                    '<a href="/diff/' + url_pas(name) + '?first=' + flask.request.form.get('ver', '1') + '&second=' + doc_ver + '">(r' + doc_ver + ')</a>' + \
                ''

            warring_edit += '<hr class="main_hr">'
            editor_top_text = warring_edit + editor_top_text

        editor_top_text += '' + \
            '<a href="/edit_filter">(' + load_lang('edit_filter_rule') + ')</a>' + \
            '<hr class="main_hr">' + \
        ''

        curs.execute(db_change('select data from other where name = "edit_bottom_text"'))
        sql_d = curs.fetchall()
        b_text = ('<hr class="main_hr">' + sql_d[0][0]) if sql_d and sql_d[0][0] != '' else ''
        
        curs.execute(db_change('select data from other where name = "copyright_checkbox_text"'))
        sql_d = curs.fetchall()
        if sql_d and sql_d[0][0] != '':
            cccb_text = '' + \
                '<hr class="main_hr">' + \
                '<input type="checkbox" name="copyright_agreement" value="yes"> ' + sql_d[0][0] + \
                '<hr class="main_hr">' + \
            ''
        else:
            cccb_text = ''

        curs.execute(db_change('select data from other where name = "edit_help"'))
        sql_d = curs.fetchall()
        p_text = html.escape(sql_d[0][0]) if sql_d and sql_d[0][0] != '' else load_lang('default_edit_help')

        data = re.sub(r'\n+$', '', data)

        if flask.request.cookies.get('main_css_monaco', '0') == '1':
            editor_display = 'style="display: none;"'
            monaco_display = ''
            add_get_file = '''
                <link   rel="stylesheet"
                        data-name="vs/editor/editor.main" 
                        href="https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.20.0/min/vs/editor/editor.main.min.css">
                <script src="https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.20.0/min/vs/loader.min.js"></script>
            '''
            add_script = '''
                require.config({ paths: { 'vs': 'https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.20.0/min/vs' }});
                require(["vs/editor/editor.main"], function () {
                    window.editor = monaco.editor.create(document.getElementById('monaco_editor'), {
                        value: document.getElementById('content').value,
                        language: 'plaintext',
                        theme: 'vs-dark'
                    });
                });

                function monaco_to_content() {
                    document.getElementById('content').innerHTML = window.editor.getValue();
                }
            '''
        else:
            editor_display = ''
            monaco_display = 'style="display: none;"'
            add_get_file = ''
            add_script = 'function monaco_to_content() {}'

        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2(['(' + load_lang('edit') + ')', 0])],
            data =  editor_top_text + add_get_file + '''
                <form method="post">
                    <script>
                        do_paste_image();
                        do_not_out();
                        ''' + add_script + '''
                    </script>
                    <div ''' + editor_display + '''>''' + edit_button() + '''</div>
                    <div    id="monaco_editor"
                            class="content" 
                            ''' + monaco_display + '''></div>
                    <textarea   id="content"
                                ''' + editor_display + '''
                                class="content" 
                                placeholder="''' + p_text + '''" 
                                name="content">''' + html.escape(data) + '''</textarea>
                    <hr class="main_hr">
                    <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                    <input id="origin" name="ver" value="''' + doc_ver + '''">
                    <hr class="main_hr">
                    ''' + captcha_get() + ip_warring() + cccb_text + '''
                    <button id="save"
                            type="submit" 
                            onclick="monaco_to_content(); save_stop_exit();">''' + load_lang('save') + '''</button>
                    <button id="preview" 
                            type="button" 
                            onclick="monaco_to_content(); load_preview(\'''' + url_pas(name) + '\');">' + load_lang('preview') + '''</button>
                </form>
                ''' + b_text + '''
                <hr class="main_hr">
                <div id="see_preview"></div>
            ''',
            menu = [
                ['w/' + url_pas(name), load_lang('return')],
                ['delete/' + url_pas(name), load_lang('delete')], 
                ['move/' + url_pas(name), load_lang('move')], 
                ['upload', load_lang('upload')]
            ]
        ))