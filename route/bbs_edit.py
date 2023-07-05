from .tool.func import *

from .api_bbs_w_post import api_bbs_w_post

def bbs_edit(bbs_num = '', post_num = '', do_type = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        bbs_num_str = str(bbs_num)
        post_num_str = str(post_num)

        ip = ip_check()

        curs.execute(db_change('select set_id from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num_str])
        if not curs.fetchall():
            return redirect('/bbs/main')
        
        if post_num != '':
            curs.execute(db_change('select set_data from bbs_data where set_name = "user_id" and set_id = ? and set_code = ?'), [bbs_num, post_num])
            db_data = curs.fetchall()
            if not db_data:
                return redirect('/bbs/main')
            else:
                if not db_data[0][0] == ip and admin_check() != 1:
                    return re_error('/ban')
            
        if acl_check(bbs_num_str, 'bbs_edit') == 1:
            return redirect('/bbs/set/' + bbs_num_str)
        
        i_list = ['post_view_acl', 'post_comment_acl']

        if flask.request.method == 'POST' and do_type != 'preview':
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

            if post_num == '':
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('title', ?, ?, ?)"), [id_data, bbs_num_str, title])
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('data', ?, ?, ?)"), [id_data, bbs_num_str, data])
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('date', ?, ?, ?)"), [id_data, bbs_num_str, date])
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('user_id', ?, ?, ?)"), [id_data, bbs_num_str, ip])
            else:
                curs.execute(db_change("update bbs_data set set_data = ? where set_name = 'title' and set_code = ? and set_id = ?"), [title, post_num, bbs_num_str])
                curs.execute(db_change("update bbs_data set set_data = ? where set_name = 'data' and set_code = ? and set_id = ?"), [data, id_data, bbs_num_str])
                curs.execute(db_change("update bbs_data set set_data = ? where set_name = 'date' and set_code = ? and set_id = ?"), [date, id_data, bbs_num_str])

            return redirect('/bbs/w/' + bbs_num_str + '/' + id_data)
        else:
            d_list = ['' for _ in range(0, len(i_list))]
            if do_type == 'preview':
                title = flask.request.form.get('title', '')
                data = flask.request.form.get('content', '')
                data = data.replace('\r', '')
                data_preview = render_set(
                    doc_data = data,
                    data_type = 'thread',
                    data_in = 'bbs'
                ) + '<hr>'
                for for_a in range(0, len(i_list)):
                    d_list[for_a] = flask.request.form.get(i_list[for_a], 'normal')
            else:
                if post_num == '':
                    title = ''
                    data = ''
                    data_preview = ''
                else:
                    curs.execute(db_change('select set_name, set_data, set_code from bbs_data where set_id = ? and set_code = ?'), [bbs_num, post_num])
                    db_data = curs.fetchall()
                    db_data = list(db_data) if db_data else []

                    temp_dict = json.loads(api_bbs_w_post(bbs_num_str + '-' + post_num_str).data)

                    title = temp_dict['title']
                    data = temp_dict['data']
                    data_preview = ''

            acl_div = ['' for _ in range(0, len(i_list))]
            acl_list = get_acl_list()
            for for_a in range(0, len(i_list)):
                for data_list in acl_list:
                    if data_list == d_list[for_a]:
                        check = 'selected="selected"'
                    else:
                        check = ''

                    acl_div[for_a] += '<option value="' + data_list + '" ' + check + '>' + (data_list if data_list != '' else 'normal') + '</option>'
            
            if post_num == '':
                form_action = 'formaction="/bbs/edit/' + bbs_num_str + '"'
                form_action_preview = 'formaction="/bbs/edit/preview/' + bbs_num_str + '"'
            else:
                form_action = 'formaction="/bbs/edit/' + bbs_num_str + '/' + post_num_str + '"'
                form_action_preview = 'formaction="/bbs/edit/preview/' + bbs_num_str + '/' + post_num_str + '"'
    
            editor_top_text = '<a href="/edit_filter">(' + load_lang('edit_filter_rule') + ')</a>'
            
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

            if post_num == '':
                bbs_title = load_lang('post_add')
            else:
                bbs_title = load_lang('post_edit')
    
            return easy_minify(flask.render_template(skin_check(), 
                imp = [bbs_title, wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data =  editor_top_text + add_get_file + '''
                    <form method="post">
                        <textarea style="display: none;" id="opennamu_edit_origin" name="doc_data_org"></textarea>

                        <div>''' + edit_button('opennamu_edit_textarea', 'opennamu_monaco_editor') + '''</div>
                        
                        <input placeholder="''' + load_lang('title') + '''" name="title" value="''' + html.escape(title) + '''">
                        <hr class="main_hr">

                        <div id="opennamu_monaco_editor" class="opennamu_textarea_500" ''' + monaco_display + '''></div>
                        <textarea id="opennamu_edit_textarea" ''' + editor_display + ''' class="opennamu_textarea_500" name="content">''' + html.escape(data) + '''</textarea>
                        <hr class="main_hr">
                        
                        ''' + captcha_get() + ip_warning() + '''
                        
                        <button id="opennamu_save_button" type="submit" ''' + form_action + ''' onclick="do_monaco_to_textarea(); do_stop_exit_release();">''' + load_lang('save') + '''</button>
                        <button id="opennamu_preview_button" type="submit" ''' + form_action_preview + ''' onclick="do_monaco_to_textarea(); do_stop_exit_release();">''' + load_lang('preview') + '''</button>
                    
                        <hr class="main_hr">
                        <div id="opennamu_preview_area">''' + data_preview + '''</div>
                        
                        <script>
                            do_stop_exit();
                            do_paste_image('opennamu_edit_textarea', 'opennamu_monaco_editor');
                            ''' + add_script + '''
                        </script>

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
                    </form>
                ''',
                menu = [['bbs/w/' + bbs_num_str, load_lang('return')]]
            ))