from .tool.func import *

def main_setting_main(db_set):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/ban')
        
        setting_list = {
            0 : ['name', 'Wiki'],
            2 : ['frontpage', 'FrontPage'],
            4 : ['upload', '2'],
            5 : ['skin', ''],
            7 : ['reg', ''],
            8 : ['ip_view', ''],
            9 : ['back_up', ''],
            10 : ['port', '3000'],
            11 : ['key', load_random_key()],
            12 : ['update', 'stable'],
            15 : ['encode', 'sha3'],
            16 : ['host', '0.0.0.0'],
            19 : ['slow_edit', ''],
            20 : ['requires_approval', ''],
            21 : ['backup_where', ''],
            22 : ['domain', ''],
            23 : ['ua_get', ''],
            24 : ['enable_comment', ''],
            25 : ['enable_challenge', ''],
            26 : ['edit_bottom_compulsion', ''],
            27 : ['http_select', 'http'],
            28 : ['title_max_length', ''],
            29 : ['title_topic_max_length', ''],
            30 : ['password_min_length', ''],
            31 : ['wiki_access_password_need', ''],
            32 : ['wiki_access_password', ''],
            33 : ['history_recording_off', ''],
            34 : ['namumark_compatible', ''],
            35 : ['user_name_view', ''],
            36 : ['link_case_insensitive', ''],
            37 : ['move_with_redirect', ''],
            38 : ['slow_thread', ''],
            39 : ['edit_timeout', '5'],
            40 : ['document_content_max_length', ''],
            41 : ['backup_count', ''],
            42 : ['ua_expiration_date', ''],
            43 : ['auth_history_expiration_date', ''],
            44 : ['auth_history_off', ''],
            45 : ['user_name_level', '']
        }

        if flask.request.method == 'POST':
            for i in setting_list:
                curs.execute(db_change("update other set data = ? where name = ?"), [
                    flask.request.form.get(setting_list[i][0], setting_list[i][1]),
                    setting_list[i][0]
                ])

            conn.commit()

            admin_check(None, 'edit_set (main)')

            return redirect('/setting/main')
        else:
            d_list = {}
            for i in setting_list:
                curs.execute(db_change('select data from other where name = ?'), [setting_list[i][0]])
                db_data = curs.fetchall()
                if not db_data:
                    curs.execute(db_change('insert into other (name, data, coverage) values (?, ?, "")'), [setting_list[i][0], setting_list[i][1]])

                d_list[i] = db_data[0][0] if db_data else setting_list[i][1]
            else:
                conn.commit()

            init_set_list = get_init_set_list()
                
            # 언어도 변경 가능하도록 필요
                
            encode_select = ''
            encode_select_data = init_set_list['encode']['list'] + ['sha256']
            for encode_select_one in encode_select_data:
                if encode_select_one == d_list[15]:
                    encode_select = '<option value="' + encode_select_one + '">' + encode_select_one + '</option>' + encode_select
                else:
                    encode_select += '<option value="' + encode_select_one + '">' + encode_select_one + '</option>'
                    
            tls_select = ''
            tls_select_data = ['http', 'https']
            for tls_select_one in tls_select_data:
                if tls_select_one == d_list[27]:
                    tls_select = '<option value="' + tls_select_one + '">' + tls_select_one + '</option>' + tls_select
                else:
                    tls_select += '<option value="' + tls_select_one + '">' + tls_select_one + '</option>'

            check_box_div = [7, 8, '', 20, 23, 24, 25, 26, 31, 33, 34, 35, 36, 37, 44, 45]
            for i in range(0, len(check_box_div)):
                acl_num = check_box_div[i]
                if acl_num != '' and d_list[acl_num]:
                    check_box_div[i] = 'checked="checked"'
                else:
                    check_box_div[i] = ''

            branch_div = ''
            branch_list = ['stable', 'dev', 'beta']
            for i in branch_list:
                if d_list[12] == i:
                    branch_div = '<option value="' + i + '">' + i + '</option>' + branch_div
                else:
                    branch_div += '<option value="' + i + '">' + i + '</option>'

            sqlite_only = 'style="display:none;"' if db_set != 'sqlite' else ''

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('main_setting'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = render_simple_set('''
                    <form method="post">
                        <h2>''' + load_lang('basic_set') + '''</h2>
                        
                        <span>''' + load_lang('wiki_name') + '''</span>
                        <hr class="main_hr">
                        <input name="name" value="''' + html.escape(d_list[0]) + '''">
                        <hr class="main_hr">

                        <span><a href="/setting/main/logo">(''' + load_lang('wiki_logo') + ''')</a></span>
                        <hr class="main_hr">

                        <span>''' + load_lang('main_page') + '''</span>
                        <hr class="main_hr">
                        <input name="frontpage" value="''' + html.escape(d_list[2]) + '''">
                        <hr class="main_hr">

                        <span>''' + load_lang('tls_method') + '''</span>
                        <hr class="main_hr">
                        <select name="http_select">''' + tls_select + '''</select>
                        <hr class="main_hr">

                        <span>''' + load_lang('domain') + '''</span> (EX : 2du.pythonanywhere.com) (''' + load_lang('off') + ''' : ''' + load_lang('empty') + ''')
                        <hr class="main_hr">
                        <input name="domain" value="''' + html.escape(d_list[22]) + '''">
                        <hr class="main_hr">

                        <span>''' + load_lang('wiki_host') + '''</span>
                        <hr class="main_hr">
                        <input name="host" value="''' + html.escape(d_list[16]) + '''">
                        <hr class="main_hr">

                        <span>''' + load_lang('wiki_port') + '''</span>
                        <hr class="main_hr">
                        <input name="port" value="''' + html.escape(d_list[10]) + '''">
                        <hr class="main_hr">

                        <span>''' + load_lang('wiki_secret_key') + '''</span>
                        <hr class="main_hr">
                        <input type="password" name="key" value="''' + html.escape(d_list[11]) + '''">
                        <hr class="main_hr">

                        <span>''' + load_lang('encryption_method') + '''</span>
                        <hr class="main_hr">
                        <select name="encode">''' + encode_select + '''</select>
                        <hr class="main_hr">
                        
                        <input type="checkbox" name="wiki_access_password_need" ''' + check_box_div[8] + '''> ''' + load_lang('set_wiki_access_password_need') + ''' (''' + load_lang('restart_required') + ''')
                        <hr class="main_hr">
                        
                        <span>''' + load_lang('set_wiki_access_password') + '''</span> (''' + load_lang('restart_required') + ''')
                        <hr class="main_hr">
                        <input type="password" name="wiki_access_password" value="''' + html.escape(d_list[32]) + '''">
                        
                        <h3>''' + load_lang('authority_use_list') + '''</h3>
                        
                        <input type="checkbox" name="auth_history_off" ''' + check_box_div[14] + '''> ''' + load_lang('authority_use_list_off') + '''
                        <hr class="main_hr">
                        
                        <span>''' + load_lang('authority_use_list_expiration_date') + '''</span> (''' + load_lang('day') + ''') (''' + load_lang('off') + ''' : ''' + load_lang('empty') + ''')
                        <hr class="main_hr">
                        <input name="auth_history_expiration_date" value="''' + html.escape(d_list[43]) + '''">
                        <hr class="main_hr">

                        <h3>''' + load_lang('communication_set') + '''</h3>
                        
                        <input type="checkbox" name="enable_comment" ''' + check_box_div[5] + '''> ''' + load_lang('enable_comment_function') + ''' (''' + load_lang('not_working') + ''')
                        <hr class="main_hr">

                        <input type="checkbox" name="enable_challenge" ''' + check_box_div[6] + '''> ''' + load_lang('enable_challenge_function') + ''' (''' + load_lang('not_working') + ''')
                        <hr class="main_hr">

                        <input type="checkbox" name="user_name_level" ''' + check_box_div[15] + '''> ''' + load_lang('display_level_in_user_name') + '''
                        <hr class="main_hr">

                        <h2>''' + load_lang('design_set') + '''</h2>
                        
                        <span>''' + load_lang('wiki_skin') + '''</span>
                        <hr class="main_hr">
                        <select name="skin">''' + load_skin(d_list[5] if d_list[5] != '' else 'ringo') + '''</select>

                        <h2>''' + load_lang('render_set') + '''</h2>
                        
                        <input type="checkbox" name="namumark_compatible" ''' + check_box_div[10] + '''> ''' + load_lang('namumark_fully_compatible_mode') + '''
                        <hr class="main_hr">
                        
                        <input type="checkbox" name="link_case_insensitive" ''' + check_box_div[12] + '''> ''' + load_lang('link_case_insensitive') + '''
                        <hr class="main_hr">

                        <h2>''' + load_lang('login_set') + '''</h2>
                        
                        <input type="checkbox" name="reg" ''' + check_box_div[0] + '''> ''' + load_lang('no_register') + '''
                        <hr class="main_hr">

                        <input type="checkbox" name="ip_view" ''' + check_box_div[1] + '''> ''' + load_lang('hide_ip') + '''
                        <hr class="main_hr">

                        <input type="checkbox" name="user_name_view" ''' + check_box_div[11] + '''> ''' + load_lang('hide_user_name') + '''
                        <hr class="main_hr">

                        <input type="checkbox" name="requires_approval" ''' + check_box_div[3] + '''> ''' + load_lang('requires_approval') + '''
                        <hr class="main_hr">

                        <span>''' + load_lang('password_min_length') + '''</span> (''' + load_lang('off') + ''' : ''' + load_lang('empty') + ''')
                        <hr class="main_hr">
                        <input name="password_min_length" value="''' + html.escape(d_list[30]) + '''">

                        <h3>''' + load_lang('ua') + '''</h3>
                        
                        <input type="checkbox" name="ua_get" ''' + check_box_div[4] + '''> ''' + load_lang('ua_get_off') + '''
                        <hr class="main_hr">
                        
                        <span>''' + load_lang('ua_expiration_date') + '''</span> (''' + load_lang('day') + ''') (''' + load_lang('off') + ''' : ''' + load_lang('empty') + ''')
                        <hr class="main_hr">
                        <input name="ua_expiration_date" value="''' + html.escape(d_list[42]) + '''">
                        <hr class="main_hr">
                        
                        <h2>''' + load_lang('server_set') + '''</h2>

                        <span>''' + load_lang('update_branch') + '''</span>
                        <hr class="main_hr">
                        <select name="update">''' + branch_div + '''</select>

                        <span ''' + sqlite_only + '''>
                            <h3>''' + load_lang('backup') + ''' (''' + load_lang('sqlite_only') + ''')</h3>
                            
                            <span>''' + load_lang('backup_warning') + ''' (EX : data_YYYYMMDDHHMMSS.db)</span>
                            <hr class="main_hr">
                            <hr class="main_hr">
                            
                            <span>''' + load_lang('backup_interval') + '''</span> (''' + load_lang('hour') + ''') (''' + load_lang('off') + ''' : ''' + load_lang('empty') + ''')
                            <hr class="main_hr">
                            <input name="back_up" value="''' + html.escape(d_list[9]) + '''">
                            <hr class="main_hr">
                            
                            <span>''' + load_lang('backup_where') + '''</span> (''' + load_lang('default') + ''' : ''' + load_lang('empty') + ''') (''' + load_lang('example') + ''' : ./data/backup.db)
                            <hr class="main_hr">
                            <input name="backup_where" value="''' + html.escape(d_list[21]) + '''">
                            <hr class="main_hr">

                            <span>''' + load_lang('backup_count') + '''</span> (''' + load_lang('default') + ''' : ''' + load_lang('empty') + ''')
                            <hr class="main_hr">
                            <input name="backup_count" value="''' + html.escape(d_list[41]) + '''">
                            <hr class="main_hr">
                        </span>

                        <h2>''' + load_lang('edit_set') + '''</h2>
                        
                        <span>''' + load_lang('slow_edit') + '''</span> (''' + load_lang('second') + ''') (''' + load_lang('off') + ''' : ''' + load_lang('empty') + ''')
                        <hr class="main_hr">
                        <input name="slow_edit" value="''' + html.escape(d_list[19]) + '''">
                        <hr class="main_hr">
                        
                        <input type="checkbox" name="edit_bottom_compulsion" ''' + check_box_div[7] + '''> ''' + load_lang('edit_bottom_compulsion') + '''
                        <hr class="main_hr">
                        
                        <span>''' + load_lang('title_max_length') + '''</span> (''' + load_lang('off') + ''' : ''' + load_lang('empty') + ''')
                        <hr class="main_hr">
                        <input name="title_max_length" value="''' + html.escape(d_list[28]) + '''">
                        <hr class="main_hr">
                        
                        <span>''' + load_lang('title_topic_max_length') + '''</span> (''' + load_lang('off') + ''' : ''' + load_lang('empty') + ''')
                        <hr class="main_hr">
                        <input name="title_topic_max_length" value="''' + html.escape(d_list[29]) + '''">
                        <hr class="main_hr">
                        
                        <span>''' + load_lang('max_file_size') + ''' (MB)</span>
                        <hr class="main_hr">
                        <input name="upload" value="''' + html.escape(d_list[4]) + '''">
                        <hr class="main_hr">
                        
                        <input type="checkbox" name="history_recording_off" ''' + check_box_div[9] + '''> ''' + load_lang('set_history_recording_off') + '''
                        <hr class="main_hr">

                        <input type="checkbox" name="move_with_redirect" ''' + check_box_div[13] + '''> ''' + load_lang('move_with_redirect') + ''' (''' + load_lang('not_working') + ''')
                        <hr class="main_hr">

                        <span>''' + load_lang('slow_thread') + '''</span> (''' + load_lang('second') + ''') (''' + load_lang('off') + ''' : ''' + load_lang('empty') + ''')
                        <hr class="main_hr">
                        <input name="slow_thread" value="''' + html.escape(d_list[38]) + '''">
                        <hr class="main_hr">

                        <span>''' + load_lang('edit_timeout') + '''</span> (''' + load_lang('second') + ''') (''' + load_lang('off') + ''' : ''' + load_lang('empty') + ''') (''' + load_lang('linux_only') + ''')
                        <hr class="main_hr">
                        <input name="edit_timeout" value="''' + html.escape(d_list[39]) + '''">
                        <hr class="main_hr">

                        <span>''' + load_lang('document_content_max_length') + '''</span> (''' + load_lang('off') + ''' : ''' + load_lang('empty') + ''')
                        <hr class="main_hr">
                        <input name="document_content_max_length" value="''' + html.escape(d_list[40]) + '''">
                        <hr class="main_hr">

                        <button id="opennamu_save_button" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                '''),
                menu = [['setting', load_lang('return')]]
            ))