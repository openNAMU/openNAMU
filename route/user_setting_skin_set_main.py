from .tool.func import *

def user_setting_skin_set_main():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if ban_check(ip) == 1:
            return re_error('/ban')
            
        set_list = {
            'main_css_strike' : [
                ['normal', load_lang('default')],
                ['change', load_lang('change_to_normal')],
                ['delete', load_lang('delete')]
            ], 'main_css_bold' : [
                ['normal', load_lang('default')],
                ['change', load_lang('change_to_normal')],
                ['delete', load_lang('delete')]
            ], 'main_css_include_link' : [
                ['normal', load_lang('default')],
                ['use', load_lang('use')]
            ], 'main_css_image_paste' : [
                ['normal', load_lang('default')],
                ['use', load_lang('use')]
            ], 'main_css_category_set' : [
                ['bottom', load_lang('default')],
                ['top', load_lang('top')]
            ], 'main_css_footnote_set' : [
                ['normal', load_lang('default')],
                ['spread', load_lang('spread')]
            ], 'main_css_image_set' : [
                ['normal', load_lang('default')],
                ['click', load_lang('change_to_link')],
                ['new_click', load_lang('click_load')]
            ], 'main_css_toc_set' : [
                ['normal', load_lang('default')],
                ['off', load_lang('all_off')],
                ['half_off', load_lang('in_content')]
            ], 'main_css_monaco' : [
                ['normal', load_lang('default')],
                ['use', load_lang('use')]
            ], 'main_css_exter_link' : [
                ['blank', load_lang('default')],
                ['self', load_lang('self_tab')]
            ], 'main_css_link_delimiter' : [
                ['normal', load_lang('default')],
                ['use', load_lang('use')]
            ], 'main_css_darkmode' : [
                ['0', load_lang('default')],
                ['1', load_lang('use')]
            ], 'main_css_font_size' : [
                ['']
            ]
        }
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

                if set_list[for_b][0] == ['']:
                    set_data[for_b] = get_data
                else:
                    for for_a in set_list[for_b]:
                        if get_data == for_a[0]:
                            set_data[for_b] = '<option value="' + for_a[0] + '">' + for_a[1] + '</option>' + set_data[for_b]
                        else:
                            set_data[for_b] += '<option value="' + for_a[0] + '">' + for_a[1] + '</option>'

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('main_skin_set'), wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('beta') + ')', 0])],
                data = render_simple_set('''
                    <form method="post">
                        <h2>''' + load_lang("render") + '''</h2>
                        <h3>''' + load_lang("strike") + '''</h3>
                        <select name="main_css_strike">
                            ''' + set_data["main_css_strike"] + '''
                        </select>
                        <h3>''' + load_lang("bold") + '''</h3>
                        <select name="main_css_bold">
                            ''' + set_data["main_css_bold"] + '''
                        </select>
                        <h3>''' + load_lang("category") + '''</h3>
                        <select name="main_css_category_set">
                            ''' + set_data["main_css_category_set"] + '''
                        </select>
                        <h3>''' + load_lang("footnote") + '''</h3>
                        <select name="main_css_footnote_set">
                            ''' + set_data["main_css_footnote_set"] + '''
                        </select>
                        <h3>''' + load_lang("include_link") + '''</h3>
                        <select name="main_css_include_link">
                            ''' + set_data["main_css_include_link"] + '''
                        </select>
                        <h3>''' + load_lang("image") + ''' (''' + load_lang("not_working") + ''')</h3>
                        <select name="main_css_image_set">
                            ''' + set_data["main_css_image_set"] + '''
                        </select>
                        <h3>''' + load_lang("toc") + '''</h3>
                        <select name="main_css_toc_set">
                            ''' + set_data["main_css_toc_set"] + '''
                        </select>
                        <h3>''' + load_lang("exter_link") + '''</h3>
                        <select name="main_css_exter_link">
                            ''' + set_data["main_css_exter_link"] + '''
                        </select>
                        <h3>''' + load_lang("link_delimiter") + '''</h3>
                        <select name="main_css_link_delimiter">
                            ''' + set_data["main_css_link_delimiter"] + '''
                        </select>
                        <h3>''' + load_lang("force_darkmode") + '''</h3>
                        <select name="main_css_darkmode">
                            ''' + set_data["main_css_darkmode"] + '''
                        </select>
                        <h3>''' + load_lang("font_size") + '''</h3>
                        (EX : 11) (''' + load_lang('off') + ''' : ''' + load_lang('empty') + ''')
                        <hr class="main_hr">
                        <input id="main_css_font_size" value="''' + set_data["main_css_font_size"] + '''">
                        <h2>''' + load_lang("editor") + '''</h2>
                        <h3>''' + load_lang("image_paste") + '''</h3>
                        <sup>''' + load_lang('only_korean') + '''</sup> <sup>''' + load_lang('unavailable_in_monaco') + '''</sup>
                        <hr class="main_hr">
                        <select name="main_css_image_paste">
                            ''' + set_data["main_css_image_paste"] + '''
                        </select>
                        <h3>''' + load_lang("monaco_editor") + '''</h3>
                        <select name="main_css_monaco">
                            ''' + set_data["main_css_monaco"] + '''
                        </select>
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('save') + '''</button>
                    </form>
                '''),
                menu = [['change', load_lang('user_setting')], ['change/skin_set', load_lang('skin_set')]]
            ))