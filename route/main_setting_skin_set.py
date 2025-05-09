from .tool.func import *

from .user_setting_skin_set_main import user_setting_skin_set_main_set_list

async def main_setting_skin_set():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 0)
            
        set_list = user_setting_skin_set_main_set_list(conn)

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

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'main_skin_set_default'), await wiki_set(), await wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'beta') + ')', 0])],
                data = render_simple_set(conn, '''
                    <form method="post">
                        <h2>''' + get_lang(conn, "render") + '''</h2>
                        <h3>''' + get_lang(conn, "strike") + '''</h3>
                        ''' + set_data_main["main_css_strike"] + '''
                        <select name="main_css_strike">
                            ''' + set_data["main_css_strike"] + '''
                        </select>
                        <h3>''' + get_lang(conn, "bold") + '''</h3>
                        ''' + set_data_main["main_css_bold"] + '''
                        <select name="main_css_bold">
                            ''' + set_data["main_css_bold"] + '''
                        </select>
                        <h3>''' + get_lang(conn, "category") + '''</h3>
                        <h4>''' + get_lang(conn, "position") + '''</h4>
                        ''' + set_data_main["main_css_category_set"] + '''
                        <select name="main_css_category_set">
                            ''' + set_data["main_css_category_set"] + '''
                        </select>
                        <h4>''' + get_lang(conn, "category_change_title") + '''</h4>
                        ''' + set_data_main["main_css_category_change_title"] + '''
                        <select name="main_css_category_change_title">
                            ''' + set_data["main_css_category_change_title"] + '''
                        </select>
                        <h3>''' + get_lang(conn, "footnote") + ''' (''' + get_lang(conn, 'beta') + ''')</h3>
                        <h4>''' + get_lang(conn, "footnote_render") + '''</h4>
                        ''' + set_data_main["main_css_footnote_set"] + '''
                        <select name="main_css_footnote_set">
                            ''' + set_data["main_css_footnote_set"] + '''
                        </select>
                        <h4>''' + get_lang(conn, "footnote_number") + '''</h4>
                        ''' + set_data_main["main_css_footnote_number"] + '''
                        <select name="main_css_footnote_number">
                            ''' + set_data["main_css_footnote_number"] + '''
                        </select>
                        <h4>''' + get_lang(conn, "footnote_real_num_view") + '''</h4>
                        ''' + set_data_main["main_css_view_real_footnote_num"] + '''
                        <select name="main_css_view_real_footnote_num">
                            ''' + set_data["main_css_view_real_footnote_num"] + '''
                        </select>
                        <h3>''' + get_lang(conn, "include_link") + '''</h3>
                        ''' + set_data_main["main_css_include_link"] + '''
                        <select name="main_css_include_link">
                            ''' + set_data["main_css_include_link"] + '''
                        </select>
                        <h3>''' + get_lang(conn, "image") + ''' (''' + get_lang(conn, 'beta') + ''')</h3>
                        ''' + set_data_main["main_css_image_set"] + '''
                        <select name="main_css_image_set">
                            ''' + set_data["main_css_image_set"] + '''
                        </select>
                        <h3>''' + get_lang(conn, "toc") + '''</h3>
                        ''' + set_data_main["main_css_toc_set"] + '''
                        <select name="main_css_toc_set">
                            ''' + set_data["main_css_toc_set"] + '''
                        </select>
                        <h3>''' + get_lang(conn, "exter_link") + '''</h3>
                        ''' + set_data_main["main_css_exter_link"] + '''
                        <select name="main_css_exter_link">
                            ''' + set_data["main_css_exter_link"] + '''
                        </select>
                        <h3>''' + get_lang(conn, "link_delimiter") + '''</h3>
                        ''' + set_data_main["main_css_link_delimiter"] + '''
                        <select name="main_css_link_delimiter">
                            ''' + set_data["main_css_link_delimiter"] + '''
                        </select>
                        <h3>''' + get_lang(conn, "force_darkmode") + '''</h3>
                        ''' + set_data_main["main_css_darkmode"] + '''
                        <select name="main_css_darkmode">
                            ''' + set_data["main_css_darkmode"] + '''
                        </select>
                        <h3>''' + get_lang(conn, "table") + '''</h3>
                        <h4>''' + get_lang(conn, "table_scroll") + '''</h4>
                        ''' + set_data_main["main_css_table_scroll"] + '''
                        <select name="main_css_table_scroll">
                            ''' + set_data["main_css_table_scroll"] + '''
                        </select>
                        <h4>''' + get_lang(conn, "table_transparent") + '''</h4>
                        ''' + set_data_main["main_css_table_transparent"] + '''
                        <select name="main_css_table_transparent">
                            ''' + set_data["main_css_table_transparent"] + '''
                        </select>
                        <h3>''' + get_lang(conn, "list_view_change") + '''</h3>
                        ''' + set_data_main["main_css_list_view_change"] + '''
                        <select name="main_css_list_view_change">
                            ''' + set_data["main_css_list_view_change"] + '''
                        </select>
                        <h3>''' + get_lang(conn, "view_joke") + '''</h3>
                        ''' + set_data_main["main_css_view_joke"] + '''
                        <select name="main_css_view_joke">
                            ''' + set_data["main_css_view_joke"] + '''
                        </select>
                        <h3>''' + get_lang(conn, "math_scroll") + '''</h3>
                        ''' + set_data_main["main_css_math_scroll"] + '''
                        <select name="main_css_math_scroll">
                            ''' + set_data["main_css_math_scroll"] + '''
                        </select>
                        <h3>''' + get_lang(conn, "view_history") + '''</h3>
                        ''' + set_data_main["main_css_view_history"] + '''
                        <select name="main_css_view_history">
                            ''' + set_data["main_css_view_history"] + '''
                        </select>
                        <h3>''' + get_lang(conn, "font_size") + '''</h3>
                        ''' + set_data_main["main_css_font_size"] + '''
                        <select name="main_css_font_size">
                            ''' + set_data["main_css_font_size"] + '''
                        </select>
                        <h2>''' + get_lang(conn, "edit") + '''</h2>
                        <h3>''' + get_lang(conn, "monaco_editor") + '''</h3>
                        ''' + set_data_main["main_css_monaco"] + '''
                        <select name="main_css_monaco">
                            ''' + set_data["main_css_monaco"] + '''
                        </select>
                        <hr class="main_hr">
                        <button type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                '''),
                menu = [['setting', get_lang(conn, 'return')]]
            ))