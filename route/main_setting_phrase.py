from .tool.func import *

async def main_setting_phrase():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 0)
        
        i_list = [
            'contract',
            'no_login_warning',
            'edit_bottom_text',
            'copyright_checkbox_text',
            'check_key_text',
            'email_title',
            'email_text',
            'email_insert_text',
            'password_search_text',
            'reset_user_text',
            'error_401',
            'error_404',
            'approval_question',
            'edit_help',
            'upload_help',
            'upload_default',
            'license',
            'topic_text',
            'phrase_user_page_admin',
            'phrase_user_page_owner',
            'phrase_old_page_warning',
            'bbs_help',
            'bbs_comment_help',
            'outdated_doc_warning',
            'outdated_doc_warning_date',
            'category_text',
            'redirect_text',
            'template_var_1',
            'template_var_2',
            'template_var_3',
            'edit_only_bottom_text',
            'move_bottom_text',
            'delete_bottom_text',
            'revert_bottom_text',
        ]
        if flask.request.method == 'POST':
            curs.executemany(db_change("update other set data = ? where name = ?"), [[flask.request.form.get(for_a, ''), for_a] for for_a in i_list])

            await acl_check(tool = 'owner_auth', memo = 'edit_set (phrase)')

            return redirect(conn, '/setting/phrase')
        else:
            d_list = []
            for i in i_list:
                curs.execute(db_change('select data from other where name = ?'), [i])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list += [sql_d[0][0]]
                else:
                    curs.execute(db_change('insert into other (name, data, coverage) values (?, ?, "")'), [i, ''])
                    d_list += ['']

            return easy_minify(flask.render_template(await skin_check(conn),
                imp = [await get_lang('text_setting'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = await render_simple_set('''
                    <form method="post">
                        <h2>''' + await get_lang('register_text') + ''' (HTML)</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[0] + '''">''' + html.escape(d_list[0]) + '''</textarea>

                        <h2>''' + await get_lang('non_login_alert') + ''' (HTML)</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[1] + '''">''' + html.escape(d_list[1]) + '''</textarea>

                        <h2>''' + await get_lang('copyright_checkbox_text') + ''' (HTML)</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[3] + '''">''' + html.escape(d_list[3]) + '''</textarea>

                        <h2>''' + await get_lang('check_key_text') + ''' (HTML)</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[4] + '''">''' + html.escape(d_list[4]) + '''</textarea>

                        <h2>''' + await get_lang('email_title') + '''</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[5] + '''">''' + html.escape(d_list[5]) + '''</textarea>

                        <h2>''' + await get_lang('email_text') + '''</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[6] + '''">''' + html.escape(d_list[6]) + '''</textarea>

                        <h2>''' + await get_lang('email_insert_text') + '''</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[7] + '''">''' + html.escape(d_list[7]) + '''</textarea>

                        <h2>''' + await get_lang('password_search_text') + '''</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[8] + '''">''' + html.escape(d_list[8]) + '''</textarea>

                        <h2>''' + await get_lang('reset_user_text') + '''</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[9] + '''">''' + html.escape(d_list[9]) + '''</textarea>

                        <h2>''' + await get_lang('error_401') + '''</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[10] + '''">''' + html.escape(d_list[10]) + '''</textarea>

                        <h2>''' + await get_lang('error_404') + '''</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[11] + '''">''' + html.escape(d_list[11]) + '''</textarea>

                        <h2>''' + await get_lang('approval_question') + '''</h2>
                        <sup><a href="/setting/main">''' + await get_lang('approval_question_visible_only_when_approval_on') + '''</a></sup>
                        <hr class="main_hr">
                        <textarea class="opennamu_textarea_100" name="''' + i_list[12] + '''">''' + html.escape(d_list[12]) + '''</textarea>

                        <h2>''' + await get_lang('edit_help') + '''</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[13] + '''">''' + html.escape(d_list[13]) + '''</textarea>

                        <h2>''' + await get_lang('upload_help') + ''' (HTML)</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[14] + '''">''' + html.escape(d_list[14]) + '''</textarea>

                        <h2>''' + await get_lang('upload_default') + '''</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[15] + '''">''' + html.escape(d_list[15]) + '''</textarea>

                        <h2>''' + await get_lang('bottom_text') + ''' (HTML)</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[16] + '''">''' + html.escape(d_list[16]) + '''</textarea>

                        <h2>''' + await get_lang('topic_text') + '''</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[17] + '''">''' + html.escape(d_list[17]) + '''</textarea>
                        
                        <h2>''' + await get_lang('phrase_user_page_admin') + ''' (HTML)</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[18] + '''">''' + html.escape(d_list[18]) + '''</textarea>
                        
                        <h2>''' + await get_lang('phrase_user_page_owner') + ''' (HTML)</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[19] + '''">''' + html.escape(d_list[19]) + '''</textarea>

                        <h2>''' + await get_lang('phrase_old_page_warning') + ''' (''' + await get_lang('beta') + ''') (HTML)</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[20] + '''">''' + html.escape(d_list[20]) + '''</textarea>
                        
                        <h2>''' + await get_lang('bbs_help') + '''</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[21] + '''">''' + html.escape(d_list[21]) + '''</textarea>

                        <h2>''' + await get_lang('bbs_comment_help') + '''</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[22] + '''">''' + html.escape(d_list[22]) + '''</textarea>

                        <h2>''' + await get_lang('outdated_doc_warning') + '''  (HTML)</h2>
                        <span>''' + await get_lang('period') + '''</span> (''' + await get_lang('day') + ''') (''' + await get_lang('off') + ''' : ''' + await get_lang('empty') + ''')
                        <hr class="main_hr">
                        <input name="''' + i_list[24] + '''" value="''' + html.escape(d_list[24]) + '''">
                        <hr class="main_hr">
                        <textarea class="opennamu_textarea_100" name="''' + i_list[23] + '''" placeholder="''' + await get_lang('old_page_warning') + '''">''' + html.escape(d_list[23]) + '''</textarea>

                        <h2>''' + await get_lang('category') + '''</h2>
                        <input name="''' + i_list[25] + '''" value="''' + html.escape(d_list[25]) + '''">

                        <h2>''' + await get_lang('redirect') + '''</h2>
                        <span>EX : {0} ➤ {1}</span>
                        <hr class="main_hr">
                        <input name="''' + i_list[26] + '''" value="''' + html.escape(d_list[26]) + '''">

                        <h2>''' + await get_lang('template_var') + '''</h2>
                        <h3>''' + await get_lang('template_var_1') + ''' (''' + await get_lang('default') + ''' : ''' + await get_lang('top') + ''') (HTML)</h3>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[27] + '''">''' + html.escape(d_list[27]) + '''</textarea>

                        <h3>''' + await get_lang('template_var_2') + ''' (''' + await get_lang('default') + ''' : ''' + await get_lang('sidebar') + ''') (HTML)</h3>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[28] + '''">''' + html.escape(d_list[28]) + '''</textarea>

                        <h3>''' + await get_lang('template_var_3') + ''' (''' + await get_lang('default') + ''' : ''' + await get_lang('bottom') + ''') (HTML)</h3>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[29] + '''">''' + html.escape(d_list[29]) + '''</textarea>

                        <h2>''' + await get_lang('edit_bottom_text') + ''' (HTML)</h2>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[2] + '''">''' + html.escape(d_list[2]) + '''</textarea>

                        <h3>''' + await get_lang('edit_only_bottom_text') + ''' (HTML)</h3>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[30] + '''">''' + html.escape(d_list[30]) + '''</textarea>

                        <h3>''' + await get_lang('move_bottom_text') + ''' (HTML)</h3>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[31] + '''">''' + html.escape(d_list[31]) + '''</textarea>

                        <h3>''' + await get_lang('delete_bottom_text') + ''' (HTML)</h3>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[32] + '''">''' + html.escape(d_list[32]) + '''</textarea>

                        <h3>''' + await get_lang('revert_bottom_text') + ''' (HTML)</h3>
                        <textarea class="opennamu_textarea_100" name="''' + i_list[33] + '''">''' + html.escape(d_list[33]) + '''</textarea>

                        <hr class="main_hr">
                        <button id="opennamu_save_button" type="submit">''' + await get_lang('save') + '''</button>
                    </form>
                '''),
                menu = [['setting', await get_lang('return')]]
            ))