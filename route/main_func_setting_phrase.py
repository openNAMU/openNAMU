from .tool.func import *

def main_func_setting_phrase():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/ban')
        
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
            'topic_text'
        ]
        if flask.request.method == 'POST':
            for i in i_list:
                curs.execute(db_change("update other set data = ? where name = ?"), [
                    flask.request.form.get(i, ''),
                    i
                ])

            conn.commit()

            admin_check(None, 'edit_set (phrase)')

            return redirect('/setting/phrase')
        else:
            d_list = []

            for i in i_list:
                curs.execute(db_change('select data from other where name = ?'), [i])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list += [sql_d[0][0]]
                else:
                    curs.execute(db_change('insert into other (name, data) values (?, ?)'), [i, ''])

                    d_list += ['']

            conn.commit()

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('text_setting'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post" id="main_set_data">
                        <h2>1. ''' + load_lang('register_text') + ''' (HTML)</h2>
                        <textarea rows="3" name="''' + i_list[0] + '''">''' + html.escape(d_list[0]) + '''</textarea>

                        <h2>2. ''' + load_lang('non_login_alert') + ''' (HTML)</h2>
                        <textarea rows="3" name="''' + i_list[1] + '''">''' + html.escape(d_list[1]) + '''</textarea>

                        <h2>3. ''' + load_lang('edit_bottom_text') + ''' (HTML)</h2>
                        <textarea rows="3" name="''' + i_list[2] + '''">''' + html.escape(d_list[2]) + '''</textarea>

                        <h2>4. ''' + load_lang('copyright_checkbox_text') + ''' (HTML)</h2>
                        <textarea rows="3" name="''' + i_list[3] + '''">''' + html.escape(d_list[3]) + '''</textarea>

                        <h2>5. ''' + load_lang('check_key_text') + ''' (HTML)</h2>
                        <textarea rows="3" name="''' + i_list[4] + '''">''' + html.escape(d_list[4]) + '''</textarea>

                        <h2>6. ''' + load_lang('email_title') + '''</h2>
                        <textarea rows="3" name="''' + i_list[5] + '''">''' + html.escape(d_list[5]) + '''</textarea>

                        <h2>7. ''' + load_lang('email_text') + '''</h2>
                        <textarea rows="3" name="''' + i_list[6] + '''">''' + html.escape(d_list[6]) + '''</textarea>

                        <h2>8. ''' + load_lang('email_insert_text') + '''</h2>
                        <textarea rows="3" name="''' + i_list[7] + '''">''' + html.escape(d_list[7]) + '''</textarea>

                        <h2>9. ''' + load_lang('password_search_text') + '''</h2>
                        <textarea rows="3" name="''' + i_list[8] + '''">''' + html.escape(d_list[8]) + '''</textarea>

                        <h2>10. ''' + load_lang('reset_user_text') + '''</h2>
                        <textarea rows="3" name="''' + i_list[9] + '''">''' + html.escape(d_list[9]) + '''</textarea>

                        <h2>11. ''' + load_lang('error_401') + '''</h2>
                        <textarea rows="3" name="''' + i_list[10] + '''">''' + html.escape(d_list[10]) + '''</textarea>

                        <h2>12. ''' + load_lang('error_404') + '''</h2>
                        <textarea rows="3" name="''' + i_list[11] + '''">''' + html.escape(d_list[11]) + '''</textarea>

                        <h2>13. ''' + load_lang('approval_question') + '''</h2>
                        <sup>(1)</sup>
                        <hr class="main_hr">
                        <textarea rows="3" name="''' + i_list[12] + '''">''' + html.escape(d_list[12]) + '''</textarea>

                        <h2>14. ''' + load_lang('edit_help') + '''</h2>
                        <textarea rows="3" name="''' + i_list[13] + '''">''' + html.escape(d_list[13]) + '''</textarea>

                        <h2>15. ''' + load_lang('upload_help') + ''' (HTML)</h2>
                        <textarea rows="3" name="''' + i_list[14] + '''">''' + html.escape(d_list[14]) + '''</textarea>

                        <h2>16. ''' + load_lang('upload_default') + '''</h2>
                        <textarea rows="3" name="''' + i_list[15] + '''">''' + html.escape(d_list[15]) + '''</textarea>

                        <h2>17. ''' + load_lang('bottom_text') + ''' (HTML)</h2>
                        <textarea rows="3" name="''' + i_list[16] + '''">''' + html.escape(d_list[16]) + '''</textarea>

                        <h2>18. ''' + load_lang('topic_text') + '''</h2>
                        <textarea rows="3" name="''' + i_list[17] + '''">''' + html.escape(d_list[17]) + '''</textarea>

                        <hr class="main_hr">
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                    <ul id="footnote_data">
                        <li><a href="#note_1" id="note_1_end">(1)</a> ''' + load_lang('approval_question_visible_only_when_approval_on') + '''</li>
                    </ul>
                    <script>simple_render('main_set_data');</script>
                ''',
                menu = [['setting', load_lang('return')]]
            ))