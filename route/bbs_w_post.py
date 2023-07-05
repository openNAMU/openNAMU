from .tool.func import *

from .api_bbs_w_post import api_bbs_w_post
from .api_bbs_w_comment import api_bbs_w_comment

def bbs_w_post_make_thread(user_id : str, date : str, data : str, code : str, color : str = '', blind : str = '', add_style : str = '') -> str:
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

def bbs_w_post_comment(user_id : str, sub_code : str, comment_num : str, bbs_num_str : str, post_num_str : str) -> typing.Tuple[str, str, int, int]:
    comment_data : str = ''
    comment_select : str = ''

    comment_count : int = 0
    comment_add_count : int = 0

    thread_data : typing.List[dict[str, str]] = json.loads(api_bbs_w_comment(sub_code).data)
    
    comment_count += len(thread_data)
    comment_add_count += comment_count

    temp_dict : dict[str, str]
    for temp_dict in thread_data:
        color : str = 'default'
        if user_id == temp_dict['comment_user_id']:
            color = 'green'

        sub_code_check : str = re.sub(r'^[0-9]+-[0-9]+-', '', sub_code + '-' + temp_dict['code'])
        margin_count : int = sub_code_check.count('-')

        date : str = ''
        date += '<a href="javascript:opennamu_change_comment(\'' + sub_code_check + '\');">(' + load_lang('comment') + ')</a> '
        date += '<a href="/bbs/w/' + bbs_num_str + '/' + post_num_str + '/comment/' + sub_code_check + '/tool">(' + load_lang('tool') + ')</a> '
        date += temp_dict['comment_date']

        comment_data += '<span style="padding-left: 20px;"></span>' * margin_count
        comment_data += bbs_w_post_make_thread(
            ip_pas(temp_dict['comment_user_id']),
            date,
            render_set(
                doc_data = temp_dict['comment'],
                data_in = 'bbs_comment_' + sub_code_check
            ),
            sub_code_check,
            color = color,
            add_style = 'width: calc(100% - ' + str(margin_count * 20) + 'px);'
        )

        comment_default : str = ''
        if comment_num == sub_code_check:
            comment_default = 'selected'

        comment_select += '<option value="' + sub_code_check + '" ' + comment_default + '>' + sub_code_check + '</option>'
        comment_data += '<hr class="main_hr">'

        temp_data : typing.Tuple[str, str, int, int] = bbs_w_post_comment(user_id, sub_code + '-' + temp_dict['code'], comment_num, bbs_num_str, post_num_str)

        comment_data += temp_data[0]
        comment_select += temp_data[1]
        comment_add_count += temp_data[3]

    return (comment_data, comment_select, comment_count, comment_add_count)

