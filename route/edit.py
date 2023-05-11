from .tool.func import *

def edit(name = 'Test', section = 0, do_type = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()
    
        ip = ip_check()
        if acl_check(name, 'document_edit') == 1:
            return redirect('/raw_acl/' + url_pas(name))
        
        if do_title_length_check(name) == 1:
            return re_error('/error/38')
        
        curs.execute(db_change("select id from history where title = ? order by id + 0 desc"), [name])
        doc_ver = curs.fetchall()
        doc_ver = doc_ver[0][0] if doc_ver else '0'
        
        section = '' if section == 0 else section
        post_ver = flask.request.form.get('ver', '')
        if flask.request.method == 'POST':
            edit_repeat = 'error' if post_ver != doc_ver else 'post'
            edit_repeat = 'error' if do_type == 'preview' else 'post'
        else:
            edit_repeat = 'get'
        
        if edit_repeat == 'post':
            if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return re_error('/error/13')
            else:
                captcha_post('', 0)
    
            if do_edit_slow_check() == 1:
                return re_error('/error/24')
    
            today = get_time()
            content = flask.request.form.get('content', '').replace('\r', '')
            send = flask.request.form.get('send', '')
            agree = flask.request.form.get('copyright_agreement', '')
            
            if do_edit_filter(content) == 1:
                return re_error('/error/21')

            if do_edit_send_check(send) == 1:
                return re_error('/error/37')

            if do_edit_text_bottom_check_box_check(agree) == 1:
                return re_error('/error/29')
            
            curs.execute(db_change("select data from data where title = ?"), [name])
            db_data = curs.fetchall()
            if db_data:
                o_data = db_data[0][0].replace('\r', '')

                if section != '':
                    if flask.request.form.get('doc_section_edit_apply', 'X') != 'X':
                        if flask.request.form.get('doc_section_data_where', '') != '':
                            data_match_where = flask.request.form.get('doc_section_data_where', '').split(',')
                            if len(data_match_where) == 2:
                                data_match_a = int(number_check(data_match_where[0]))
                                if data_match_where[1] != 'inf':
                                    data_match_b = int(number_check(data_match_where[1]))
                                else:
                                    data_match_b = 'inf'

                                try:
                                    if data_match_b != 'inf':
                                        content = o_data[ : data_match_a] + content + o_data[data_match_b : ]
                                    else:
                                        content = o_data[ : data_match_a] + content
                                except:
                                    pass
    
                leng = leng_check(len(o_data), len(content))
            else:
                leng = '+' + str(len(content))

            render_set(
                doc_name = name,
                doc_data = content,
                data_in = ''
            )
                
            if db_data:
                curs.execute(db_change("update data set data = ? where title = ?"), [content, name])
            else:    
                curs.execute(db_change("insert into data (title, data) values (?, ?)"), [name, content])
    
                curs.execute(db_change('select data from other where name = "count_all_title"'))
                curs.execute(db_change("update other set data = ? where name = 'count_all_title'"), [str(int(curs.fetchall()[0][0]) + 1)])
    
            curs.execute(db_change("select user from scan where title = ? and type = ''"), [name])
            for scan_user in curs.fetchall():
                add_alarm(scan_user[0], ip + ' | <a href="/w/' + url_pas(name) + '">' + html.escape(name) + '</a> | Edit')
                    
            history_plus(
                name,
                content,
                today,
                ip,
                send,
                leng
            )
            
            curs.execute(db_change("delete from back where link = ?"), [name])
            curs.execute(db_change("delete from back where title = ? and type = 'no'"), [name])
            
            render_set(
                doc_name = name,
                doc_data = content,
                data_type = 'backlink'
            )
            
            conn.commit()
            
            section = (('#edit_load_' + str(section)) if section != '' else '')
            
            return redirect('/w/' + url_pas(name) + section)
        else:
            editor_top_text = ''

            doc_section_edit_apply = 'X'
            data_section = ''
            data_section_where = ''
            data_preview = ''

            if edit_repeat == 'get':
                if do_type == 'load':
                    if flask.session and 'edit_load_document' in flask.session:
                        load_title = flask.session['edit_load_document']
                    else:
                        load_title = 0
                else:
                    load_title = 0
                
                if load_title == 0 and section == '':
                    load_title = name
                    editor_top_text += '<a href="/manager/15/' + url_pas(name) + '">(' + load_lang('load') + ')</a> '
                elif section != '':
                    load_title = name
                    
                curs.execute(db_change("select data from data where title = ?"), [load_title])
                db_data = curs.fetchall()
                data = db_data[0][0] if db_data else ''
                data = data.replace('\r', '')

                if section != '':
                    curs.execute(db_change('select data from other where name = "markup"'))
                    db_data = curs.fetchall()
                    db_data = db_data[0][0] if db_data else 'namumark'
                    if db_data in ('namumark', 'namumark_beta'):
                        count = 1
                        data_section = '\n' + data + '\n'
                        
                        while 1:
                            data_match_re = r'\n((={1,6})(#?) ?([^\n]+))\n'
                            data_match = re.search(data_match_re, data_section)
                            if not data_match:
                                data_section = ''

                                break
                            elif count > section:
                                data_section = ''

                                break

                            if section == count:
                                data_section_sub = data_section
                                data_section_sub = re.sub(data_match_re, ('.' * (len(data_match.group(0)) - 1)) + '\n', data_section_sub, 1)

                                data_match_plus = re.search(data_match_re, data_section_sub)
                                if data_match_plus:
                                    data_section = data[data_match.span()[0] : data_match_plus.span()[0] - 1]
                                    data_section_where = str(data_match.span()[0]) + ',' + str(data_match_plus.span()[0] - 1)
                                else:
                                    data_section = data[data_match.span()[0] : ]
                                    data_section_where = str(data_match.span()[0]) + ',inf'

                                doc_section_edit_apply = 'O'

                                break
                            else:
                                data_section = re.sub(data_match_re, ('.' * (len(data_match.group(0)) - 1)) + '\n', data_section, 1)

                            count += 1
            else:
                data = flask.request.form.get('content', '')
                data = data.replace('\r', '')
                
                data_section_where = flask.request.form.get('doc_section_data_where', '')
                doc_section_edit_apply = flask.request.form.get('doc_section_edit_apply', '')

                doc_ver = flask.request.form.get('ver', '')

                if do_type != 'preview':
                    warning_edit = load_lang('exp_edit_conflict') + ' '
        
                    if flask.request.form.get('ver', '0') == '0':
                        warning_edit += '<a href="/raw/' + url_pas(name) + '">(r' + doc_ver + ')</a>'
                    else:
                        warning_edit += '' + \
                            '<a href="/diff/' + flask.request.form.get('ver', '1') + '/' + doc_ver + '/' + url_pas(name) + '">' + \
                                '(r' + doc_ver + ')' + \
                            '</a>' + \
                        ''
        
                    warning_edit += '<hr class="main_hr">'
                    editor_top_text = warning_edit + editor_top_text
                else:
                    data_preview = render_set(
                        doc_name = name, 
                        doc_data = data,
                        data_in = 'from'
                    )

            if data_section == '':
                data_section = data

            if section == '':
                form_action = 'formaction="/edit/' + url_pas(name) + '"'
                form_action_preview = 'formaction="/edit_preview/' + url_pas(name) + '"'
            else:
                form_action = 'formaction="/edit_section/' + str(section) + '/' + url_pas(name) + '"'
                form_action_preview = 'formaction="/edit_section_preview/' + str(section) + '/' + url_pas(name) + '"'
    
            editor_top_text += '<a href="/edit_filter">(' + load_lang('edit_filter_rule') + ')</a>'
    
            curs.execute(db_change('select data from other where name = "edit_help"'))
            sql_d = curs.fetchall()
            p_text = html.escape(sql_d[0][0]) if sql_d and sql_d[0][0] != '' else load_lang('default_edit_help')
            
            monaco_on = get_main_skin_set(curs, flask.session, 'main_css_monaco', ip)
            if monaco_on == 'use':
                editor_display = 'style="display: none;"'
                monaco_display = ''
                add_get_file = '''
                    <link   rel="stylesheet"
                            data-name="vs/editor/editor.main" 
                            href="https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.37.1/min/vs/editor/editor.main.min.css">
                    <script src="https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.37.1/min/vs/loader.min.js"></script>
                '''

                editor_top_text += ' <a href="javascript:opennamu_edit_turn_off_monaco();">(' + load_lang('turn_off_monaco') + ')</a>'
                
                if flask.request.cookies.get('main_css_darkmode', '0') == '1':
                    monaco_thema = 'vs-dark'
                else:
                    monaco_thema = ''
                
                add_script = 'do_monaco_init("' + monaco_thema + '");'
            else:
                editor_display = ''
                monaco_display = 'style="display: none;"'
                add_get_file = ''
                add_script = 'opennamu_edit_turn_off_monaco();'

            if editor_top_text != '':
                editor_top_text += '<hr class="main_hr">'

            sub_menu = ' (' + str(section) + ')' if section != '' else ''
    
            return easy_minify(flask.render_template(skin_check(), 
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('edit') + ')' + sub_menu, 0])],
                data =  editor_top_text + add_get_file + '''
                    <script>
                        
                    </script>
                    <form method="post">
                        <textarea style="display: none;" id="opennamu_edit_origin" name="doc_data_org">''' + html.escape(data_section) + '''</textarea>
                        <textarea style="display: none;" name="doc_section_data_where">''' + data_section_where + '''</textarea>
                        <input style="display: none;" name="doc_section_edit_apply" value="''' + doc_section_edit_apply + '''">

                        <input style="display: none;" name="ver" value="''' + doc_ver + '''">
                        
                        <div>''' + edit_button('opennamu_edit_textarea', 'opennamu_monaco_editor') + '''</div>
                        
                        <div id="opennamu_monaco_editor" class="opennamu_textarea_500" ''' + monaco_display + '''></div>
                        <textarea id="opennamu_edit_textarea" ''' + editor_display + ''' class="opennamu_textarea_500" name="content" placeholder="''' + p_text + '''">''' + html.escape(data_section) + '''</textarea>
                        <hr class="main_hr">

                        <input placeholder="''' + load_lang('why') + '''" name="send">
                        <hr class="main_hr">
                        
                        ''' + captcha_get() + ip_warning() + get_edit_text_bottom_check_box() + get_edit_text_bottom() + '''
                        
                        <button id="opennamu_save_button" type="submit" ''' + form_action + ''' onclick="do_monaco_to_textarea(); do_stop_exit_release();">''' + load_lang('save') + '''</button>
                        <button id="opennamu_preview_button" type="submit" ''' + form_action_preview + ''' onclick="do_monaco_to_textarea(); do_stop_exit_release();">''' + load_lang('preview') + '''</button>
                    </form>
                    
                    <hr class="main_hr">
                    <div id="opennamu_preview_area">''' + data_preview + '''</div>
                    
                    <script>
                        do_stop_exit();
                        do_paste_image('opennamu_edit_textarea', 'opennamu_monaco_editor');
                        ''' + add_script + '''
                    </script>
                ''',
                menu = [
                    ['w/' + url_pas(name), load_lang('return')],
                    ['delete/' + url_pas(name), load_lang('delete')], 
                    ['move/' + url_pas(name), load_lang('move')], 
                    ['upload', load_lang('upload')]
                ]
            ))