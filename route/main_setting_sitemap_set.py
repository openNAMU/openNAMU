from .tool.func import *

async def main_setting_sitemap_set():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 0)
        
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

            await acl_check(tool = 'owner_auth', memo = 'edit_set (sitemap)')

            return redirect(conn, '/setting/sitemap_set')
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

            check_box_div = [0, 1, 2, 3, 4, '']
            for i in range(0, len(check_box_div)):
                acl_num = check_box_div[i]
                if acl_num != '' and d_list[acl_num]:
                    check_box_div[i] = 'checked="checked"'
                else:
                    check_box_div[i] = ''

            sitemap_list = ''
            if os.path.exists('sitemap.xml'):
                sitemap_list += '<a href="/sitemap.xml">(' + get_lang(conn, 'view') + ')</a>'

                for_a = 0
                while os.path.exists('sitemap_' + str(for_a) + '.xml'):
                    sitemap_list += ' <a href="/sitemap_' + str(for_a) + '.xml">(sitemap_' + str(for_a) + '.xml)</a>'

                    for_a += 1

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'sitemap_management'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    ''' + sitemap_list + '''
                    <hr class="main_hr">
                    <form method="post">
                        <a href="/setting/sitemap">(''' + get_lang(conn, 'sitemap_manual_create') + ''')</a>
                        <hr class="main_hr">

                        <label><input type="checkbox" ''' + check_box_div[4] + ''' name="sitemap_auto_make"> ''' + get_lang(conn, 'sitemap_auto_make') + '''</label>
                        <hr class="main_hr">

                        <label><input type="checkbox" ''' + check_box_div[0] + ''' name="sitemap_auto_exclude_domain"> ''' + get_lang(conn, 'stiemap_exclude_domain') + '''</label>
                        <hr class="main_hr">

                        <label><input type="checkbox" ''' + check_box_div[1] + ''' name="sitemap_auto_exclude_user_page"> ''' + get_lang(conn, 'stiemap_exclude_user_page') + '''</label>
                        <hr class="main_hr">

                        <label><input type="checkbox" ''' + check_box_div[2] + ''' name="sitemap_auto_exclude_file_page"> ''' + get_lang(conn, 'stiemap_exclude_file_page') + '''</label>
                        <hr class="main_hr">

                        <label><input type="checkbox" ''' + check_box_div[3] + ''' name="sitemap_auto_exclude_category_page"> ''' + get_lang(conn, 'stiemap_exclude_category_page') + '''</label>
                        <hr class="main_hr">

                        <button id="opennamu_save_button" type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                ''',
                menu = [['setting', get_lang(conn, 'return')]]
            ))