from .tool.func import *

def bbs_w_post_make_thread(user_id, date, data, code, color = '', blind = '', add_style = ''):
    if blind != '':
        if data == '':
            color_b = 'opennamu_comment_blind'
        else:
            color_b = 'opennamu_comment_blind_admin'
    else:
        color_b = 'opennamu_comment_blind_not'

    return '''
        <table class="opennamu_comment" style="''' + add_style + '''">
            <tr>
                <td class="opennamu_comment_color_''' + color + '''">
                    <a href="#thread_shortcut" id="''' + code + '''">#''' + code + '''</a>
                    ''' + user_id + '''
                    <span style="float: right;">''' + date + '''</span>
                </td>
            </tr>
            <tr>
                <td class="''' + color_b + '''" id="opennamu_comment_data_main">
                    ''' + data + '''
                </td>
            </tr>
        </table>
    '''

def bbs_w_post(bbs_num = '', post_num = '', do_type = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change('select set_name, set_data, set_code from bbs_data where set_id = ? and set_code = ?'), [bbs_num, post_num])
        db_data = curs.fetchall()
        if not db_data:
            return redirect('/bbs/main')

        bbs_num_str = str(bbs_num)
        post_num_str = str(post_num)
        bbs_comment_acl = acl_check(bbs_num_str, 'bbs_comment')
        ip = ip_check()
        
        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_type"'), [bbs_num])
        db_data_2 = curs.fetchall()
        if not db_data_2:
            return redirect('/bbs/main')
        elif db_data_2[0][0] == 'thread':
            if flask.request.method == 'POST' and do_type != 'preview':
                if bbs_comment_acl == 1:
                    return redirect('/bbs/set/' + bbs_num_str)
                
                if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                    return re_error('/error/13')
                else:
                    captcha_post('', 0)

                set_id = bbs_num_str + '-' + post_num_str

                curs.execute(db_change('select set_code from bbs_data where set_name = "comment" and set_id = ? order by set_code + 0 desc'), [set_id])
                db_data = curs.fetchall()
                id_data = str(int(db_data[0][0]) + 1) if db_data else '1'

                data = flask.request.form.get('content', '')
                if data == '':
                    # re_error로 대체 예정
                    return redirect('/bbs/w/' + bbs_num_str + '/' + post_num_str)
                
                date = get_time()

                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment', ?, ?, ?)"), [id_data, set_id, data])
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_date', ?, ?, ?)"), [id_data, set_id, date])
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_user_id', ?, ?, ?)"), [id_data, set_id, ip])

                conn.commit()

                return redirect('/bbs/w/' + bbs_num_str + '/' + post_num_str + '#comment_' + str(int(id_data) + 1))
            else:
                if acl_check(bbs_num_str, 'bbs_view') == 1:
                    return re_error('/ban')

                if do_type == 'preview':
                    text = flask.request.form.get('content', '')
                    text = text.replace('\r', '')

                    data_preview = render_set(
                        doc_name = '', 
                        doc_data = text,
                        data_in = 'from'
                    )
                else:
                    text = ''
                    data_preview = ''
                
                temp_id = ''
                temp_dict = {}

                db_data = list(db_data) if db_data else []
                for for_a in db_data + [['', '', '']]:
                    if temp_id != for_a[2]:
                        temp_id = for_a[2]
                        temp_dict['code'] = for_a[2]

                    temp_dict[for_a[0]] = for_a[1]

                count = 1

                data = ''
                data += '<h2>' + html.escape(temp_dict['title']) + '</h2>'
                data += bbs_w_post_make_thread(
                    ip_pas(temp_dict['user_id']),
                    temp_dict['date'],
                    render_set(
                        doc_name = '', 
                        doc_data = temp_dict['data'],
                        data_in = 'from'
                    ),
                    '1',
                    color = 'green'
                )
                data += '<hr class="main_hr">'

                user_id = temp_dict['user_id']

                temp_id = ''
                temp_dict = {}

                curs.execute(db_change('select set_name, set_data, set_code, set_id from bbs_data where (set_name = "comment" or set_name = "comment_date" or set_name = "comment_user_id") and set_id = ? order by set_code + 0 asc'), [bbs_num_str + '-' + post_num_str])
                db_data = curs.fetchall()
                db_data = list(db_data) if db_data else []

                for for_a in db_data + [['', '', '']]:
                    if temp_id == '':
                        temp_id = for_a[2]

                    if temp_id != for_a[2]:
                        temp_id = for_a[2]
                        temp_dict['code'] = for_a[2]
                        count += 1

                        if user_id == temp_dict['comment_user_id']:
                            color = 'green'
                        else:
                            color = 'default'

                        data += bbs_w_post_make_thread(
                            ip_pas(temp_dict['comment_user_id']),
                            temp_dict['comment_date'],
                            render_set(
                                doc_name = '', 
                                doc_data = temp_dict['comment'],
                                data_in = 'from'
                            ),
                            str(count),
                            color = color
                        )
                        data += '<hr class="main_hr">'

                    temp_dict[for_a[0]] = for_a[1]

                bbs_comment_form = ''
                if bbs_comment_acl == 0:
                    bbs_comment_form = '''                        
                        <textarea name="content" id="opennamu_edit_textarea" class="opennamu_textarea_100">''' + html.escape(text) + '''</textarea>
                        <hr class="main_hr">
                        
                        ''' + captcha_get() + ip_warning() + '''

                        <button id="opennamu_save_button" formaction="/bbs/w/''' + bbs_num_str + '''/''' + post_num_str + '''" type="submit">''' + load_lang('send') + '''</button>
                        <button id="opennamu_preview_button" formaction="/bbs/w/preview/''' + bbs_num_str + '''/''' + post_num_str + '''#opennamu_edit_textarea" type="submit">''' + load_lang('preview') + '''</button>
                        <hr class="main_hr">
                    '''

                data += '''
                    <form method="post">
                        ''' + bbs_comment_form + '''
                        ''' + data_preview + '''
                    </form>
                '''

                return easy_minify(flask.render_template(skin_check(),
                    imp = [load_lang('bbs_main'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                    data = data,
                    menu = [['bbs/w/' + bbs_num_str, load_lang('return')], ['bbs/edit/' + bbs_num_str + '/' + post_num_str, load_lang('edit')]]
                ))
        else:
            # db_data_2[0][0] == 'comment'
            if flask.request.method == 'POST' and do_type != 'preview':
                if bbs_comment_acl == 1:
                    return redirect('/bbs/set/' + bbs_num_str)
                
                if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                    return re_error('/error/13')
                else:
                    captcha_post('', 0)
                
                select = flask.request.form.get('comment_select', 'default')
                select = '' if select == 'default' else select
                if select != '':
                    select = select.split('-')
                    if len(select) < 2:
                        curs.execute(db_change('select set_code from bbs_data where set_name = "comment" and set_id = ? and set_code = ? limit 1'), [bbs_num_str + '-' + post_num_str, select[0]])
                        if not curs.fetchall():
                            return ''
                        else:
                            set_id = bbs_num_str + '-' + post_num_str + '-' + select[0]
                    else:
                        curs.execute(db_change('select set_code from bbs_data where set_name = "comment" and set_id = ? and set_code = ? limit 1'), [bbs_num_str + '-' + post_num_str + '-' + '-'.join(select[0:len(select) - 1]), select[len(select) - 1]])
                        if not curs.fetchall():
                            return ''
                        else:
                            set_id = bbs_num_str + '-' + post_num_str + '-' + '-'.join(select)
                else:
                    set_id = bbs_num_str + '-' + post_num_str

                curs.execute(db_change('select set_code from bbs_data where set_name = "comment" and set_id = ? order by set_code + 0 desc'), [set_id])
                db_data = curs.fetchall()
                id_data = str(int(db_data[0][0]) + 1) if db_data else '1'

                data = flask.request.form.get('content', '')
                if data == '':
                    # re_error로 대체 예정
                    return redirect('/bbs/w/' + bbs_num_str + '/' + post_num_str)

                date = get_time()

                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment', ?, ?, ?)"), [id_data, set_id, data])
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_date', ?, ?, ?)"), [id_data, set_id, date])
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_user_id', ?, ?, ?)"), [id_data, set_id, ip])

                conn.commit()
            
                if set_id == '':
                    return redirect('/bbs/w/' + bbs_num_str + '/' + post_num_str + '#comment_' + id_data)
                else:
                    set_id = re.sub(r'^[0-9]+-[0-9]+-?', '', set_id)
                    return redirect('/bbs/w/' + bbs_num_str + '/' + post_num_str + '#comment_' + set_id + '-' + id_data)
            else:
                if acl_check(bbs_num_str, 'bbs_view') == 1:
                    return re_error('/ban')
                    
                if do_type == 'preview':
                    text = flask.request.form.get('content', '')
                    text = text.replace('\r', '')

                    data_preview = render_set(
                        doc_name = '', 
                        doc_data = text,
                        data_in = 'from'
                    )
                else:
                    text = ''
                    data_preview = ''

                temp_id = ''
                temp_dict = {}

                db_data = list(db_data) if db_data else []
                for for_a in db_data + [['', '', '']]:
                    if temp_id != for_a[2]:
                        temp_id = for_a[2]
                        temp_dict['code'] = for_a[2]

                    temp_dict[for_a[0]] = for_a[1]

                data = ''
                data += '<h2>' + html.escape(temp_dict['title']) + '</h2>'
                data += bbs_w_post_make_thread(
                    ip_pas(temp_dict['user_id']),
                    temp_dict['date'],
                    render_set(
                        doc_name = '', 
                        doc_data = temp_dict['data'],
                        data_in = 'from'
                    ),
                    '0',
                    color = 'red'
                )

                user_id = temp_dict['user_id']
                temp_id = ''
                temp_dict = {}
                comment_data = ''

                comment_select = '<hr class="main_hr"><select name="comment_select">'
                comment_select += '<option value="default">' + load_lang('normal') + '</option>'

                curs.execute(db_change('select set_name, set_data, set_code, set_id from bbs_data where (set_name = "comment" or set_name = "comment_date" or set_name = "comment_user_id") and set_id = ? order by set_code + 0 asc'), [bbs_num_str + '-' + post_num_str])
                db_data = curs.fetchall()
                if db_data:
                    data += '<hr class="main_hr"><hr>'
                else:
                    db_data = []

                for_a = 0
                db_data_2 = db_data + [['', '', '', '']]
                db_data_len = len(db_data_2)
                comment_count = 0
                comment_add_count = 0
                while(for_a < db_data_len):
                    if temp_id != (db_data_2[for_a][3] + '-' + db_data_2[for_a][2]):
                        if temp_id != '':
                            temp_dict['code'] = temp_id
                            temp_dict['code'] = re.sub(r'^[0-9]+-[0-9]+-', '', temp_dict['code'])

                            if user_id == temp_dict['comment_user_id']:
                                color = 'green'
                            else:
                                color = 'default'

                            margin_count = temp_dict['code'].count('-')
                            if margin_count == 0:
                                comment_count += 1
                            else:
                                comment_add_count += 1

                            comment_data += '<span style="padding-left: 20px;"></span>' * margin_count
                            comment_data += bbs_w_post_make_thread(
                                ip_pas(temp_dict['comment_user_id']),
                                temp_dict['comment_date'],
                                render_set(
                                    doc_name = '', 
                                    doc_data = temp_dict['comment'],
                                    data_in = 'from'
                                ),
                                temp_dict['code'],
                                color = color,
                                add_style = 'width: calc(100% - ' + str(margin_count * 20) + 'px);'
                            )

                            comment_select += '<option value="' + temp_dict['code'] + '">' + temp_dict['code'] + '</option>'

                            curs.execute(db_change('select set_name, set_data, set_code, set_id from bbs_data where (set_name = "comment" or set_name = "comment_date" or set_name = "comment_user_id") and set_id = ? order by set_code + 0 asc'), [bbs_num_str + '-' + post_num_str + '-' + temp_dict['code']])
                            db_data = curs.fetchall()
                            if db_data:
                                db_data_2 = db_data_2[:for_a] + db_data + db_data_2[for_a:]
                                db_data_len += len(db_data)

                            if db_data_2[for_a][0] != '':
                                comment_data += '<hr class="main_hr">'

                        temp_id = db_data_2[for_a][3] + '-' + db_data_2[for_a][2]

                    temp_dict[db_data_2[for_a][0]] = db_data_2[for_a][1]
                    for_a += 1

                comment_select += '</select>'
                if comment_data != '':
                    data += load_lang('comment') + ' : ' + str(comment_count) + '<hr class="main_hr">'
                    data += load_lang('reply') + ' : ' + str(comment_add_count) + '<hr class="main_hr">'
                    data += comment_data

                bbs_comment_form = ''
                if bbs_comment_acl == 0:
                    bbs_comment_form = '''
                        ''' + comment_select + '''
                        <hr class="main_hr">
                        
                        <textarea name="content" id="opennamu_edit_textarea" class="opennamu_textarea_100">''' + html.escape(text) + '''</textarea>
                        <hr class="main_hr">
                        
                        ''' + captcha_get() + ip_warning() + '''

                        <button id="opennamu_save_button" formaction="/bbs/w/''' + bbs_num_str + '''/''' + post_num_str + '''" type="submit">''' + load_lang('send') + '''</button>
                        <button id="opennamu_preview_button" formaction="/bbs/w/preview/''' + bbs_num_str + '''/''' + post_num_str + '''#opennamu_edit_textarea" type="submit">''' + load_lang('preview') + '''</button>
                        <hr class="main_hr">
                    '''

                data += '''
                    <form method="post">
                        ''' + bbs_comment_form + '''
                        ''' + data_preview + '''
                    </form>
                '''

                return easy_minify(flask.render_template(skin_check(),
                    imp = [load_lang('bbs_main'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                    data = data,
                    menu = [['bbs/w/' + bbs_num_str, load_lang('return')], ['bbs/edit/' + bbs_num_str + '/' + post_num_str, load_lang('edit')]]
                ))