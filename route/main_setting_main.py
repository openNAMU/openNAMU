from .tool.func import *

async def main_setting_main():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 0)
        
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
            45 : ['user_name_level', ''],
            46 : ['load_ip_select', ''],
            47 : ['not_use_view_count', '']
        }

        if flask.request.method == 'POST':
            for i in setting_list:
                curs.execute(db_change("update other set data = ? where name = ?"), [
                    flask.request.form.get(setting_list[i][0], setting_list[i][1]),
                    setting_list[i][0]
                ])

            await acl_check(tool = 'owner_auth', memo = 'edit_set (main)')

            return redirect(conn, '/setting/main')
        else:
            d_list = {}
            for i in setting_list:
                curs.execute(db_change('select data from other where name = ?'), [setting_list[i][0]])
                db_data = curs.fetchall()
                if not db_data:
                    curs.execute(db_change('insert into other (name, data, coverage) values (?, ?, "")'), [setting_list[i][0], setting_list[i][1]])

                d_list[i] = db_data[0][0] if db_data else setting_list[i][1]

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

            check_box_div = [7, 8, '', 20, 23, 24, '', 26, 31, 33, 34, 35, 36, 37, 44, 45, 47]
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

            set_data = global_some_set_do('db_type')
            
            sqlite_only = ''
            if set_data != 'sqlite':
                sqlite_only = 'style="display:none;"'

            ip_load_select_data = ''
            ip_load_option = ['default', 'HTTP_X_REAL_IP', 'HTTP_CF_CONNECTING_IP', 'REMOTE_ADDR']
            for for_a in ip_load_option:
                view_ip_option = for_a
                if for_a == 'default':
                    view_ip_option = get_lang(conn, 'default')

                if d_list[46] == for_a:
                    ip_load_select_data = '<option value="' + for_a + '">' + view_ip_option + '</option>' + ip_load_select_data
                else:
                    ip_load_select_data += '<option value="' + for_a + '">' + view_ip_option + '</option>'

            basic_set = '''
                <h2>''' + get_lang(conn, 'basic_set') + '''</h2>
                            
                <span>''' + get_lang(conn, 'wiki_name') + '''</span>
                <hr class="main_hr">
                <input name="name" value="''' + html.escape(d_list[0]) + '''">
                <hr class="main_hr">

                <span><a href="/setting/main/logo">(''' + get_lang(conn, 'wiki_logo') + ''')</a></span>
                <hr class="main_hr">

                <span>''' + get_lang(conn, 'main_page') + '''</span>
                <hr class="main_hr">
                <input name="frontpage" value="''' + html.escape(d_list[2]) + '''">
                <hr class="main_hr">

                <span>''' + get_lang(conn, 'tls_method') + '''</span>
                <hr class="main_hr">
                <select name="http_select">''' + tls_select + '''</select>
                <hr class="main_hr">

                <span>''' + get_lang(conn, 'domain') + '''</span> (EX : 2du.pythonanywhere.com) (''' + get_lang(conn, 'off') + ''' : ''' + get_lang(conn, 'empty') + ''')
                <hr class="main_hr">
                <input name="domain" value="''' + html.escape(d_list[22]) + '''">
                <hr class="main_hr">

                <span>''' + get_lang(conn, 'wiki_host') + '''</span>
                <hr class="main_hr">
                <input name="host" value="''' + html.escape(d_list[16]) + '''">
                <hr class="main_hr">

                <span>''' + get_lang(conn, 'wiki_port') + '''</span>
                <hr class="main_hr">
                <input name="port" value="''' + html.escape(d_list[10]) + '''">
                <hr class="main_hr">

                <span>''' + get_lang(conn, 'wiki_secret_key') + '''</span>
                <hr class="main_hr">
                <input type="password" name="key" value="''' + html.escape(d_list[11]) + '''">
                <hr class="main_hr">
                
                <label><input type="checkbox" name="wiki_access_password_need" ''' + check_box_div[8] + '''> ''' + get_lang(conn, 'set_wiki_access_password_need') + ''' (''' + get_lang(conn, 'restart_required') + ''')</label>
                <hr class="main_hr">
                
                <span>''' + get_lang(conn, 'set_wiki_access_password') + '''</span> (''' + get_lang(conn, 'restart_required') + ''')
                <hr class="main_hr">
                <input type="password" name="wiki_access_password" value="''' + html.escape(d_list[32]) + '''">
                <hr class="main_hr">

                <span>''' + get_lang(conn, 'wiki_load_ip_select') + '''</span> (''' + get_lang(conn, 'restart_required') + ''')
                <hr class="main_hr">
                <select name="load_ip_select">''' + ip_load_select_data + '''</select>
                
                <h3>''' + get_lang(conn, 'authority_use_list') + '''</h3>
                
                <label><input type="checkbox" name="auth_history_off" ''' + check_box_div[14] + '''> ''' + get_lang(conn, 'authority_use_list_off') + '''</label>
                <hr class="main_hr">
                
                <span>''' + get_lang(conn, 'authority_use_list_expiration_date') + '''</span> (''' + get_lang(conn, 'day') + ''') (''' + get_lang(conn, 'off') + ''' : ''' + get_lang(conn, 'empty') + ''')
                <hr class="main_hr">
                <input name="auth_history_expiration_date" value="''' + html.escape(d_list[43]) + '''">
                <hr class="main_hr">

                <h3>''' + get_lang(conn, 'communication_set') + '''</h3>
                
                <label><input type="checkbox" name="enable_comment" ''' + check_box_div[5] + '''> ''' + get_lang(conn, 'enable_comment_function') + '''</label>
                <hr class="main_hr">

                <label><input type="checkbox" name="user_name_level" ''' + check_box_div[15] + '''> ''' + get_lang(conn, 'display_level_in_user_name') + '''</label>
                <hr class="main_hr">

                <label><input type="checkbox" name="not_use_view_count" ''' + check_box_div[16] + '''> ''' + get_lang(conn, 'not_use_view_count') + '''</label>
                <hr class="main_hr">
            '''

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'main_setting'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = render_simple_set(conn, '''
                    <form method="post">
                        ''' + basic_set + '''
                        <h2>''' + get_lang(conn, 'design_set') + '''</h2>
                        
                        <span>''' + get_lang(conn, 'wiki_skin') + '''</span>
                        <hr class="main_hr">
                        <select name="skin">''' + load_skin(conn, d_list[5] if d_list[5] != '' else 'ringo') + '''</select>

                        <h2>''' + get_lang(conn, 'render_set') + '''</h2>
                        
                        <label><input type="checkbox" name="namumark_compatible" ''' + check_box_div[10] + '''> ''' + get_lang(conn, 'namumark_fully_compatible_mode') + '''</label>
                        <hr class="main_hr">
                        
                        <label><input type="checkbox" name="link_case_insensitive" ''' + check_box_div[12] + '''> ''' + get_lang(conn, 'link_case_insensitive') + '''</label>
                        <hr class="main_hr">

                        <h2>''' + get_lang(conn, 'login_set') + '''</h2>
                        
                        <label><input type="checkbox" name="reg" ''' + check_box_div[0] + '''> ''' + get_lang(conn, 'no_register') + '''</label>
                        <hr class="main_hr">

                        <label><input type="checkbox" name="ip_view" ''' + check_box_div[1] + '''> ''' + get_lang(conn, 'hide_ip') + '''</label>
                        <hr class="main_hr">

                        <label><input type="checkbox" name="user_name_view" ''' + check_box_div[11] + '''> ''' + get_lang(conn, 'hide_user_name') + '''</label>
                        <hr class="main_hr">

                        <label><input type="checkbox" name="requires_approval" ''' + check_box_div[3] + '''> ''' + get_lang(conn, 'requires_approval') + '''</label>
                        <hr class="main_hr">

                        <span>''' + get_lang(conn, 'password_min_length') + '''</span> (''' + get_lang(conn, 'off') + ''' : ''' + get_lang(conn, 'empty') + ''')
                        <hr class="main_hr">
                        <input name="password_min_length" value="''' + html.escape(d_list[30]) + '''">
                        <hr class="main_hr">

                        <span>''' + get_lang(conn, 'encryption_method') + '''</span>
                        <hr class="main_hr">
                        <select name="encode">''' + encode_select + '''</select>

                        <h3>''' + get_lang(conn, 'ua') + '''</h3>
                        
                        <label><input type="checkbox" name="ua_get" ''' + check_box_div[4] + '''> ''' + get_lang(conn, 'ua_get_off') + '''</label>
                        <hr class="main_hr">
                        
                        <span>''' + get_lang(conn, 'ua_expiration_date') + '''</span> (''' + get_lang(conn, 'day') + ''') (''' + get_lang(conn, 'off') + ''' : ''' + get_lang(conn, 'empty') + ''')
                        <hr class="main_hr">
                        <input name="ua_expiration_date" value="''' + html.escape(d_list[42]) + '''">
                        <hr class="main_hr">
                        
                        <h2>''' + get_lang(conn, 'server_set') + '''</h2>

                        <span>''' + get_lang(conn, 'update_branch') + '''</span>
                        <hr class="main_hr">
                        <select name="update">''' + branch_div + '''</select>

                        <span ''' + sqlite_only + '''>
                            <h3>''' + get_lang(conn, 'backup') + ''' (''' + get_lang(conn, 'sqlite_only') + ''')</h3>
                            
                            <span>''' + get_lang(conn, 'backup_warning') + ''' (EX : data_YYYYMMDDHHMMSS.db)</span>
                            <hr class="main_hr">
                            <hr class="main_hr">
                            
                            <span>''' + get_lang(conn, 'backup_interval') + '''</span> (''' + get_lang(conn, 'hour') + ''') (''' + get_lang(conn, 'off') + ''' : ''' + get_lang(conn, 'empty') + ''')
                            <hr class="main_hr">
                            <input name="back_up" value="''' + html.escape(d_list[9]) + '''">
                            <hr class="main_hr">
                            
                            <span>''' + get_lang(conn, 'backup_where') + '''</span> (''' + get_lang(conn, 'default') + ''' : ''' + get_lang(conn, 'empty') + ''') (''' + get_lang(conn, 'example') + ''' : ./data/backup.db)
                            <hr class="main_hr">
                            <input name="backup_where" value="''' + html.escape(d_list[21]) + '''">
                            <hr class="main_hr">

                            <span>''' + get_lang(conn, 'backup_count') + '''</span> (''' + get_lang(conn, 'default') + ''' : ''' + get_lang(conn, 'empty') + ''')
                            <hr class="main_hr">
                            <input name="backup_count" value="''' + html.escape(d_list[41]) + '''">
                            <hr class="main_hr">
                        </span>

                        <h2>''' + get_lang(conn, 'edit_set') + '''</h2>
                        
                        <span>''' + get_lang(conn, 'slow_edit') + '''</span> (''' + get_lang(conn, 'second') + ''') (''' + get_lang(conn, 'off') + ''' : ''' + get_lang(conn, 'empty') + ''')
                        <hr class="main_hr">
                        <input name="slow_edit" value="''' + html.escape(d_list[19]) + '''">
                        <hr class="main_hr">
                        
                        <label><input type="checkbox" name="edit_bottom_compulsion" ''' + check_box_div[7] + '''> ''' + get_lang(conn, 'edit_bottom_compulsion') + '''</label>
                        <hr class="main_hr">
                        
                        <span>''' + get_lang(conn, 'title_max_length') + '''</span> (''' + get_lang(conn, 'off') + ''' : ''' + get_lang(conn, 'empty') + ''')
                        <hr class="main_hr">
                        <input name="title_max_length" value="''' + html.escape(d_list[28]) + '''">
                        <hr class="main_hr">
                        
                        <span>''' + get_lang(conn, 'title_topic_max_length') + '''</span> (''' + get_lang(conn, 'off') + ''' : ''' + get_lang(conn, 'empty') + ''')
                        <hr class="main_hr">
                        <input name="title_topic_max_length" value="''' + html.escape(d_list[29]) + '''">
                        <hr class="main_hr">
                        
                        <span>''' + get_lang(conn, 'max_file_size') + ''' (MB)</span>
                        <hr class="main_hr">
                        <input name="upload" value="''' + html.escape(d_list[4]) + '''">
                        <hr class="main_hr">
                        
                        <label><input type="checkbox" name="history_recording_off" ''' + check_box_div[9] + '''> ''' + get_lang(conn, 'set_history_recording_off') + '''</label>
                        <hr class="main_hr">

                        <label><input type="checkbox" name="move_with_redirect" ''' + check_box_div[13] + '''> ''' + get_lang(conn, 'move_with_redirect') + ''' (''' + get_lang(conn, 'not_working') + ''')</label>
                        <hr class="main_hr">

                        <span>''' + get_lang(conn, 'slow_thread') + '''</span> (''' + get_lang(conn, 'second') + ''') (''' + get_lang(conn, 'off') + ''' : ''' + get_lang(conn, 'empty') + ''')
                        <hr class="main_hr">
                        <input name="slow_thread" value="''' + html.escape(d_list[38]) + '''">
                        <hr class="main_hr">

                        <span>''' + get_lang(conn, 'edit_timeout') + '''</span> (''' + get_lang(conn, 'second') + ''') (''' + get_lang(conn, 'off') + ''' : ''' + get_lang(conn, 'empty') + ''') (''' + get_lang(conn, 'linux_only') + ''')
                        <hr class="main_hr">
                        <input name="edit_timeout" value="''' + html.escape(d_list[39]) + '''">
                        <hr class="main_hr">

                        <span>''' + get_lang(conn, 'document_content_max_length') + '''</span> (''' + get_lang(conn, 'off') + ''' : ''' + get_lang(conn, 'empty') + ''')
                        <hr class="main_hr">
                        <input name="document_content_max_length" value="''' + html.escape(d_list[40]) + '''">
                        <hr class="main_hr">

                        <button id="opennamu_save_button" type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                '''),
                menu = [['setting', get_lang(conn, 'return')]]
            ))