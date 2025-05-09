from .tool.func import *

from .go_api_bbs_w import api_bbs_w
from .go_api_bbs_w_comment_one import api_bbs_w_comment_one

from .edit import edit_editor

async def bbs_w_edit(bbs_num = '', post_num = '', comment_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        bbs_num_str = str(bbs_num)
        post_num_str = str(post_num)

        ip = ip_check()

        curs.execute(db_change('select set_id from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num_str])
        if not curs.fetchall():
            return redirect(conn, '/bbs/main')
        
        if comment_num != '':
            temp_dict = await api_bbs_w_comment_one(bbs_num_str + '-' + post_num_str + '-' + comment_num)
            if 'comment_user_id' in temp_dict:
                if not temp_dict['comment_user_id'] == ip and await acl_check('', 'owner_auth', '', '') == 1:
                    return await re_error(conn, 0)
            else:
                return redirect(conn, '/bbs/main')
        elif post_num != '':
            temp_dict = await api_bbs_w(bbs_num_str + '-' + post_num_str)
            if 'user_id' in temp_dict:
                if not temp_dict['user_id'] == ip and await acl_check('', 'owner_auth', '', '') == 1:
                    return await re_error(conn, 0)
            else:
                return redirect(conn, '/bbs/main')
            
        if await acl_check(bbs_num_str, 'bbs_edit') == 1:
            return redirect(conn, '/bbs/set/' + bbs_num_str)
        
        i_list = ['post_view_acl', 'post_comment_acl']

        if flask.request.method == 'POST':
            if await captcha_post(conn, flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return await re_error(conn, 13)
        
            if post_num == '':
                curs.execute(db_change('select set_code from bbs_data where set_name = "title" and set_id = ? order by set_code + 0 desc'), [bbs_num_str])
                db_data = curs.fetchall()
                id_data = str(int(db_data[0][0]) + 1) if db_data else '1'
            else:
                id_data = post_num_str

            title = flask.request.form.get('title', 'test')
            title = 'test' if title == '' else title
            data = flask.request.form.get('content', '')
            if data == '':
                # re_error로 대체 예정
                return redirect(conn, '/bbs/in/' + bbs_num_str)
            
            if await do_edit_filter(conn, title) == 1:
                return await re_error(conn, 21)

            if await do_edit_filter(conn, data) == 1:
                return await re_error(conn, 21)
            
            date = get_time()

            if comment_num != '':
                sub_code = (bbs_num_str + '-' + post_num_str + '-' + comment_num).split('-')
                sub_code_last = ''
                if len(sub_code) > 2:
                    sub_code_last = sub_code[len(sub_code) - 1]
                    del sub_code[len(sub_code) - 1]
                    
                sub_code = '-'.join(sub_code)

                curs.execute(db_change("update bbs_data set set_data = ? where set_name = 'comment' and set_code = ? and set_id = ?"), [data, sub_code_last, sub_code])
            elif post_num == '':
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('title', ?, ?, ?)"), [id_data, bbs_num_str, title])
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('data', ?, ?, ?)"), [id_data, bbs_num_str, data])
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('date', ?, ?, ?)"), [id_data, bbs_num_str, date])
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('user_id', ?, ?, ?)"), [id_data, bbs_num_str, ip])
            else:
                curs.execute(db_change("update bbs_data set set_data = ? where set_name = 'title' and set_code = ? and set_id = ?"), [title, post_num, bbs_num_str])
                curs.execute(db_change("update bbs_data set set_data = ? where set_name = 'data' and set_code = ? and set_id = ?"), [data, id_data, bbs_num_str])
                curs.execute(db_change("update bbs_data set set_data = ? where set_name = 'date' and set_code = ? and set_id = ?"), [date, id_data, bbs_num_str])

            if comment_num != '':
                return redirect(conn, '/bbs/w/' + bbs_num_str + '/' + id_data + '#' + url_pas(comment_num))
            else:
                return redirect(conn, '/bbs/w/' + bbs_num_str + '/' + id_data)
        else:
            option_display = ''

            if comment_num != '':
                temp_dict = await api_bbs_w_comment_one(bbs_num_str + '-' + post_num_str + '-' + comment_num)

                title = ''
                data = temp_dict['comment']
                option_display = 'display: none;'
            elif post_num == '':
                title = ''
                data = ''
            else:
                temp_dict = await api_bbs_w(bbs_num_str + '-' + post_num_str)

                title = temp_dict['title']
                data = temp_dict['data']

            acl_div = ['' for _ in range(0, len(i_list))]
            acl_list = await get_acl_list()
            for for_a in range(0, len(i_list)):
                for data_list in acl_list:
                    acl_div[for_a] += '<option value="' + data_list + '">' + (data_list if data_list != '' else 'normal') + '</option>'
    
            editor_top_text = '<a href="/filter/edit_filter">(' + get_lang(conn, 'edit_filter_rule') + ')</a>'

            if editor_top_text != '':
                editor_top_text += '<hr class="main_hr">'

            if comment_num != '':
                bbs_title = get_lang(conn, 'bbs_comment_edit')
            elif post_num == '':
                bbs_title = get_lang(conn, 'post_add')
            else:
                bbs_title = get_lang(conn, 'post_edit')
    
            return easy_minify(conn, flask.render_template(skin_check(conn), 
                imp = [bbs_title, await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data =  editor_top_text + '''
                    <form method="post">                        
                        <input style="''' + option_display + '''" placeholder="''' + get_lang(conn, 'title') + '''" name="title" value="''' + html.escape(title) + '''">
                        <hr style="''' + option_display + '''" class="main_hr">

                        ''' + await edit_editor(conn, ip, data, 'bbs') + '''

                        <!--
                        <div style="''' + option_display + '''">
                            ''' + render_simple_set(conn, '''
                                <hr class="main_hr">
                                <a href="/acl/TEST#exp">(''' + get_lang(conn, 'reference') + ''')</a>
                                <h2>''' + get_lang(conn, 'acl') + '''</h2>
                                <h3>''' + get_lang(conn, 'post_view_acl') + '''</h3>
                                <select name="post_view_acl">''' + acl_div[0] + '''</select>

                                <h4>''' + get_lang(conn, 'post_comment_acl') + '''</h4>
                                <select name="post_comment_acl">''' + acl_div[1] + '''</select>

                                <h2>''' + get_lang(conn, 'markup') + '''</h2>
                                ''' + get_lang(conn, 'not_working') + '''
                            ''') + '''
                        </div>
                        -->
                    </form>
                ''',
                menu = [['bbs/in/' + bbs_num_str, get_lang(conn, 'return')]]
            ))