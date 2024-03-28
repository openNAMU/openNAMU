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
            13 : 'recaptcha_one_check_five_pass_acl',
            14 : 'document_edit_request_acl'
        }
        default_list = {
            12 : 'user'
        }

        if flask.request.method == 'POST':
            if admin_check(conn, None, 'edit_set (acl)') != 1:
                return re_error(conn, '/ban')
            else:
                curs.executemany(db_change("update other set data = ? where name = ?"), [[flask.request.form.get(i_list[for_a], 'normal'), i_list[for_a]] for for_a in i_list])

                return redirect(conn, '/setting/acl')
        else:
            d_list = {}
            disable = 'disabled' if admin_check(conn) != 1 else ''
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

            acl_list = get_acl_list()
            for for_a in range(0, len(i_list)):
                for data_list in acl_list:
                    acl_div[for_a] += '<option value="' + data_list + '" ' + ('selected="selected"' if data_list == d_list[for_a + 1] else '') + '>' + (data_list if data_list != '' else 'normal') + '</option>'

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'main_acl_setting'), wiki_set(conn), wiki_custom(conn), wiki_css([0, 0])],
                data = render_simple_set(conn, '''
                    <form method="post">
                        <hr class="main_hr">
                        <a href="/acl/TEST#exp">(''' + get_lang(conn, 'reference') + ''')</a>
                        
                        <h2>''' + get_lang(conn, 'document_acl') + '''</h2>
                        <select ''' + disable + ''' name="edit">''' + acl_div[0] + '''</select>

                        <h3>''' + get_lang(conn, 'document_edit_acl') + '''</h3>
                        <select ''' + disable + ''' name="document_edit_acl">''' + acl_div[6] + '''</select>

                        <h3>''' + get_lang(conn, 'document_edit_request_acl') + '''</h3>
                        <select ''' + disable + ''' name="document_edit_request_acl">''' + acl_div[13] + '''</select>

                        <h3>''' + get_lang(conn, 'document_move_acl') + '''</h3>
                        <select ''' + disable + ''' name="document_move_acl">''' + acl_div[7] + '''</select>

                        <h3>''' + get_lang(conn, 'document_delete_acl') + '''</h3>
                        <select ''' + disable + ''' name="document_delete_acl">''' + acl_div[8] + '''</select>
                        
                        <h2>''' + get_lang(conn, 'discussion_acl') + '''</h2>
                        <select ''' + disable + ''' name="discussion">''' + acl_div[1] + '''</select>
                        
                        <h2>''' + get_lang(conn, 'upload_acl') + '''</h2>
                        <select ''' + disable + ''' name="upload_acl">''' + acl_div[2] + '''</select>
                        
                        <h3>''' + get_lang(conn, 'many_upload_acl') + '''</h3>
                        <select ''' + disable + ''' name="many_upload_acl">''' + acl_div[4] + '''</select>
                        
                        <h2>''' + get_lang(conn, 'view_acl') + '''</h2>
                        <select ''' + disable + ''' name="all_view_acl">''' + acl_div[3] + '''</select>
                        
                        <h2>''' + get_lang(conn, 'vote_acl') + '''</h2>
                        <select ''' + disable + ''' name="vote_acl">''' + acl_div[5] + '''</select>

                        <h2>''' + get_lang(conn, 'slow_edit_acl') + '''</h2>
                        <select ''' + disable + ''' name="slow_edit_acl">''' + acl_div[9] + '''</select>

                        <h2>''' + get_lang(conn, 'edit_bottom_compulsion_acl') + '''</h2>
                        <select ''' + disable + ''' name="edit_bottom_compulsion_acl">''' + acl_div[10] + '''</select>

                        <h2>''' + get_lang(conn, 'recaptcha_pass_acl') + '''</h2>
                        <select ''' + disable + ''' name="recaptcha_pass_acl">''' + acl_div[11] + '''</select>

                        <h3>''' + get_lang(conn, 'recaptcha_one_check_five_pass_acl') + '''</h3>
                        <select ''' + disable + ''' name="recaptcha_one_check_five_pass_acl">''' + acl_div[12] + '''</select>
                        
                        <hr class="main_hr">
                        <button id="opennamu_save_button" type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                '''),
                menu = [['setting/main', get_lang(conn, 'return')]]
            ))