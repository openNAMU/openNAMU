from .tool.func import *

from .user_setting_skin_set_main import user_setting_skin_set_main_set_list

async def main_setting_skin_set():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 0)
            
        set_list = await user_setting_skin_set_main_set_list()

        if flask.request.method == 'POST':
            for for_b in set_list:
                curs.execute(db_change('select data from other where name = ?'), [for_b])
                if curs.fetchall():
                    curs.execute(db_change("update other set data = ? where name = ?"), [flask.request.form.get(for_b, set_list[for_b][0][0]), for_b])
                else:
                    curs.execute(db_change('insert into other (name, data, coverage) values (?, ?, "")'), [for_b, flask.request.form.get(for_b, set_list[for_b][0][0])])

            await acl_check(tool = 'owner_auth', memo = 'edit_set (skin_set)')

            return redirect(conn, '/setting/skin_set')
        else:
            set_data = {}
            for for_b in set_list:
                set_data[for_b] = ''

                curs.execute(db_change('select data from other where name = ?'), [for_b])
                db_data = curs.fetchall()
                get_data = db_data[0][0] if db_data else ''

                for for_a in set_list[for_b]:
                    if get_data == for_a[0]:
                        set_data[for_b] = '<option value="' + for_a[0] + '">' + for_a[1] + '</option>' + set_data[for_b]
                    else:
                        set_data[for_b] += '<option value="' + for_a[0] + '">' + for_a[1] + '</option>'

            set_data_main = { for_b : '' for for_b in set_list }

            return await render_template(
                await get_lang('main_skin_set_default'),
                await render_simple_set('''
                    <form method="post">
                        <h2>''' + await get_lang("render") + '''</h2>
                        <h3>''' + await get_lang("strike") + '''</h3>
                        ''' + set_data_main["main_css_strike"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_strike">
                            ''' + set_data["main_css_strike"] + '''
                        </select></span>
                        <h3>''' + await get_lang("bold") + '''</h3>
                        ''' + set_data_main["main_css_bold"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_bold">
                            ''' + set_data["main_css_bold"] + '''
                        </select></span>
                        <h3>''' + await get_lang("category") + '''</h3>
                        <h4>''' + await get_lang("position") + '''</h4>
                        ''' + set_data_main["main_css_category_set"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_category_set">
                            ''' + set_data["main_css_category_set"] + '''
                        </select></span>
                        <h4>''' + await get_lang("category_change_title") + '''</h4>
                        ''' + set_data_main["main_css_category_change_title"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_category_change_title">
                            ''' + set_data["main_css_category_change_title"] + '''
                        </select></span>
                        <h3>''' + await get_lang("footnote") + ''' (''' + await get_lang('beta') + ''')</h3>
                        <h4>''' + await get_lang("footnote_render") + '''</h4>
                        ''' + set_data_main["main_css_footnote_set"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_footnote_set">
                            ''' + set_data["main_css_footnote_set"] + '''
                        </select></span>
                        <h4>''' + await get_lang("footnote_number") + '''</h4>
                        ''' + set_data_main["main_css_footnote_number"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_footnote_number">
                            ''' + set_data["main_css_footnote_number"] + '''
                        </select></span>
                        <h4>''' + await get_lang("footnote_real_num_view") + '''</h4>
                        ''' + set_data_main["main_css_view_real_footnote_num"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_view_real_footnote_num">
                            ''' + set_data["main_css_view_real_footnote_num"] + '''
                        </select></span>
                        <h3>''' + await get_lang("include_link") + '''</h3>
                        ''' + set_data_main["main_css_include_link"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_include_link">
                            ''' + set_data["main_css_include_link"] + '''
                        </select></span>
                        <h3>''' + await get_lang("image") + ''' (''' + await get_lang('beta') + ''')</h3>
                        ''' + set_data_main["main_css_image_set"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_image_set">
                            ''' + set_data["main_css_image_set"] + '''
                        </select></span>
                        <h3>''' + await get_lang("toc") + '''</h3>
                        ''' + set_data_main["main_css_toc_set"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_toc_set">
                            ''' + set_data["main_css_toc_set"] + '''
                        </select></span>
                        <h3>''' + await get_lang("exter_link") + '''</h3>
                        ''' + set_data_main["main_css_exter_link"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_exter_link">
                            ''' + set_data["main_css_exter_link"] + '''
                        </select></span>
                        <h3>''' + await get_lang("link_delimiter") + '''</h3>
                        ''' + set_data_main["main_css_link_delimiter"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_link_delimiter">
                            ''' + set_data["main_css_link_delimiter"] + '''
                        </select></span>
                        <h3>''' + await get_lang("force_darkmode") + '''</h3>
                        ''' + set_data_main["main_css_darkmode"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_darkmode">
                            ''' + set_data["main_css_darkmode"] + '''
                        </select></span>
                        <h3>''' + await get_lang("table") + '''</h3>
                        <h4>''' + await get_lang("table_scroll") + '''</h4>
                        ''' + set_data_main["main_css_table_scroll"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_table_scroll">
                            ''' + set_data["main_css_table_scroll"] + '''
                        </select></span>
                        <h4>''' + await get_lang("table_transparent") + '''</h4>
                        ''' + set_data_main["main_css_table_transparent"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_table_transparent">
                            ''' + set_data["main_css_table_transparent"] + '''
                        </select></span>
                        <h3>''' + await get_lang("list_view_change") + '''</h3>
                        ''' + set_data_main["main_css_list_view_change"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_list_view_change">
                            ''' + set_data["main_css_list_view_change"] + '''
                        </select></span>
                        <h3>''' + await get_lang("view_joke") + '''</h3>
                        ''' + set_data_main["main_css_view_joke"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_view_joke">
                            ''' + set_data["main_css_view_joke"] + '''
                        </select></span>
                        <h3>''' + await get_lang("math_scroll") + '''</h3>
                        ''' + set_data_main["main_css_math_scroll"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_math_scroll">
                            ''' + set_data["main_css_math_scroll"] + '''
                        </select></span>
                        <h3>''' + await get_lang("view_history") + '''</h3>
                        ''' + set_data_main["main_css_view_history"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_view_history">
                            ''' + set_data["main_css_view_history"] + '''
                        </select></span>
                        <h3>''' + await get_lang("font_size") + '''</h3>
                        ''' + set_data_main["main_css_font_size"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_font_size">
                            ''' + set_data["main_css_font_size"] + '''
                        </select></span>
                        <h2>''' + await get_lang("edit") + '''</h2>
                        <h3>''' + await get_lang("monaco_editor") + '''</h3>
                        ''' + set_data_main["main_css_monaco"] + '''
                        <span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="main_css_monaco">
                            ''' + set_data["main_css_monaco"] + '''
                        </select></span>
                        <hr class="main_hr">
                        <button class="__ON_BUTTON__" type="submit">''' + await get_lang('save') + '''</button>
                    </form>
                '''),
                '(' + await get_lang('beta') + ')',
                [['setting', await get_lang('return')]]
            )
