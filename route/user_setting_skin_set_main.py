from .tool.func import *

async def user_setting_skin_set_main_set_list():
    set_list = {
        'main_css_strike' : [
            ['default', await get_lang('default')],
            ['normal', await get_lang('off')],
            ['change', await get_lang('change_to_normal')],
            ['delete', await get_lang('delete')]
        ], 'main_css_bold' : [
            ['default', await get_lang('default')],
            ['normal', await get_lang('off')],
            ['change', await get_lang('change_to_normal')],
            ['delete', await get_lang('delete')]
        ], 'main_css_include_link' : [
            ['default', await get_lang('default')],
            ['normal', await get_lang('off')],
            ['use', await get_lang('use')]
        ], 'main_css_category_set' : [
            ['default', await get_lang('default')],
            ['bottom', await get_lang('bottom')],
            ['top', await get_lang('top')]
        ], 'main_css_footnote_set' : [
            ['default', await get_lang('default')],
            ['normal', await get_lang('normal')],
            ['spread', await get_lang('spread')],
            ['popup', await get_lang('popup') + ' (' + await get_lang('not_working') + ')'],
            ['popover', await get_lang('popover')]
        ], 'main_css_image_set' : [
            ['default', await get_lang('default')],
            ['normal', await get_lang('normal')],
            ['click', await get_lang('change_to_link')],
            ['new_click', await get_lang('click_load')]
        ], 'main_css_toc_set' : [
            ['default', await get_lang('default')],
            ['normal', await get_lang('normal')],
            ['off', await get_lang('all_off')],
            ['half_off', await get_lang('in_content')]
        ], 'main_css_monaco' : [
            ['default', await get_lang('default')],
            ['normal', await get_lang('off')],
            ['use', await get_lang('use')]
        ], 'main_css_exter_link' : [
            ['default', await get_lang('default')],
            ['blank', await get_lang('normal')],
            ['self', await get_lang('self_tab')]
        ], 'main_css_link_delimiter' : [
            ['default', await get_lang('default')],
            ['normal', await get_lang('off')],
            ['use', await get_lang('use')]
        ], 'main_css_darkmode' : [
            ['default', await get_lang('default')],
            ['0', await get_lang('off')],
            ['1', await get_lang('use')]
        ], 'main_css_footnote_number' : [
            ['default', await get_lang('default')],
            ['all', await get_lang('all')],
            ['only_number', await get_lang('only_number')]
        ], 'main_css_view_real_footnote_num' : [
            ['default', await get_lang('default')],
            ['off', await get_lang('off')],
            ['on', await get_lang('use')]
        ], 'main_css_table_scroll' : [
            ['default', await get_lang('default')],
            ['off', await get_lang('off')],
            ['on', await get_lang('use')]
        ], 'main_css_category_change_title' : [
            ['default', await get_lang('default')],
            ['off', await get_lang('off')],
            ['on', await get_lang('use')]
        ], 'main_css_list_view_change' : [
            ['default', await get_lang('default')],
            ['off', await get_lang('off')],
            ['on', await get_lang('use')]
        ], 'main_css_view_joke' : [
            ['default', await get_lang('default')],
            ['on', await get_lang('use')],
            ['off', await get_lang('off')]
        ], 'main_css_math_scroll' : [
            ['default', await get_lang('default')],
            ['off', await get_lang('off')],
            ['on', await get_lang('use')]
        ], 'main_css_view_history' : [
            ['default', await get_lang('default')],
            ['off', await get_lang('off')],
            ['on', await get_lang('use')]
        ], 'main_css_table_transparent' : [
            ['default', await get_lang('default')],
            ['off', await get_lang('off')],
            ['on', await get_lang('use')]
        ], 'main_css_font_size' : [
            ['default', await get_lang('default')],
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
            
        set_list = await user_setting_skin_set_main_set_list()
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
                set_data_main[for_b] = await get_lang('default') + ' : ' + ''.join([for_a[1] for for_a in set_list[for_b] if for_a[0] == server_default]) + '<hr class="main_hr">'

            return await render_template(
                await get_lang('main_skin_set'),
                await render_simple_set('''
                    <form method="post">
                        <h2>''' + await get_lang("render") + '''</h2>
                        <h3>''' + await get_lang("strike") + '''</h3>
                        ''' + set_data_main["main_css_strike"] + '''
                        <select name="main_css_strike">
                            ''' + set_data["main_css_strike"] + '''
                        </select>
                        <h3>''' + await get_lang("bold") + '''</h3>
                        ''' + set_data_main["main_css_bold"] + '''
                        <select name="main_css_bold">
                            ''' + set_data["main_css_bold"] + '''
                        </select>
                        <h3>''' + await get_lang("category") + '''</h3>
                        <h4>''' + await get_lang("position") + '''</h4>
                        ''' + set_data_main["main_css_category_set"] + '''
                        <select name="main_css_category_set">
                            ''' + set_data["main_css_category_set"] + '''
                        </select>
                        <h4>''' + await get_lang("category_change_title") + '''</h4>
                        ''' + set_data_main["main_css_category_change_title"] + '''
                        <select name="main_css_category_change_title">
                            ''' + set_data["main_css_category_change_title"] + '''
                        </select>
                        <h3>''' + await get_lang("footnote") + ''' (''' + await get_lang('beta') + ''')</h3>
                        <h4>''' + await get_lang("footnote_render") + '''</h4>
                        ''' + set_data_main["main_css_footnote_set"] + '''
                        <select name="main_css_footnote_set">
                            ''' + set_data["main_css_footnote_set"] + '''
                        </select>
                        <h4>''' + await get_lang("footnote_number") + '''</h4>
                        ''' + set_data_main["main_css_footnote_number"] + '''
                        <select name="main_css_footnote_number">
                            ''' + set_data["main_css_footnote_number"] + '''
                        </select>
                        <h4>''' + await get_lang("footnote_real_num_view") + '''</h4>
                        ''' + set_data_main["main_css_view_real_footnote_num"] + '''
                        <select name="main_css_view_real_footnote_num">
                            ''' + set_data["main_css_view_real_footnote_num"] + '''
                        </select>
                        <h3>''' + await get_lang("include_link") + '''</h3>
                        ''' + set_data_main["main_css_include_link"] + '''
                        <select name="main_css_include_link">
                            ''' + set_data["main_css_include_link"] + '''
                        </select>
                        <h3>''' + await get_lang("image") + ''' (''' + await get_lang('beta') + ''')</h3>
                        ''' + set_data_main["main_css_image_set"] + '''
                        <select name="main_css_image_set">
                            ''' + set_data["main_css_image_set"] + '''
                        </select>
                        <h3>''' + await get_lang("toc") + '''</h3>
                        ''' + set_data_main["main_css_toc_set"] + '''
                        <select name="main_css_toc_set">
                            ''' + set_data["main_css_toc_set"] + '''
                        </select>
                        <h3>''' + await get_lang("exter_link") + '''</h3>
                        ''' + set_data_main["main_css_exter_link"] + '''
                        <select name="main_css_exter_link">
                            ''' + set_data["main_css_exter_link"] + '''
                        </select>
                        <h3>''' + await get_lang("link_delimiter") + '''</h3>
                        ''' + set_data_main["main_css_link_delimiter"] + '''
                        <select name="main_css_link_delimiter">
                            ''' + set_data["main_css_link_delimiter"] + '''
                        </select>
                        <h3>''' + await get_lang("force_darkmode") + '''</h3>
                        ''' + set_data_main["main_css_darkmode"] + '''
                        <select name="main_css_darkmode">
                            ''' + set_data["main_css_darkmode"] + '''
                        </select>
                        <h3>''' + await get_lang("table") + '''</h3>
                        <h4>''' + await get_lang("table_scroll") + '''</h4>
                        ''' + set_data_main["main_css_table_scroll"] + '''
                        <select name="main_css_table_scroll">
                            ''' + set_data["main_css_table_scroll"] + '''
                        </select>
                        <h4>''' + await get_lang("table_transparent") + '''</h4>
                        ''' + set_data_main["main_css_table_transparent"] + '''
                        <select name="main_css_table_transparent">
                            ''' + set_data["main_css_table_transparent"] + '''
                        </select>
                        <h3>''' + await get_lang("list_view_change") + '''</h3>
                        ''' + set_data_main["main_css_list_view_change"] + '''
                        <select name="main_css_list_view_change">
                            ''' + set_data["main_css_list_view_change"] + '''
                        </select>
                        <h3>''' + await get_lang("view_joke") + '''</h3>
                        ''' + set_data_main["main_css_view_joke"] + '''
                        <select name="main_css_view_joke">
                            ''' + set_data["main_css_view_joke"] + '''
                        </select>
                        <h3>''' + await get_lang("math_scroll") + '''</h3>
                        ''' + set_data_main["main_css_math_scroll"] + '''
                        <select name="main_css_math_scroll">
                            ''' + set_data["main_css_math_scroll"] + '''
                        </select>
                        <h3>''' + await get_lang("view_history") + '''</h3>
                        ''' + set_data_main["main_css_view_history"] + '''
                        <select name="main_css_view_history">
                            ''' + set_data["main_css_view_history"] + '''
                        </select>
                        <h3>''' + await get_lang("font_size") + '''</h3>
                        ''' + set_data_main["main_css_font_size"] + '''
                        <select name="main_css_font_size">
                            ''' + set_data["main_css_font_size"] + '''
                        </select>
                        <h2>''' + await get_lang("edit") + '''</h2>
                        <h3>''' + await get_lang("monaco_editor") + '''</h3>
                        ''' + set_data_main["main_css_monaco"] + '''
                        <select name="main_css_monaco">
                            ''' + set_data["main_css_monaco"] + '''
                        </select>
                        <hr class="main_hr">
                        <button class="__ON_BUTTON__" type="submit">''' + await get_lang('save') + '''</button>
                    </form>
                '''),
                0,
                [['change', await get_lang('user_setting')], ['change/skin_set', await get_lang('skin_set')], ['setting/skin_set', await get_lang('main_skin_set_default')]]
            )
