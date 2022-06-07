from .tool.func import *

def main_func_setting_main(db_set):
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
            9 : ['back_up', '0'],
            10 : ['port', '3000'],
            11 : ['key', load_random_key()],
            12 : ['update', 'stable'],
            15 : ['encode', 'sha3'],
            16 : ['host', '0.0.0.0'],
            19 : ['slow_edit', '0'],
            20 : ['requires_approval', ''],
            21 : ['backup_where', ''],
            22 : ['domain', flask.request.host],
            23 : ['ua_get', ''],
            24 : ['enable_comment', ''],
            25 : ['enable_challenge', ''],
            26 : ['edit_bottom_compulsion', ''],
            27 : ['http_select', 'http'],
            28 : ['title_max_length', ''],
            29 : ['title_topic_max_length', '']
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
                    curs.execute(db_change('insert into other (name, data) values (?, ?)'), [setting_list[i][0], setting_list[i][1]])

                d_list[i] = db_data[0][0] if db_data else setting_list[i][1]
            else:
                conn.commit()

            encode_select = ''
            encode_select_data = ['sha256', 'sha3']
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

            check_box_div = ['', '', '', '', '', '', '', '']
            for i in range(0, len(check_box_div)):
                if i == 0:
                    acl_num = 7
                elif i == 1:
                    acl_num = 8
                elif i == 3:
                    acl_num = 20
                elif i == 4:
                    acl_num = 23
                elif i == 5:
                    acl_num = 24
                elif i == 6:
                    acl_num = 25
                elif i == 7:
                    acl_num = 26

                if d_list[acl_num]:
                    check_box_div[i] = 'checked="checked"'

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
                data = '''
                    <form method="post" id="main_set_data">
                        <h2>1. ''' + load_lang('basic_set') + '''</h2>
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

                        <span>''' + load_lang('domain') + '''</span> (EX : 2du.pythonanywhere.com)
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

                        <h3>1.1. ''' + load_lang('communication_set') + '''</h3>
                        <input type="checkbox" name="enable_comment" ''' + check_box_div[5] + '''> ''' + load_lang('enable_comment_function') + ''' (''' + load_lang('not_working') + ''')
                        <hr class="main_hr">

                        <input type="checkbox" name="enable_challenge" ''' + check_box_div[6] + '''> ''' + load_lang('enable_challenge_function') + ''' (''' + load_lang('not_working') + ''')
                        <hr class="main_hr">

                        <h2>2. ''' + load_lang('design_set') + '''</h2>
                        <span>''' + load_lang('wiki_skin') + '''</span>
                        <hr class="main_hr">
                        <select name="skin">''' + load_skin(d_list[5] if d_list[5] != '' else 'tenshi') + '''</select>

                        <h2>3. ''' + load_lang('login_set') + '''</h2>
                        <input type="checkbox" name="reg" ''' + check_box_div[0] + '''> ''' + load_lang('no_register') + '''
                        <hr class="main_hr">

                        <input type="checkbox" name="ip_view" ''' + check_box_div[1] + '''> ''' + load_lang('hide_ip') + '''
                        <hr class="main_hr">

                        <input type="checkbox" name="requires_approval" ''' + check_box_div[3] + '''> ''' + load_lang('requires_approval') + '''
                        <hr class="main_hr">

                        <input type="checkbox" name="ua_get" ''' + check_box_div[4] + '''> ''' + load_lang('ua_get_off') + '''

                        <h2>4. ''' + load_lang('server_set') + '''</h2>
                        <span>''' + load_lang('max_file_size') + ''' (MB)</span>
                        <hr class="main_hr">
                        <input name="upload" value="''' + html.escape(d_list[4]) + '''">
                        <hr class="main_hr">

                        <span>''' + load_lang('update_branch') + '''</span>
                        <hr class="main_hr">
                        <select name="update">''' + branch_div + '''</select>

                        <span ''' + sqlite_only + '''>
                            <h3>4.1. ''' + load_lang('sqlite_only') + '''</h3>
                            <span>
                                ''' + load_lang('backup_interval') + ' (' + load_lang('hour') + ') (' + load_lang('off') + ' : 0) ' + \
                                '(' + load_lang('restart_required') + ''')</span>
                            <hr class="main_hr">
                            <input name="back_up" value="''' + html.escape(d_list[9]) + '''">
                            <hr class="main_hr">

                            <span>
                                ''' + load_lang('backup_where') + ' (' + load_lang('empty') + ' : ' + load_lang('default') + ') ' + \
                                '(' + load_lang('restart_required') + ''') (''' + load_lang('example') + ''' : ./data/backup.db)
                            </span>
                            <hr class="main_hr">
                            <input name="backup_where" value="''' + html.escape(d_list[21]) + '''">
                            <hr class="main_hr">
                        </span>

                        <h2>5. ''' + load_lang('edit_set') + '''</h2>
                        <span><a href="/setting/acl">(''' + load_lang('main_acl_setting') + ''')</a></span>
                        <hr class="main_hr">

                        <span>''' + load_lang('slow_edit') + ' (' + load_lang('second') + ') (' + load_lang('off') + ''' : 0)</span>
                        <hr class="main_hr">
                        <input name="slow_edit" value="''' + html.escape(d_list[19]) + '''">
                        <hr class="main_hr">
                        
                        <input type="checkbox" name="edit_bottom_compulsion" ''' + check_box_div[7] + '''> ''' + load_lang('edit_bottom_compulsion') + ''' (''' + load_lang('beta') + ''')
                        <hr class="main_hr">
                        
                        <span>''' + load_lang('title_max_length') + ''' (''' + load_lang('beta') + ''')</span>
                        <hr class="main_hr">
                        <input name="title_max_length" value="''' + html.escape(d_list[28]) + '''">
                        <hr class="main_hr">
                        
                        <span>''' + load_lang('title_topic_max_length') + ''' (''' + load_lang('not_working') + ''')</span>
                        <hr class="main_hr">
                        <input name="title_topic_max_length" value="''' + html.escape(d_list[29]) + '''">
                        <hr class="main_hr">

                        <hr class="main_hr">
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                    <script>simple_render('main_set_data');</script>
                ''',
                menu = [['setting', load_lang('return')]]
            ))