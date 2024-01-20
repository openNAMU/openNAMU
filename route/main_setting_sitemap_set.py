from .tool.func import *

def main_setting_sitemap_set():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/ban')
        
        setting_list = {
            0 : ['sitemap_auto_exclude_domain', ''],
            1 : ['sitemap_auto_exclude_user_page', ''],
            2 : ['sitemap_auto_exclude_file_page', ''],
            3 : ['sitemap_auto_exclude_category_page', ''],
            4 : ['sitemap_auto_make', '']
        }

        if flask.request.method == 'POST':
            for i in setting_list:
                curs.execute(db_change("update other set data = ? where name = ?"), [
                    flask.request.form.get(setting_list[i][0], setting_list[i][1]),
                    setting_list[i][0]
                ])

            conn.commit()

            admin_check(None, 'edit_set (sitemap)')

            return redirect('/setting/sitemap_set')
        else:
            d_list = {}
            for i in setting_list:
                curs.execute(db_change('select data from other where name = ?'), [setting_list[i][0]])
                db_data = curs.fetchall()
                if not db_data:
                    curs.execute(db_change('insert into other (name, data, coverage) values (?, ?, "")'), [
                        setting_list[i][0],
                        setting_list[i][1]
                    ])

                d_list[i] = db_data[0][0] if db_data else setting_list[i][1]
            else:
                conn.commit()

            check_box_div = [0, 1, 2, 3, 4, '']
            for i in range(0, len(check_box_div)):
                acl_num = check_box_div[i]
                if acl_num != '' and d_list[acl_num]:
                    check_box_div[i] = 'checked="checked"'
                else:
                    check_box_div[i] = ''

            sitemap_list = ''
            if os.path.exists('sitemap.xml'):
                sitemap_list += '<a href="/sitemap.xml">(' + load_lang('view') + ')</a>'

                for_a = 0
                while os.path.exists('sitemap_' + str(for_a) + '.xml'):
                    sitemap_list += ' <a href="/sitemap_' + str(for_a) + '.xml">(sitemap_' + str(for_a) + '.xml)</a>'

                    for_a += 1

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('sitemap_management'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    ''' + sitemap_list + '''
                    <hr class="main_hr">
                    <form method="post">
                        <a href="/setting/sitemap">(''' + load_lang('sitemap_manual_create') + ''')</a>
                        <hr class="main_hr">

                        <input type="checkbox" ''' + check_box_div[4] + ''' name="sitemap_auto_make"> ''' + load_lang('sitemap_auto_make') + '''
                        <hr class="main_hr">

                        <input type="checkbox" ''' + check_box_div[0] + ''' name="sitemap_auto_exclude_domain"> ''' + load_lang('stiemap_exclude_domain') + '''
                        <hr class="main_hr">

                        <input type="checkbox" ''' + check_box_div[1] + ''' name="sitemap_auto_exclude_user_page"> ''' + load_lang('stiemap_exclude_user_page') + '''
                        <hr class="main_hr">

                        <input type="checkbox" ''' + check_box_div[2] + ''' name="sitemap_auto_exclude_file_page"> ''' + load_lang('stiemap_exclude_file_page') + '''
                        <hr class="main_hr">

                        <input type="checkbox" ''' + check_box_div[3] + ''' name="sitemap_auto_exclude_category_page"> ''' + load_lang('stiemap_exclude_category_page') + '''
                        <hr class="main_hr">

                        <button id="opennamu_save_button" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['setting', load_lang('return')]]
            ))