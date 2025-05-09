import multiprocessing

from .tool.func import *

from .view_set import view_set_markup

def edit_render_set(name, content):
    with get_db_connect() as conn:
        render_set(conn, 
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
        
async def edit_editor(conn, ip, data_main = '', do_type = 'edit', addon = '', name = ''):
    curs = conn.cursor()

    monaco_editor_top = ''
    div = ''

    if do_type == 'edit':
        curs.execute(db_change('select data from other where name = "edit_help"'))
        sql_d = curs.fetchall()

        curs.execute(db_change("select set_data from data_set where doc_name = ? and set_name = 'document_top'"), [name])
        body = curs.fetchall()
        div = body[0][0] if body else ''
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
            
    p_text = html.escape(sql_d[0][0]) if sql_d and sql_d[0][0] != '' else get_lang(conn, 'default_edit_help')
    
    monaco_editor_top += '<a href="javascript:opennamu_do_editor_temp_save();">(' + get_lang(conn, 'load_temp_save') + ')</a> <a href="javascript:opennamu_do_editor_temp_save_load();">(' + get_lang(conn, 'load_temp_save_load') + ')</a>'
    monaco_editor_top += '<hr class="main_hr">'
    
    darkmode = flask.request.cookies.get('main_css_darkmode', '0')
    monaco_thema = 'vs-dark' if darkmode == '1' else ''
    
    monaco_on = get_main_skin_set(conn, flask.session, 'main_css_monaco', ip)
    editor_display = ['style="display: none;"' for _ in range(3)]
    if monaco_on == 'use':
        editor_display[1] = ''
    else:
        editor_display[0] = ''

    # 에디터 선택창
    monaco_editor_top += '<select onclick="do_sync_monaco_and_textarea();" id="opennamu_select_editor" onchange="opennamu_edit_turn_off_monaco();">'
    monaco_editor_top += '<option value="default" ' + ('selected' if editor_display[0] == '' else '') + '>' + get_lang(conn, 'default') + '</option>'
    monaco_editor_top += '<option value="monaco" ' + ('selected' if editor_display[1] == '' else '') + '>' + get_lang(conn, 'monaco_editor') + '</option>'
    monaco_editor_top += '</select> '

    # 문법 선택창
    if do_type == 'edit':
        monaco_editor_top += view_set_markup(conn, document_name = name, addon = 'id="opennamu_editor_markup" onclick="opennamu_do_sync_monaco_markup();"')
    else:
        monaco_editor_top += view_set_markup(conn, addon = 'id="opennamu_editor_markup" onclick="opennamu_do_sync_monaco_markup();"', disable = 'disabled')

    textarea_size = 'opennamu_textarea_500' if do_type == 'edit' else 'opennamu_textarea_100'

    out_field = await captcha_get(conn) + ip_warning(conn) + addon
    if out_field != '':
        out_field += '<hr class="main_hr">'

    return '''
        <textarea style="display: none;" id="opennamu_edit_origin" name="doc_data_org">''' + html.escape(data_main) + '''</textarea>
        <div>
            ''' + monaco_editor_top + '''
            <hr class="main_hr">
            ''' + edit_button(conn) + '''
            <div id="opennamu_editor_user_button"></div>
        </div>
        
        ''' + div + '''

        <div id="opennamu_monaco_editor" class="''' + textarea_size + '''" ''' + editor_display[1] + '''></div>
        <textarea id="opennamu_edit_textarea" class="''' + textarea_size + '''" ''' + editor_display[0] + ''' name="content" placeholder="''' + p_text + '''">''' + html.escape(data_main) + '''</textarea>
        <hr class="main_hr">
        ''' + out_field + '''
        
        <script>
            window.addEventListener('DOMContentLoaded', function() {
                do_stop_exit();
                do_paste_image();
                do_monaco_init("''' + monaco_thema + '''");
                opennnamu_do_user_editor();
            });
        </script>
                        
        <button id="opennamu_save_button" type="submit" onclick="do_stop_exit_release();">''' + get_lang(conn, 'send') + '''</button>
        <button id="opennamu_preview_button" type="button" onclick="opennamu_do_editor_preview();">''' + get_lang(conn, 'preview') + '''</button>
        <hr class="main_hr">

        <div id="opennamu_preview_area"></div>
    '''

async def edit(name = 'Test', section = 0, do_type = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()
    
        ip = ip_check()

        edit_req_mode = 0
        if await acl_check(name, 'document_edit') == 1:
            edit_req_mode = 1
            if await acl_check(name, 'document_edit_request') == 1:
                return redirect(conn, '/raw_acl/' + url_pas(name))
            
        if do_title_length_check(conn, name) == 1:
            return await re_error(conn, 38)
        
        curs.execute(db_change("select id from history where title = ? order by id + 0 desc"), [name])
        doc_ver = curs.fetchall()
        doc_ver = doc_ver[0][0] if doc_ver else '0'

        if doc_ver == '0':
            if await acl_check(name, 'document_make_acl') == 1:
                edit_req_mode = 1

        curs.execute(db_change("select set_data from data_set where doc_name = ? and doc_rev = ? and set_name = 'edit_request_data'"), [name, doc_ver])
        if curs.fetchall():
            return redirect(conn, '/edit_request_from/' + url_pas(name))
        
        section = '' if section == 0 else section
        post_ver = flask.request.form.get('ver', '')
        if flask.request.method == 'POST':
            edit_repeat = 'error' if post_ver != doc_ver else 'post'
        else:
            edit_repeat = 'get'
        
        if edit_repeat == 'post':
            if await captcha_post(conn, flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return await re_error(conn, 13)
    
            if await do_edit_slow_check(conn) == 1:
                return await re_error(conn, 24)
    
            today = get_time()
            content = flask.request.form.get('content', '').replace('\r', '')
            send = flask.request.form.get('send', '')
            agree = flask.request.form.get('copyright_agreement', '')
            
            if await do_edit_filter(conn, content) == 1:
                return await re_error(conn, 21)
            
            if await do_edit_filter(conn, send) == 1:
                return await re_error(conn, 21)

            if await do_edit_send_check(conn, send) == 1:
                return await re_error(conn, 37)

            if do_edit_text_bottom_check_box_check(conn, agree) == 1:
                return await re_error(conn, 29)
            
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
                    return await re_error(conn, 44)

            curs.execute(db_change("select data from other where name = 'edit_timeout'"))
            db_data_2 = curs.fetchall()
            db_data_2 = number_check(db_data_2[0][0]) if db_data_2 and db_data_2[0][0] != '' else ''
            if db_data_2 != '' and platform.system() in ('Linux', 'Darwin'):
                timeout = edit_timeout(edit_render_set, (name, content), timeout = int(db_data_2))
            else:
                timeout = 0

            if timeout == 1:
                return await re_error(conn, 41)
            
            if edit_req_mode == 0:
                # 진짜 기록 부분
                curs.execute(db_change("delete from data where title = ?"), [name])
                curs.execute(db_change("insert into data (title, data) values (?, ?)"), [name, content])
        
                curs.execute(db_change("select id from user_set where name = 'watchlist' and data = ?"), [name])
                for scan_user in curs.fetchall():
                    await add_alarm(scan_user[0], ip, '<a href="/w/' + url_pas(name) + '">' + html.escape(name) + '</a>')
                        
                history_plus(conn, 
                    name,
                    content,
                    today,
                    ip,
                    send,
                    leng
                )
                
                render_set(conn, 
                    doc_name = name,
                    doc_data = content,
                    data_type = 'backlink'
                )
                
                section = (('#edit_load_' + str(section)) if section != '' else '')
                return redirect(conn, '/w/' + url_pas(name) + section)
            else:
                curs.execute(db_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, 'edit_request_data', ?)"), [name, doc_ver, content])
                curs.execute(db_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, 'edit_request_user', ?)"), [name, doc_ver, ip])
                curs.execute(db_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, 'edit_request_date', ?)"), [name, doc_ver, today])
                curs.execute(db_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, 'edit_request_send', ?)"), [name, doc_ver, send])
                curs.execute(db_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, 'edit_request_leng', ?)"), [name, doc_ver, leng])
                curs.execute(db_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, 'edit_request_doing', ?)"), [name, doc_ver, today])

                curs.execute(db_change("select id from user_set where name = 'watchlist' and data = ?"), [name])
                for scan_user in curs.fetchall():
                    await add_alarm(scan_user[0], ip, '<a href="/edit_request/' + url_pas(name) + '">' + html.escape(name) + '</a> edit_request')
            
                return redirect(conn, '/edit_request_from/' + url_pas(name))
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
                    editor_top_text += '<a href="/manager/15/' + url_pas(name) + '">(' + get_lang(conn, 'load') + ')</a> '
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

                warning_edit = get_lang(conn, 'exp_edit_conflict') + ' '
    
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
    
            editor_top_text += '<a href="/filter/edit_filter">(' + get_lang(conn, 'edit_filter_rule') + ')</a>'
    
            if editor_top_text != '':
                editor_top_text += '<hr class="main_hr">'

            sub_menu = ' (' + str(section) + ')' if section != '' else ''
            sub_title = '(' + get_lang(conn, 'edit_request') + ')' if edit_req_mode == 1 else '(' + get_lang(conn, 'edit') + ')'

            return easy_minify(conn, flask.render_template(skin_check(conn), 
                imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css([sub_title + sub_menu, 0])],
                data = editor_top_text + '''
                    <form method="post">
                        <textarea style="display: none;" name="doc_section_data_where">''' + data_section_where + '''</textarea>
                        <input style="display: none;" name="doc_section_edit_apply" value="''' + doc_section_edit_apply + '''">

                        <input style="display: none;" id="opennamu_editor_doc_name" value="''' + html.escape(name) + '''">
                        <input style="display: none;" name="ver" value="''' + doc_ver + '''">
                        
                        <input placeholder="''' + get_lang(conn, 'why') + '''" name="send">
                        <hr class="main_hr">
                        
                        ''' + await edit_editor(conn, ip, data_section, addon = get_edit_text_bottom_check_box(conn) + get_edit_text_bottom(conn, 'edit') , name = name) + '''
                    </form>
                ''',
                menu = [
                    ['w/' + url_pas(name), get_lang(conn, 'return')],
                    ['delete/' + url_pas(name), get_lang(conn, 'delete')], 
                    ['move/' + url_pas(name), get_lang(conn, 'move')], 
                    ['upload', get_lang(conn, 'upload')]
                ]
            ))