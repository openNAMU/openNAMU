from .tool.func import *

from .api_bbs_w_post import api_bbs_w_post
from .api_bbs_w_comment_one import api_bbs_w_comment_one

from .edit import edit_editor

def bbs_w_edit(bbs_num = '', post_num = '', comment_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        bbs_num_str = str(bbs_num)
        post_num_str = str(post_num)

        ip = ip_check()

        curs.execute(db_change('select set_id from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num_str])
        if not curs.fetchall():
            return redirect('/bbs/main')
        
        if comment_num != '':
            temp_dict = json.loads(api_bbs_w_comment_one(bbs_num_str + '-' + post_num_str + '-' + comment_num).data)
            if 'comment_user_id' in temp_dict:
                if not temp_dict['comment_user_id'] == ip and admin_check() != 1:
                    return re_error('/ban')
            else:
                return redirect('/bbs/main')
        elif post_num != '':
            temp_dict = json.loads(api_bbs_w_post(bbs_num_str + '-' + post_num_str).data)
            if 'user_id' in temp_dict:
                if not temp_dict['user_id'] == ip and admin_check() != 1:
                    return re_error('/ban')
            else:
                return redirect('/bbs/main')
            
        if acl_check(bbs_num_str, 'bbs_edit') == 1:
            return redirect('/bbs/set/' + bbs_num_str)
        
        i_list = ['post_view_acl', 'post_comment_acl']

        if flask.request.method == 'POST':
            if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return re_error('/error/13')
            else:
                captcha_post('', 0)
        
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
                return redirect('/bbs/w/' + bbs_num_str)
            
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
                return redirect('/bbs/w/' + bbs_num_str + '/' + id_data + '#' + url_pas(comment_num))
            else:
                return redirect('/bbs/w/' + bbs_num_str + '/' + id_data)
        else:
            option_display = ''

            if comment_num != '':
                temp_dict = json.loads(api_bbs_w_comment_one(bbs_num_str + '-' + post_num_str + '-' + comment_num).data)

                title = ''
                data = temp_dict['comment']
                option_display = 'display: none;'
            elif post_num == '':
                title = ''
                data = ''
            else:
                temp_dict = json.loads(api_bbs_w_post(bbs_num_str + '-' + post_num_str).data)

                title = temp_dict['title']
                data = temp_dict['data']

            acl_div = ['' for _ in range(0, len(i_list))]
            acl_list = get_acl_list()
            for for_a in range(0, len(i_list)):
                for data_list in acl_list:
                    acl_div[for_a] += '<option value="' + data_list + '">' + (data_list if data_list != '' else 'normal') + '</option>'
    
            editor_top_text = '<a href="/filter/edit_filter">(' + load_lang('edit_filter_rule') + ')</a>'

            if editor_top_text != '':
                editor_top_text += '<hr class="main_hr">'

            if comment_num != '':
                bbs_title = load_lang('bbs_comment_edit')
            elif post_num == '':
                bbs_title = load_lang('post_add')
            else:
                bbs_title = load_lang('post_edit')
    
            return easy_minify(flask.render_template(skin_check(), 
                imp = [bbs_title, wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data =  editor_top_text + '''
                    <form method="post">                        
                        <input style="''' + option_display + '''" placeholder="''' + load_lang('title') + '''" name="title" value="''' + html.escape(title) + '''">
                        <hr style="''' + option_display + '''" class="main_hr">

                        ''' + edit_editor(curs, ip, data, 'bbs') + '''

                        <!--
                        <div style="''' + option_display + '''">
                            ''' + render_simple_set('''
                                <hr class="main_hr">
                                <a href="/acl/TEST#exp">(''' + load_lang('reference') + ''')</a>
                                <h2>''' + load_lang('acl') + '''</h2>
                                <h3>''' + load_lang('post_view_acl') + '''</h3>
                                <select name="post_view_acl">''' + acl_div[0] + '''</select>

                                <h4>''' + load_lang('post_comment_acl') + '''</h4>
                                <select name="post_comment_acl">''' + acl_div[1] + '''</select>

                                <h2>''' + load_lang('markup') + '''</h2>
                                ''' + load_lang('not_working') + '''
                            ''') + '''
                        </div>
                        -->
                    </form>
                ''',
                menu = [['bbs/w/' + bbs_num_str, load_lang('return')]]
            ))