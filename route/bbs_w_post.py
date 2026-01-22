from .tool.func import *

from .go_api_bbs_w import api_bbs_w
from .go_api_bbs_w_comment import api_bbs_w_comment

from .go_api_topic import api_topic_thread_make, api_topic_thread_pre_render

from .edit import edit_editor

async def bbs_w_post(bbs_num = '', post_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
        db_data_3 = curs.fetchall()
        if not db_data_3:
            return redirect(conn, '/bbs/main')
        
        bbs_name = db_data_3[0][0]

        bbs_num_str = str(bbs_num)
        post_num_str = str(post_num)
        bbs_comment_acl = await acl_check(bbs_num_str, 'bbs_comment')
        ip = ip_check()

        temp_dict = await api_bbs_w(bbs_num_str + '-' + post_num_str)
        if temp_dict == {}:
            return redirect(conn, '/bbs/main')
        
        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_type"'), [bbs_num])
        db_data_2 = curs.fetchall()
        if not db_data_2:
            return redirect(conn, '/bbs/main')
        else:
            if flask.request.method == 'POST':
                if db_data_2[0][0] == 'thread':
                    if bbs_comment_acl == 1:
                        return redirect(conn, '/bbs/set/' + bbs_num_str)
                    
                    if await captcha_post(conn, flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                        return await re_error(conn, 13)

                    set_id = bbs_num_str + '-' + post_num_str

                    curs.execute(db_change('select set_code from bbs_data where set_name = "comment" and set_id = ? order by set_code + 0 desc'), [set_id])
                    db_data_4 = curs.fetchall()
                    id_data = str(int(db_data_4[0][0]) + 1) if db_data_4 else '1'

                    data = flask.request.form.get('content', '')
                    if data == '':
                        # re_error로 대체 예정
                        return redirect(conn, '/bbs/w/' + bbs_num_str + '/' + post_num_str)
                    
                    data = data.replace('\r', '')
                    data = await api_topic_thread_pre_render(conn, data, id_data, ip, set_id, bbs_name, temp_dict['title'], 'post')
                    
                    date = get_time()

                    curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment', ?, ?, ?)"), [id_data, set_id, data])
                    curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_date', ?, ?, ?)"), [id_data, set_id, date])
                    curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_user_id', ?, ?, ?)"), [id_data, set_id, ip])

                    await add_alarm(temp_dict['user_id'], ip, 'BBS <a href="/bbs/w/' + bbs_num_str + '/' + post_num_str + '#' + id_data + '">' + html.escape(bbs_name) + ' - ' + html.escape(temp_dict['title']) + '#' + id_data + '</a>')

                    return redirect(conn, '/bbs/w/' + bbs_num_str + '/' + post_num_str + '#' + id_data)
                else:
                    if bbs_comment_acl == 1:
                        return redirect(conn, '/bbs/set/' + bbs_num_str)
                    
                    if await captcha_post(conn, flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                        return await re_error(conn, 13)
                    
                    select = flask.request.form.get('comment_select', '0')
                    select = '' if select == '0' else select

                    comment_user_name = ''

                    if select != '':
                        select_split = select.split('-')
                        if len(select_split) < 2:
                            curs.execute(db_change('select set_data from bbs_data where set_name = "comment_user_id" and set_id = ? and set_code = ? limit 1'), [bbs_num_str + '-' + post_num_str, select_split[0]])    
                            db_data_6 = curs.fetchall()
                            if not db_data_6:
                                # re_error로 변경 예정
                                return redirect(conn, '/bbs/w/' + bbs_num_str + '/' + post_num_str)
                            else:
                                set_id = bbs_num_str + '-' + post_num_str + '-' + select_split[0]
                                comment_user_name = db_data_6[0][0]
                        else:
                            curs.execute(db_change('select set_data from bbs_data where set_name = "comment_user_id" and set_id = ? and set_code = ? limit 1'), [bbs_num_str + '-' + post_num_str + '-' + '-'.join(select_split[0:len(select_split) - 1]), select_split[len(select_split) - 1]])
                            db_data_7 = curs.fetchall()
                            if not db_data_7:
                                return redirect(conn, '/bbs/w/' + bbs_num_str + '/' + post_num_str)
                            else:
                                set_id = bbs_num_str + '-' + post_num_str + '-' + '-'.join(select_split)
                                comment_user_name = db_data_7[0][0]
                    else:
                        set_id = bbs_num_str + '-' + post_num_str

                    curs.execute(db_change('select set_code from bbs_data where set_name = "comment" and set_id = ? order by set_code + 0 desc limit 1'), [set_id])
                    db_data_5 = curs.fetchall()
                    id_data = str(int(db_data_5[0][0]) + 1) if db_data_5 else '1'

                    data = flask.request.form.get('content', '')
                    if data == '':
                        # re_error로 대체 예정
                        return redirect(conn, '/bbs/w/' + bbs_num_str + '/' + post_num_str)

                    date = get_time()

                    curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment', ?, ?, ?)"), [id_data, set_id, data])
                    curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_date', ?, ?, ?)"), [id_data, set_id, date])
                    curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_user_id', ?, ?, ?)"), [id_data, set_id, ip])
                
                    if set_id == '':
                        end_id = id_data
                    else:
                        set_id = re.sub(r'^[0-9]+-[0-9]+-?', '', set_id)
                        set_id += '-' if set_id != '' else ''
                        end_id = set_id + id_data

                    await add_alarm(temp_dict['user_id'], ip, 'BBS <a href="/bbs/w/' + bbs_num_str + '/' + post_num_str + '#' + end_id + '">' + html.escape(bbs_name) + ' - ' + html.escape(temp_dict['title']) + '#' + end_id + '</a>')
                    if comment_user_name != '':
                        await add_alarm(comment_user_name, ip, 'BBS <a href="/bbs/w/' + bbs_num_str + '/' + post_num_str + '#' + end_id + '">' + html.escape(bbs_name) + ' - ' + html.escape(temp_dict['title']) + '#' + end_id + '</a>')

                    return redirect(conn, '/bbs/w/' + bbs_num_str + '/' + post_num_str + '#' + end_id)
            else:
                if await acl_check(bbs_num_str, 'bbs_view') == 1:
                    return await re_error(conn, 0)

                date = ''
                date += '<a href="javascript:opennamu_change_comment(\'0\');">(' + await get_lang('comment') + ')</a> '
                date += temp_dict['date']

                data = '<div class="opennamu_bbs_w_post_tab">'
                data += '<big><big><big>' + html.escape(temp_dict['title']) + '</big></big></big>'
                data += '<hr class="main_hr">'
                data += await ip_pas(temp_dict['user_id']) + '<span style="float: right;">' + date + '</span>'
                data += '<hr>'
                data += '<div class="opennamu_bbs_w_post_tab_content">' + await render_set(conn, doc_data = temp_dict['data']) + '</div>'
                data += '</div>'

                if bbs_comment_acl == 0:
                    data += '<hr class="main_hr">'
                    data += '<div id="opennamu_bbs_w_post_tabom"></div>'

                data += '' + \
                    '<hr>' + \
                    '<div id="opennamu_bbs_w_post"></div>' + \
                    '<script defer src="/views/main_css/js/route/topic.js' + cache_v() + '"></script>' + \
                    '<script defer src="/views/main_css/js/route/bbs_w_post.js' + cache_v() + '"></script>' + \
                    '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_load_comment(); });</script>' + \
                ''

                bbs_comment_form = ''
                if bbs_comment_acl == 0:
                    bbs_comment_form += '''
                        <div id="opennamu_bbs_w_post_select"></div>
                        ''' + await edit_editor(conn, ip, '', 'bbs_comment') + '''
                    '''

                data += '''
                    <form method="post">
                        ''' + bbs_comment_form + '''
                    </form>
                '''
                
                await python_to_golang("get_json", path = "v2/bbs/w/page_view_post/" + url_pas(bbs_num_str) + "/" + url_pas(post_num_str))

                view_count_data = await python_to_golang("get_json", path = "v2/bbs/w/page_view/" + url_pas(bbs_num_str) + "/" + url_pas(post_num_str))
                view_count = view_count_data['data']

                return await render_template(
                    bbs_name,
                    data,
                    '(' + await get_lang('bbs') + ')',
                    [['bbs/in/' + bbs_num_str, await get_lang('return')], ['bbs/edit/' + bbs_num_str + '/' + post_num_str, await get_lang('edit')], ['bbs/tool/' + bbs_num_str + '/' + post_num_str, await get_lang('tool')]],
                    [temp_dict['date'], 0, 0, view_count],
                )
