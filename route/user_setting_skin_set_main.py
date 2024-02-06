from .tool.func import *

def user_setting_skin_set_main_set_list():
    set_list = {
        'main_css_strike' : [
            ['default', load_lang('default')],
            ['normal', load_lang('off')],
            ['change', load_lang('change_to_normal')],
            ['delete', load_lang('delete')]
        ], 'main_css_bold' : [
            ['default', load_lang('default')],
            ['normal', load_lang('off')],
            ['change', load_lang('change_to_normal')],
            ['delete', load_lang('delete')]
        ], 'main_css_include_link' : [
            ['default', load_lang('default')],
            ['normal', load_lang('off')],
            ['use', load_lang('use')]
        ], 'main_css_image_paste' : [
            ['default', load_lang('default')],
            ['normal', load_lang('off')],
            ['use', load_lang('use')]
        ], 'main_css_category_set' : [
            ['default', load_lang('default')],
            ['bottom', load_lang('bottom')],
            ['top', load_lang('top')]
        ], 'main_css_footnote_set' : [
            ['default', load_lang('default')],
            ['normal', load_lang('normal')],
            ['spread', load_lang('spread')],
            ['popup', load_lang('popup') + ' (' + load_lang('not_working') + ')'],
            ['popover', load_lang('popover')]
        ], 'main_css_image_set' : [
            ['default', load_lang('default')],
            ['normal', load_lang('normal')],
            ['click', load_lang('change_to_link')],
            ['new_click', load_lang('click_load')]
        ], 'main_css_toc_set' : [
            ['default', load_lang('default')],
            ['normal', load_lang('normal')],
            ['off', load_lang('all_off')],
            ['half_off', load_lang('in_content')]
        ], 'main_css_monaco' : [
            ['default', load_lang('default')],
            ['normal', load_lang('off')],
            ['use', load_lang('use')]
        ], 'main_css_exter_link' : [
            ['default', load_lang('default')],
            ['blank', load_lang('normal')],
            ['self', load_lang('self_tab')]
        ], 'main_css_link_delimiter' : [
            ['default', load_lang('default')],
            ['normal', load_lang('off')],
            ['use', load_lang('use')]
        ], 'main_css_darkmode' : [
            ['default', load_lang('default')],
            ['0', load_lang('off')],
            ['1', load_lang('use')]
        ], 'main_css_footnote_number' : [
            ['default', load_lang('default')],
            ['all', load_lang('all')],
            ['only_number', load_lang('only_number')]
        ], 'main_css_view_real_footnote_num' : [
            ['default', load_lang('default')],
            ['off', load_lang('off')],
            ['on', load_lang('use')]
        ], 'main_css_table_scroll' : [
            ['default', load_lang('default')],
            ['off', load_lang('off')],
            ['on', load_lang('use')]
        ], 'main_css_category_change_title' : [
            ['default', load_lang('default')],
            ['off', load_lang('off')],
            ['on', load_lang('use')]
        ], 'main_css_list_view_change' : [
            ['default', load_lang('default')],
            ['off', load_lang('off')],
            ['on', load_lang('use')]
        ], 'main_css_view_joke' : [
            ['default', load_lang('default')],
            ['on', load_lang('use')],
            ['off', load_lang('off')]
        ], 'main_css_math_scroll' : [
            ['default', load_lang('default')],
            ['off', load_lang('off')],
            ['on', load_lang('use')]
        ], 'main_css_view_history' : [
            ['default', load_lang('default')],
            ['off', load_lang('off')],
            ['on', load_lang('use')]
        ], 'main_css_table_transparent' : [
            ['default', load_lang('default')],
            ['off', load_lang('off')],
            ['on', load_lang('use')]
        ], 'main_css_font_size' : [
            ['default', load_lang('default')],
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

def user_setting_skin_set_main():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if ban_check(ip) == 1:
            return re_error('/ban')
            
        set_list = user_setting_skin_set_main_set_list()
        use_cookie = ['main_css_image_paste', 'main_css_darkmode']

        if flask.request.method == 'POST':
            html_data = flask.make_response(redirect('/change/skin_set/main'))

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

            conn.commit()

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
                set_data_main[for_b] = load_lang('default') + ' : ' + ''.join([for_a[1] for for_a in set_list[for_b] if for_a[0] == server_default]) + '<hr class="main_hr">'

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('main_skin_set'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = render_simple_set('''
                    <form method="post">
                        <h2>''' + load_lang("render") + '''</h2>
                        <h3>''' + load_lang("strike") + '''</h3>
                        ''' + set_data_main["main_css_strike"] + '''
                        <select name="main_css_strike">
                            ''' + set_data["main_css_strike"] + '''
                        </select>
                        <h3>''' + load_lang("bold") + '''</h3>
                        ''' + set_data_main["main_css_bold"] + '''
                        <select name="main_css_bold">
                            ''' + set_data["main_css_bold"] + '''
                        </select>
                        <h3>''' + load_lang("category") + '''</h3>
                        <h4>''' + load_lang("position") + '''</h4>
                        ''' + set_data_main["main_css_category_set"] + '''
                        <select name="main_css_category_set">
                            ''' + set_data["main_css_category_set"] + '''
                        </select>
                        <h4>''' + load_lang("category_change_title") + '''</h4>
                        ''' + set_data_main["main_css_category_change_title"] + '''
                        <select name="main_css_category_change_title">
                            ''' + set_data["main_css_category_change_title"] + '''
                        </select>
                        <h3>''' + load_lang("footnote") + ''' (''' + load_lang('beta') + ''')</h3>
                        <h4>''' + load_lang("footnote_render") + '''</h4>
                        ''' + set_data_main["main_css_footnote_set"] + '''
                        <select name="main_css_footnote_set">
                            ''' + set_data["main_css_footnote_set"] + '''
                        </select>
                        <h4>''' + load_lang("footnote_number") + '''</h4>
                        ''' + set_data_main["main_css_footnote_number"] + '''
                        <select name="main_css_footnote_number">
                            ''' + set_data["main_css_footnote_number"] + '''
                        </select>
                        <h4>''' + load_lang("footnote_real_num_view") + '''</h4>
                        ''' + set_data_main["main_css_view_real_footnote_num"] + '''
                        <select name="main_css_view_real_footnote_num">
                            ''' + set_data["main_css_view_real_footnote_num"] + '''
                        </select>
                        <h3>''' + load_lang("include_link") + '''</h3>
                        ''' + set_data_main["main_css_include_link"] + '''
                        <select name="main_css_include_link">
                            ''' + set_data["main_css_include_link"] + '''
                        </select>
                        <h3>''' + load_lang("image") + ''' (''' + load_lang('beta') + ''')</h3>
                        ''' + set_data_main["main_css_image_set"] + '''
                        <select name="main_css_image_set">
                            ''' + set_data["main_css_image_set"] + '''
                        </select>
                        <h3>''' + load_lang("toc") + '''</h3>
                        ''' + set_data_main["main_css_toc_set"] + '''
                        <select name="main_css_toc_set">
                            ''' + set_data["main_css_toc_set"] + '''
                        </select>
                        <h3>''' + load_lang("exter_link") + '''</h3>
                        ''' + set_data_main["main_css_exter_link"] + '''
                        <select name="main_css_exter_link">
                            ''' + set_data["main_css_exter_link"] + '''
                        </select>
                        <h3>''' + load_lang("link_delimiter") + '''</h3>
                        ''' + set_data_main["main_css_link_delimiter"] + '''
                        <select name="main_css_link_delimiter">
                            ''' + set_data["main_css_link_delimiter"] + '''
                        </select>
                        <h3>''' + load_lang("force_darkmode") + '''</h3>
                        ''' + set_data_main["main_css_darkmode"] + '''
                        <select name="main_css_darkmode">
                            ''' + set_data["main_css_darkmode"] + '''
                        </select>
                        <h3>''' + load_lang("table") + '''</h3>
                        <h4>''' + load_lang("table_scroll") + '''</h4>
                        ''' + set_data_main["main_css_table_scroll"] + '''
                        <select name="main_css_table_scroll">
                            ''' + set_data["main_css_table_scroll"] + '''
                        </select>
                        <h4>''' + load_lang("table_transparent") + '''</h4>
                        ''' + set_data_main["main_css_table_transparent"] + '''
                        <select name="main_css_table_transparent">
                            ''' + set_data["main_css_table_transparent"] + '''
                        </select>
                        <h3>''' + load_lang("list_view_change") + '''</h3>
                        ''' + set_data_main["main_css_list_view_change"] + '''
                        <select name="main_css_list_view_change">
                            ''' + set_data["main_css_list_view_change"] + '''
                        </select>
                        <h3>''' + load_lang("view_joke") + '''</h3>
                        ''' + set_data_main["main_css_view_joke"] + '''
                        <select name="main_css_view_joke">
                            ''' + set_data["main_css_view_joke"] + '''
                        </select>
                        <h3>''' + load_lang("math_scroll") + '''</h3>
                        ''' + set_data_main["main_css_math_scroll"] + '''
                        <select name="main_css_math_scroll">
                            ''' + set_data["main_css_math_scroll"] + '''
                        </select>
                        <h3>''' + load_lang("view_history") + '''</h3>
                        ''' + set_data_main["main_css_view_history"] + '''
                        <select name="main_css_view_history">
                            ''' + set_data["main_css_view_history"] + '''
                        </select>
                        <h3>''' + load_lang("font_size") + '''</h3>
                        ''' + set_data_main["main_css_font_size"] + '''
                        <select name="main_css_font_size">
                            ''' + set_data["main_css_font_size"] + '''
                        </select>
                        <h2>''' + load_lang("edit") + '''</h2>
                        <h3>''' + load_lang("image_paste") + '''</h3>
                        <sup>''' + load_lang('only_korean') + '''</sup> <sup>''' + load_lang('unavailable_in_monaco') + '''</sup>
                        <hr class="main_hr">
                        ''' + set_data_main["main_css_image_paste"] + '''
                        <select name="main_css_image_paste">
                            ''' + set_data["main_css_image_paste"] + '''
                        </select>
                        <h3>''' + load_lang("monaco_editor") + '''</h3>
                        ''' + set_data_main["main_css_monaco"] + '''
                        <select name="main_css_monaco">
                            ''' + set_data["main_css_monaco"] + '''
                        </select>
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('save') + '''</button>
                    </form>
                '''),
                menu = [['change', load_lang('user_setting')], ['change/skin_set', load_lang('skin_set')], ['setting/skin_set', load_lang('main_skin_set_default')]]
            ))