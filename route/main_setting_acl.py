from .tool.func import *

def main_setting_acl():
    with get_db_connect() as conn:
        curs = conn.cursor()

        i_list = {
            1 : 'edit',
            2 : 'discussion',
            3 : 'upload_acl',
            4 : 'all_view_acl',
            5 : 'many_upload_acl',
            6 : 'vote_acl',
            7 : 'document_edit_acl',
            8 : 'document_move_acl',
            9 : 'document_delete_acl',
            10 : 'slow_edit_acl',
            11 : 'edit_bottom_compulsion_acl',
            12 : 'recaptcha_pass_acl',
            13 : 'recaptcha_one_check_five_pass_acl'
        }
        default_list = {
            12 : 'user'
        }

        if flask.request.method == 'POST':
            if admin_check(None, 'edit_set (acl)') != 1:
                return re_error('/ban')
            else:
                curs.executemany(db_change("update other set data = ? where name = ?"), [[flask.request.form.get(i_list[for_a], 'normal'), i_list[for_a]] for for_a in i_list])
                conn.commit()

                return redirect('/setting/acl')
        else:
            d_list = {}
            disable = 'disabled' if admin_check() != 1 else ''
            acl_div = ['' for _ in range(0, len(i_list))]

            for for_a in i_list:
                curs.execute(db_change('select data from other where name = ?'), [i_list[for_a]])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list[for_a] = sql_d[0][0]
                else:
                    default_data = 'normal' if not for_a in default_list else default_list[for_a]
                    curs.execute(db_change('insert into other (name, data, coverage) values (?, ?, "")'), [i_list[for_a], default_data])
                    d_list[for_a] = default_data

            conn.commit()

            acl_list = get_acl_list()
            for for_a in range(0, len(i_list)):
                for data_list in acl_list:
                    acl_div[for_a] += '<option value="' + data_list + '" ' + ('selected="selected"' if data_list == d_list[for_a + 1] else '') + '>' + (data_list if data_list != '' else 'normal') + '</option>'

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('main_acl_setting'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = render_simple_set('''
                    <form method="post">
                        <hr class="main_hr">
                        <a href="/acl/TEST#exp">(''' + load_lang('reference') + ''')</a>
                        
                        <h2>''' + load_lang('document_acl') + '''</h2>
                        <select ''' + disable + ''' name="edit">''' + acl_div[0] + '''</select>

                        <h3>''' + load_lang('document_edit_acl') + '''</h3>
                        <select ''' + disable + ''' name="document_edit_acl">''' + acl_div[6] + '''</select>

                        <h3>''' + load_lang('document_move_acl') + '''</h3>
                        <select ''' + disable + ''' name="document_move_acl">''' + acl_div[7] + '''</select>

                        <h3>''' + load_lang('document_delete_acl') + '''</h3>
                        <select ''' + disable + ''' name="document_delete_acl">''' + acl_div[8] + '''</select>
                        
                        <h2>''' + load_lang('discussion_acl') + '''</h2>
                        <select ''' + disable + ''' name="discussion">''' + acl_div[1] + '''</select>
                        
                        <h2>''' + load_lang('upload_acl') + '''</h2>
                        <select ''' + disable + ''' name="upload_acl">''' + acl_div[2] + '''</select>
                        
                        <h3>''' + load_lang('many_upload_acl') + '''</h3>
                        <select ''' + disable + ''' name="many_upload_acl">''' + acl_div[4] + '''</select>
                        
                        <h2>''' + load_lang('view_acl') + '''</h2>
                        <select ''' + disable + ''' name="all_view_acl">''' + acl_div[3] + '''</select>
                        
                        <h2>''' + load_lang('vote_acl') + '''</h2>
                        <select ''' + disable + ''' name="vote_acl">''' + acl_div[5] + '''</select>

                        <h2>''' + load_lang('slow_edit_acl') + '''</h2>
                        <select ''' + disable + ''' name="slow_edit_acl">''' + acl_div[9] + '''</select>

                        <h2>''' + load_lang('edit_bottom_compulsion_acl') + '''</h2>
                        <select ''' + disable + ''' name="edit_bottom_compulsion_acl">''' + acl_div[10] + '''</select>

                        <h2>''' + load_lang('recaptcha_pass_acl') + '''</h2>
                        <select ''' + disable + ''' name="recaptcha_pass_acl">''' + acl_div[11] + '''</select>

                        <h3>''' + load_lang('recaptcha_one_check_five_pass_acl') + '''</h3>
                        <select ''' + disable + ''' name="recaptcha_one_check_five_pass_acl">''' + acl_div[12] + '''</select>
                        
                        <hr class="main_hr">
                        <button id="opennamu_save_button" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                '''),
                menu = [['setting/main', load_lang('return')]]
            ))