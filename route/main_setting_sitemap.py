from .tool.func import *

def main_setting_sitemap():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/ban')
        
        if flask.request.method == 'POST':
            admin_check(None, 'make sitemap')

            data = '' + \
                '<?xml version="1.0" encoding="UTF-8"?>\n' + \
                '<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n' + \
            ''

            if flask.request.form.get('exclude_domain', None):
                domain = ''
            else:
                domain = load_domain('full')

            sql_add = ''
            if flask.request.form.get('exclude_user_page', None):
                sql_add += ' title not like "user:%"'

            if flask.request.form.get('exclude_file_page', None):
                if sql_add != '':
                    sql_add += ' and'

                sql_add += ' title not like "file:%"'

            if flask.request.form.get('exclude_category_page', None):
                if sql_add != '':
                    sql_add += ' and'

                sql_add += ' title not like "category:%"'

            if sql_add != '':
                sql_add = ' where' + sql_add

            print(sql_add)
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

            return redirect('/setting/sitemap')
        else:
            sitemap_list = ''
            if os.path.exists('sitemap.xml'):
                sitemap_list += '<a href="/sitemap.xml">(' + load_lang('view') + ')</a>'

                for_a = 0
                while os.path.exists('sitemap_' + str(for_a) + '.xml'):
                    sitemap_list += ' <a href="/sitemap_' + str(for_a) + '.xml">(sitemap_' + str(for_a) + '.xml)</a>'

                    for_a += 1

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('sitemap_management'), wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('beta') + ')', 0])],
                data = '''
                    ''' + sitemap_list + '''
                    <hr class="main_hr">
                    <form method="post">
                        <input type="checkbox" name="exclude_domain"> ''' + load_lang('stiemap_exclude_domain') + '''
                        <hr class="main_hr">

                        <input type="checkbox" name="exclude_user_page"> ''' + load_lang('stiemap_exclude_user_page') + '''
                        <hr class="main_hr">

                        <input type="checkbox" name="exclude_file_page"> ''' + load_lang('stiemap_exclude_file_page') + '''
                        <hr class="main_hr">

                        <input type="checkbox" name="exclude_category_page"> ''' + load_lang('stiemap_exclude_category_page') + '''
                        <hr class="main_hr">

                        <button id="opennamu_save_button" type="submit">''' + load_lang('create') + '''</button>
                    </form>
                ''',
                menu = [['setting', load_lang('return')]]
            ))