import multiprocessing

from .tool.func import *


def edit_render_set(name, content):
    render_set(
        doc_name = name,
        doc_data = content
    )

# https://stackoverflow.com/questions/13821156/timeout-function-using-threading-in-python-does-not-work
def edit_timeout(func, args = (), timeout = 3):
    pool = multiprocessing.Pool(processes = 1)
    result = pool.apply_async(func, args = args)
    try:
        result.get(timeout = timeout)
    except multiprocessing.TimeoutError:
        pool.terminate()
        return 1
    else:
        pool.close()
        pool.join()
        return 0
        
def edit_editor(curs, ip, data_main = '', do_type = 'edit', addon = ''):
    monaco_editor_top = ''
    editor_display = ''
    add_get_file = ''
    monaco_display = ''

    if do_type == 'edit':
        curs.execute(db_change('select data from other where name = "edit_help"'))
        sql_d = curs.fetchall()
    elif do_type == 'bbs':
        curs.execute(db_change('select data from other where name = "bbs_help"'))
        sql_d = curs.fetchall()
    elif do_type == 'bbs_comment':
        curs.execute(db_change('select data from other where name = "bbs_comment_help"'))
        sql_d = curs.fetchall()
    else:
        curs.execute(db_change('select data from other where name = "topic_text"'))
        sql_d = curs.fetchall()

    if do_type == 'bbs_comment':
        do_type = 'thread'
    elif do_type == 'bbs':
        do_type = 'edit'
            
    p_text = html.escape(sql_d[0][0]) if sql_d and sql_d[0][0] != '' else load_lang('default_edit_help')
    
    monaco_editor_top += '<a href="javascript:opennamu_do_editor_temp_save();">(' + load_lang('load_temp_save') + ')</a> <a href="javascript:opennamu_do_editor_temp_save_load();">(' + load_lang('load_temp_save_load') + ')</a> '
    monaco_editor_top += '<a href="javascript:opennamu_edit_turn_off_monaco();">(' + load_lang('turn_off_monaco') + ')</a>'

    add_get_file = '''
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.41.0/min/vs/editor/editor.main.min.css" integrity="sha512-MFDhxgOYIqLdcYTXw7en/n5BshKoduTitYmX8TkQ+iJOGjrWusRi8+KmfZOrgaDrCjZSotH2d1U1e/Z1KT6nWw==" crossorigin="anonymous" referrerpolicy="no-referrer" />
        <script src="https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.41.0/min/vs/loader.min.js" integrity="sha512-A+6SvPGkIN9Rf0mUXmW4xh7rDvALXf/f0VtOUiHlDUSPknu2kcfz1KzLpOJyL2pO+nZS13hhIjLqVgiQExLJrw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
    '''
    
    darkmode = flask.request.cookies.get('main_css_darkmode', '0')
    monaco_thema = 'vs-dark' if darkmode == '1' else ''
    
    add_script = 'do_monaco_init("' + monaco_thema + '");'
    
    monaco_on = get_main_skin_set(curs, flask.session, 'main_css_monaco', ip)
    if monaco_on == 'use':
        editor_display = 'style="display: none;"'
    else:
        monaco_display = 'style="display: none;"'

    textarea_size = 'opennamu_textarea_500' if do_type == 'edit' else 'opennamu_textarea_100'

    return add_get_file + '''
        <textarea style="display: none;" id="opennamu_edit_origin" name="doc_data_org">''' + html.escape(data_main) + '''</textarea>
        <div>
            ''' + monaco_editor_top + '''
            <hr class="main_hr">
            ''' + edit_button() + '''
        </div>
        
        <div id="opennamu_monaco_editor" class="''' + textarea_size + '''" ''' + monaco_display + '''></div>
        <textarea id="opennamu_edit_textarea" ''' + editor_display + ''' class="''' + textarea_size + '''" name="content" placeholder="''' + p_text + '''">''' + html.escape(data_main) + '''</textarea>
        <hr class="main_hr">
        
        ''' + captcha_get() + ip_warning() + addon + '''
        <hr class="main_hr">

        <script>
            do_stop_exit();
            do_paste_image();
            ''' + add_script + '''
        </script>
                        
        <button id="opennamu_save_button" type="submit" onclick="do_stop_exit_release();">''' + load_lang('send') + '''</button>
        <button id="opennamu_preview_button" type="button" onclick="opennamu_do_editor_preview();">''' + load_lang('preview') + '''</button>
        <hr class="main_hr">

        <div id="opennamu_preview_area"></div>
    '''

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

            curs.execute(db_change("select data from other where name = 'document_content_max_length'"))
            db_data_3 = curs.fetchall()
            if db_data_3 and db_data_3[0][0] != '':
                if int(number_check(db_data_3[0][0])) < len(content):
                    return re_error('/error/44')

            curs.execute(db_change("select data from other where name = 'edit_timeout'"))
            db_data_2 = curs.fetchall()
            db_data_2 = '' if not db_data_2 else number_check(db_data_2[0][0])

            if db_data_2 != '' and platform.system() == 'Linux':
                timeout = edit_timeout(edit_render_set, (name, content), timeout = int(db_data_2))
            else:
                timeout = 0

            if timeout == 1:
                return re_error('/error/41')
            
            if db_data:
                curs.execute(db_change("update data set data = ? where title = ?"), [content, name])
            else:    
                curs.execute(db_change("insert into data (title, data) values (?, ?)"), [name, content])
    
            curs.execute(db_change("select user from scan where title = ? and type = ''"), [name])
            for scan_user in curs.fetchall():
                add_alarm(scan_user[0], ip, '<a href="/w/' + url_pas(name) + '">' + html.escape(name) + '</a>')
                    
            history_plus(
                name,
                content,
                today,
                ip,
                send,
                leng
            )
            
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

            if data_section == '':
                data_section = data
    
            editor_top_text += '<a href="/filter/edit_filter">(' + load_lang('edit_filter_rule') + ')</a>'
    
            if editor_top_text != '':
                editor_top_text += '<hr class="main_hr">'

            sub_menu = ' (' + str(section) + ')' if section != '' else ''

            return easy_minify(flask.render_template(skin_check(), 
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('edit') + ')' + sub_menu, 0])],
                data = editor_top_text + '''
                    <form method="post">
                        <textarea style="display: none;" name="doc_section_data_where">''' + data_section_where + '''</textarea>
                        <input style="display: none;" name="doc_section_edit_apply" value="''' + doc_section_edit_apply + '''">

                        <input style="display: none;" id="opennamu_editor_doc_name" value="''' + html.escape(name) + '''">
                        <input style="display: none;" name="ver" value="''' + doc_ver + '''">
                        
                        <input placeholder="''' + load_lang('why') + '''" name="send">
                        <hr class="main_hr">
                        
                        ''' + edit_editor(curs, ip, data_section, addon = get_edit_text_bottom_check_box() + get_edit_text_bottom()) + '''
                    </form>
                ''',
                menu = [
                    ['w/' + url_pas(name), load_lang('return')],
                    ['delete/' + url_pas(name), load_lang('delete')], 
                    ['move/' + url_pas(name), load_lang('move')], 
                    ['upload', load_lang('upload')]
                ]
            ))