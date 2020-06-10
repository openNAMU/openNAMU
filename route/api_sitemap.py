from .tool.func import *

def api_sitemap_2(conn):
    curs = conn.cursor()

    if admin_check() == 1:
        data = '' + \
            '<?xml version="1.0" encoding="UTF-8"?>\n' + \
            '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n' + \
        ''

        curs.execute(db_change("select title from data"))
        for i in curs.fetchall():
            data += '<url><loc>' + flask.request.host_url + 'w/' + url_pas(i[0]) + '</loc></url>\n'

        data += '' + \
            '</urlset>' + \
        ''

        return flask.Response(data, mimetype = 'text/xml')
    else:
        return re_error('/ban')