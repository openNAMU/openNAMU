from .tool.func import *

async def main_setting_sitemap(do_type = 0):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if not do_type == 1:
            if await acl_check('', 'owner_auth', '', '') == 1:
                return await re_error(conn, 0)
        
        if do_type == 1 or flask.request.method == 'POST':
            if not do_type == 1:
                await acl_check(tool = 'owner_auth', memo = 'make sitemap')

            data = '' + \
                '<?xml version="1.0" encoding="UTF-8"?>\n' + \
                '<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n' + \
            ''

            curs.execute(db_change('select data from other where name = "sitemap_auto_exclude_domain"'))
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '':
                domain = ''
            else:
                domain = load_domain(conn, 'full')

            sql_add = ''

            curs.execute(db_change('select data from other where name = "sitemap_auto_exclude_user_page"'))
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '':
                sql_add += ' title not like "user:%"'

            curs.execute(db_change('select data from other where name = "sitemap_auto_exclude_file_page"'))
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '':
                if sql_add != '':
                    sql_add += ' and'

                sql_add += ' title not like "file:%"'

            curs.execute(db_change('select data from other where name = "sitemap_auto_exclude_category_page"'))
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '':
                if sql_add != '':
                    sql_add += ' and'

                sql_add += ' title not like "category:%"'

            if sql_add != '':
                sql_add = ' where' + sql_add

            curs.execute(db_change("select title from data" + sql_add))
            all_data = curs.fetchall()

            len_all_data = len(all_data)
            count = int(len_all_data / 30000)
            other_count = len_all_data % 30000

            for i in range(count + 1):
                data += '<sitemap><loc>' + domain + '/sitemap_' + str(i) + '.xml</loc></sitemap>\n'

            data += '' + \
                '</sitemapindex>' + \
            ''

            f = open("sitemap.xml", 'w')
            f.write(data)
            f.close()

            for i in range(count + 1):
                data = '' + \
                    '<?xml version="1.0" encoding="UTF-8"?>\n' + \
                    '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n' + \
                ''

                if count == i:
                    for x in all_data[30000 * i:]:
                        data += '<url><loc>' + domain + '/w/' + url_pas(x[0]) + '</loc></url>\n'
                else:
                    for x in all_data[30000 * i:30000 * (i + 1)]:
                        data += '<url><loc>' + domain + '/w/' + url_pas(x[0]) + '</loc></url>\n'

                data += '' + \
                    '</urlset>' + \
                ''

                f = open("sitemap_" + str(i) + ".xml", 'w')
                f.write(data)
                f.close()

            if not do_type == 1:
                return redirect(conn, '/setting/sitemap')
            else:
                return ''
        else:
            return await render_template(
                await get_lang('sitemap_manual_create'),
                '''
                    <form method="post">
                        <button class="__ON_BUTTON__" id="opennamu_save_button" type="submit">''' + await get_lang('create') + '''</button>
                    </form>
                ''',
                0,
                [['setting/sitemap_set', await get_lang('return')]]
            )