def bbs_w_post(bbs_num : typing.Union[int, str] = '', post_num : typing.Union[int, str] = '', do_type : str = '') -> typing.Union[str, werkzeug.wrappers.response.Response]:
    conn : typing.Union[sqlite3.Connection, pymysql.connections.Connection]
    with get_db_connect() as conn:
        curs : typing.Union[sqlite3.Cursor, pymysql.cursors.Cursor] = conn.cursor()

        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
        db_data_3 : typing.Optional[typing.List[typing.Tuple[str]]] = curs.fetchall()
        if not db_data_3:
            return redirect('/bbs/main')
        
        bbs_name : str = db_data_3[0][0]

        bbs_num_str : str = str(bbs_num)
        post_num_str : str = str(post_num)
        bbs_comment_acl : int = acl_check(bbs_num_str, 'bbs_comment')
        ip : str = ip_check()

        set_id : str
        text : str
        data_preview : str
        user_id : str
        bbs_comment_form : str
        id_data : str
        data : str
        date : str
        temp_dict : dict[str, str]
        new_id_data : str

        temp_dict = json.loads(api_bbs_w_post(bbs_num_str + '-' + post_num_str).data)
        if temp_dict == {}:
            return redirect('/bbs/main')
        
        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_type"'), [bbs_num])
        db_data_2 : typing.Optional[typing.List[typing.Tuple[str]]] = curs.fetchall()
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
                db_data_4 : typing.Optional[typing.List[typing.Tuple[str]]] = curs.fetchall()
                id_data = str(int(db_data_4[0][0]) + 1) if db_data_4 else '1'

                data = flask.request.form.get('content', '')
                if data == '':
                    # re_error로 대체 예정
                    return redirect('/bbs/w/' + bbs_num_str + '/' + post_num_str)
                
                data = data.replace('\r', '')
                data = get_thread_pre_render(data, id_data, ip, set_id, bbs_name, temp_dict['title'], 'post')
                
                date = get_time()

                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment', ?, ?, ?)"), [id_data, set_id, data])
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_date', ?, ?, ?)"), [id_data, set_id, date])
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_user_id', ?, ?, ?)"), [id_data, set_id, ip])

                add_alarm(temp_dict['user_id'], ip, 'BBS <a href="/bbs/w/' + bbs_num_str + '/' + post_num_str + '#' + id_data + '">' + html.escape(bbs_name) + ' - ' + html.escape(temp_dict['title']) + '#' + id_data + '</a>')

                conn.commit()

                return redirect('/bbs/w/' + bbs_num_str + '/' + post_num_str + '#' + id_data)
            else:
                if acl_check(bbs_num_str, 'bbs_view') == 1:
                    return re_error('/ban')

                text = ''
                data_preview = ''
                if do_type == 'preview':
                    text = flask.request.form.get('content', '')
                    text = text.replace('\r', '')

                    data_preview = render_set(
                        doc_data = text,
                        data_type = 'thread',
                        data_in = 'bbs_comment_preview'
                    )

                date = ''
                date += '<a href="/bbs/w/' + bbs_num_str + '/' + post_num_str + '/tool">(' + load_lang('tool') + ')</a> '
                date += temp_dict['date']

                data = ''
                data += '<h2>' + html.escape(temp_dict['title']) + '</h2>'
                data += bbs_w_post_make_thread(
                    ip_pas(temp_dict['user_id']),
                    date,
                    render_set(
                        doc_data = temp_dict['data'],
                        data_type = 'thread',
                        data_in = 'bbs'
                    ),
                    '1',
                    color = 'green'
                )
                data += '<hr class="main_hr">'

                user_id = temp_dict['user_id']
                count : int = 1

                thread_data : typing.List[dict[str, str]] = json.loads(api_bbs_w_comment(bbs_num_str + '-' + post_num_str).data)
                for temp_dict in thread_data:
                    count += 1
                    if user_id == temp_dict['comment_user_id']:
                        color = 'green'
                    else:
                        color = 'default'
                        
                    date = ''
                    date += '<a href="/bbs/w/' + bbs_num_str + '/' + post_num_str + '/comment/' + str(count) + '/tool">(' + load_lang('tool') + ')</a> '
                    date += temp_dict['comment_date']

                    data += bbs_w_post_make_thread(
                        ip_pas(temp_dict['comment_user_id']),
                        date,
                        render_set(
                            doc_data = temp_dict['comment'],
                            data_type = 'thread',
                            data_in = 'bbs_comment_' + str(count)
                        ),
                        str(count),
                        color = color
                    )
                    data += '<hr class="main_hr">'

                bbs_comment_form = ''
                if bbs_comment_acl == 0:
                    bbs_comment_form += '''                        
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
                    imp = [bbs_name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('bbs') + ')', 0])],
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
                
                select : str = flask.request.form.get('comment_select', '0')
                select = '' if select == '0' else select

                comment_user_name : str = ''

                if select != '':
                    select_split : typing.List[str] = select.split('-')
                    if len(select_split) < 2:
                        curs.execute(db_change('select set_data from bbs_data where set_name = "comment_user_id" and set_id = ? and set_code = ? limit 1'), [bbs_num_str + '-' + post_num_str, select_split[0]])    
                        db_data_6 : typing.Optional[typing.List[typing.Tuple[str]]] = curs.fetchall()
                        if not db_data_6:
                            # re_error로 변경 예정
                            return redirect('/bbs/w/' + bbs_num_str + '/' + post_num_str)
                        else:
                            set_id = bbs_num_str + '-' + post_num_str + '-' + select_split[0]
                            comment_user_name = db_data_6[0][0]
                    else:
                        curs.execute(db_change('select set_data from bbs_data where set_name = "comment_user_id" and set_id = ? and set_code = ? limit 1'), [bbs_num_str + '-' + post_num_str + '-' + '-'.join(select_split[0:len(select_split) - 1]), select_split[len(select_split) - 1]])
                        db_data_7 : typing.Optional[typing.List[typing.Tuple[str]]] = curs.fetchall()
                        if not db_data_7:
                            return redirect('/bbs/w/' + bbs_num_str + '/' + post_num_str)
                        else:
                            set_id = bbs_num_str + '-' + post_num_str + '-' + '-'.join(select_split)
                            comment_user_name = db_data_7[0][0]
                else:
                    set_id = bbs_num_str + '-' + post_num_str

                curs.execute(db_change('select set_code from bbs_data where set_name = "comment" and set_id = ? order by set_code + 0 desc limit 1'), [set_id])
                db_data_5 : typing.Optional[typing.List[typing.Tuple[str]]] = curs.fetchall()
                id_data = str(int(db_data_5[0][0]) + 1) if db_data_5 else '1'

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
                    end_id = id_data
                else:
                    set_id = re.sub(r'^[0-9]+-[0-9]+-?', '', set_id)
                    set_id += '-' if set_id != '' else ''
                    end_id = set_id + id_data

                add_alarm(temp_dict['user_id'], ip, 'BBS <a href="/bbs/w/' + bbs_num_str + '/' + post_num_str + '#' + end_id + '">' + html.escape(bbs_name) + ' - ' + html.escape(temp_dict['title']) + '#' + end_id + '</a>')
                if comment_user_name != '':
                    add_alarm(comment_user_name, ip, 'BBS <a href="/bbs/w/' + bbs_num_str + '/' + post_num_str + '#' + end_id + '">' + html.escape(bbs_name) + ' - ' + html.escape(temp_dict['title']) + '#' + end_id + '</a>')

                return redirect('/bbs/w/' + bbs_num_str + '/' + post_num_str + '#' + end_id)
            else:
                if acl_check(bbs_num_str, 'bbs_view') == 1:
                    return re_error('/ban')
                    
                text = ''
                comment_num : str = ''
                data_preview = ''
                if do_type == 'preview':
                    text = flask.request.form.get('content', '')
                    text = text.replace('\r', '')

                    comment_num = flask.request.form.get('comment_select', '')

                    data_preview = render_set(
                        doc_data = text,
                        data_in = 'bbs_comment_preview'
                    )

                date = ''
                date += '<a href="javascript:opennamu_change_comment(\'0\');">(' + load_lang('comment') + ')</a> '
                date += '<a href="/bbs/w/' + bbs_num_str + '/' + post_num_str + '/tool">(' + load_lang('tool') + ')</a> '
                date += temp_dict['date']

                data = ''
                data += '<h2>' + html.escape(temp_dict['title']) + '</h2>'
                data += bbs_w_post_make_thread(
                    ip_pas(temp_dict['user_id']),
                    date,
                    render_set(
                        doc_data = temp_dict['data'],
                        data_in = 'bbs'
                    ),
                    '0',
                    color = 'red'
                )

                user_id = temp_dict['user_id']
                comment_data : str = ''

                comment_select : str = '<select id="opennamu_comment_select" name="comment_select">'
                comment_select += '<option value="0">' + load_lang('normal') + '</option>'

                comment_count : int = 0
                comment_add_count : int = 0

                temp_data : typing.Tuple[str, str, int, int] = bbs_w_post_comment(user_id, bbs_num_str + '-' + post_num_str, comment_num, bbs_num_str, post_num_str)

                comment_data += temp_data[0]
                comment_select += temp_data[1]
                comment_count += temp_data[2]
                comment_add_count += temp_data[3]
                comment_add_count -= comment_count

                if comment_data != '':
                    data += '<hr class="main_hr"><hr>'

                comment_select += '</select>'
                if comment_data != '':
                    data += load_lang('comment') + ' : ' + str(comment_count) + '<hr class="main_hr">'
                    data += load_lang('reply') + ' : ' + str(comment_add_count) + '<hr class="main_hr">'
                    data += comment_data
                else:
                    data += '<hr class="main_hr">'

                bbs_comment_form = ''
                if bbs_comment_acl == 0:
                    bbs_comment_form += '''
                        ''' + comment_select + ''' <a href="javascript:opennamu_return_comment();">(R)</a>
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
                    <script src="/views/main_css/js/route/bbs_w_post.js"></script>
                '''

                return easy_minify(flask.render_template(skin_check(),
                    imp = [bbs_name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('bbs') + ')', 0])],
                    data = data,
                    menu = [['bbs/w/' + bbs_num_str, load_lang('return')], ['bbs/edit/' + bbs_num_str + '/' + post_num_str, load_lang('edit')]]
                ))