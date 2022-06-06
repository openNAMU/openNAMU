from .tool.func import *

def api_sitemap():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check(None, 'make sitemap') == 1:
            data = '' + \
                '<?xml version="1.0" encoding="UTF-8"?>\n' + \
                '<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n' + \
            ''
            domain = load_domain('full')

            curs.execute(db_change("select title from data"))
            all_data = curs.fetchall()

            len_all_data = len(all_data)
            count = int(len_all_data / 30000)
            other_count = len_all_data % 30000

            for i in range(count + 1):
                data += '<sitemap><loc>' + domain + 'sitemap_' + str(i) + '.xml</loc></sitemap>\n'

            data += '' + \
                '</sitemapindex>' + \
            ''

            f = open("sitemap.xml", 'w')
            f.write(data)
            f.close()

            e = 0
            for i in range(count + 1):
                data = '' + \
                    '<?xml version="1.0" encoding="UTF-8"?>\n' + \
                    '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n' + \
                ''

                if count == i:
                    for x in all_data[30000 * i:]:
                        data += '<url><loc>' + domain + 'w/' + url_pas(x[0]) + '</loc></url>\n'
                else:
                    for x in all_data[30000 * i:30000 * (i + 1)]:
                        data += '<url><loc>' + domain + 'w/' + url_pas(x[0]) + '</loc></url>\n'

                data += '' + \
                    '</urlset>' + \
                ''

                f = open("sitemap_" + str(i) + ".xml", 'w')
                f.write(data)
                f.close()

            return redirect('/sitemap.xml')
        else:
            return re_error('/ban')