from .tool.func import *

def user_setting_skin_set_main_set_list(conn):
    set_list = {
        'main_css_strike' : [
            ['default', get_lang(conn, 'default')],
            ['normal', get_lang(conn, 'off')],
            ['change', get_lang(conn, 'change_to_normal')],
            ['delete', get_lang(conn, 'delete')]
        ], 'main_css_bold' : [
            ['default', get_lang(conn, 'default')],
            ['normal', get_lang(conn, 'off')],
            ['change', get_lang(conn, 'change_to_normal')],
            ['delete', get_lang(conn, 'delete')]
        ], 'main_css_include_link' : [
            ['default', get_lang(conn, 'default')],
            ['normal', get_lang(conn, 'off')],
            ['use', get_lang(conn, 'use')]
        ], 'main_css_category_set' : [
            ['default', get_lang(conn, 'default')],
            ['bottom', get_lang(conn, 'bottom')],
            ['top', get_lang(conn, 'top')]
        ], 'main_css_footnote_set' : [
            ['default', get_lang(conn, 'default')],
            ['normal', get_lang(conn, 'normal')],
            ['spread', get_lang(conn, 'spread')],
            ['popup', get_lang(conn, 'popup') + ' (' + get_lang(conn, 'not_working') + ')'],
            ['popover', get_lang(conn, 'popover')]
        ], 'main_css_image_set' : [
            ['default', get_lang(conn, 'default')],
            ['normal', get_lang(conn, 'normal')],
            ['click', get_lang(conn, 'change_to_link')],
            ['new_click', get_lang(conn, 'click_load')]
        ], 'main_css_toc_set' : [
            ['default', get_lang(conn, 'default')],
            ['normal', get_lang(conn, 'normal')],
            ['off', get_lang(conn, 'all_off')],
            ['half_off', get_lang(conn, 'in_content')]
        ], 'main_css_monaco' : [
            ['default', get_lang(conn, 'default')],
            ['normal', get_lang(conn, 'off')],
            ['use', get_lang(conn, 'use')]
        ], 'main_css_exter_link' : [
            ['default', get_lang(conn, 'default')],
            ['blank', get_lang(conn, 'normal')],
            ['self', get_lang(conn, 'self_tab')]
        ], 'main_css_link_delimiter' : [
            ['default', get_lang(conn, 'default')],
            ['normal', get_lang(conn, 'off')],
            ['use', get_lang(conn, 'use')]
        ], 'main_css_darkmode' : [
            ['default', get_lang(conn, 'default')],
            ['0', get_lang(conn, 'off')],
            ['1', get_lang(conn, 'use')]
        ], 'main_css_footnote_number' : [
            ['default', get_lang(conn, 'default')],
            ['all', get_lang(conn, 'all')],
            ['only_number', get_lang(conn, 'only_number')]
        ], 'main_css_view_real_footnote_num' : [
            ['default', get_lang(conn, 'default')],
            ['off', get_lang(conn, 'off')],
            ['on', get_lang(conn, 'use')]
        ], 'main_css_table_scroll' : [
            ['default', get_lang(conn, 'default')],
            ['off', get_lang(conn, 'off')],
            ['on', get_lang(conn, 'use')]
        ], 'main_css_category_change_title' : [
            ['default', get_lang(conn, 'default')],
            ['off', get_lang(conn, 'off')],
            ['on', get_lang(conn, 'use')]
        ], 'main_css_list_view_change' : [
            ['default', get_lang(conn, 'default')],
            ['off', get_lang(conn, 'off')],
            ['on', get_lang(conn, 'use')]
        ], 'main_css_view_joke' : [
            ['default', get_lang(conn, 'default')],
            ['on', get_lang(conn, 'use')],
            ['off', get_lang(conn, 'off')]
        ], 'main_css_math_scroll' : [
            ['default', get_lang(conn, 'default')],
            ['off', get_lang(conn, 'off')],
            ['on', get_lang(conn, 'use')]
        ], 'main_css_view_history' : [
            ['default', get_lang(conn, 'default')],
            ['off', get_lang(conn, 'off')],
            ['on', get_lang(conn, 'use')]
        ], 'main_css_table_transparent' : [
            ['default', get_lang(conn, 'default')],
            ['off', get_lang(conn, 'off')],
            ['on', get_lang(conn, 'use')]
        ], 'main_css_font_size' : [
            ['default', get_lang(conn, 'default')],
            ['10', '10'],
            ['12', '12'],
            ['14', '14'],
            ['16', '16'],
            ['18', '18'],
            ['20', '20'],
            ['22', '22'],
        ]
    }

    return set_list

async def user_setting_skin_set_main():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if (await ban_check(ip))[0] == 1:
            return await re_error(conn, 0)
            
        set_list = user_setting_skin_set_main_set_list(conn)
        use_cookie = ['main_css_darkmode']

        if flask.request.method == 'POST':
            html_data = flask.make_response(redirect(conn, '/change/skin_set/main'))

            for for_b in set_list:
                if for_b in use_cookie:
                    html_data.set_cookie(for_b, flask.request.form.get(for_b, set_list[for_b][0][0]))
                elif ip_or_user(ip) == 0:
                    curs.execute(db_change('select data from user_set where name = ? and id = ?'), [for_b, ip])
                    if curs.fetchall():
                        curs.execute(db_change("update user_set set data = ? where name = ? and id = ?"), [
                            flask.request.form.get(for_b, set_list[for_b][0][0]),
                            for_b,
                            ip
                        ])
                    else:
                        curs.execute(db_change('insert into user_set (name, id, data) values (?, ?, ?)'), [
                            for_b, 
                            ip,
                            flask.request.form.get(for_b, set_list[for_b][0][0])
                        ])
                else:
                    flask.session[for_b] = flask.request.form.get(for_b, set_list[for_b][0][0])

            return html_data
        else:
            set_data = {}
            for for_b in set_list:
                set_data[for_b] = ''
                if for_b in use_cookie:
                    get_data = flask.request.cookies.get(for_b, '')
                elif ip_or_user(ip) == 0:
                    curs.execute(db_change('select data from user_set where name = ? and id = ?'), [for_b, ip])
                    db_data = curs.fetchall()
                    get_data = db_data[0][0] if db_data else ''
                else:
                    get_data = flask.session[for_b] if for_b in flask.session else ''

                for for_a in set_list[for_b]:
                    if get_data == for_a[0]:
                        set_data[for_b] = '<option value="' + for_a[0] + '">' + for_a[1] + '</option>' + set_data[for_b]
                    else:
                        set_data[for_b] += '<option value="' + for_a[0] + '">' + for_a[1] + '</option>'

            set_data_main = {}
            for for_b in set_list:
                curs.execute(db_change('select data from other where name = ?'), [for_b])
                db_data = curs.fetchall()
                server_default = db_data[0][0] if db_data else 'default'
                set_data_main[for_b] = get_lang(conn, 'default') + ' : ' + ''.join([for_a[1] for for_a in set_list[for_b] if for_a[0] == server_default]) + '<hr class="main_hr">'

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'main_skin_set'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
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
                menu = [['change', get_lang(conn, 'user_setting')], ['change/skin_set', get_lang(conn, 'skin_set')], ['setting/skin_set', get_lang(conn, 'main_skin_set_default')]]
            ))