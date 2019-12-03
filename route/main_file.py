from .tool.func import *
from . import main_error_404

def main_file_2(conn, data):
    curs = conn.cursor()

    print(data)
    if re.search('\.txt$', data) or data == 'sitemap.xml':
        if data == 'robots.txt' and not os.path.exists('robots.txt'):
            return flask.Response('User-agent: *\nDisallow: /\nAllow: /$\nAllow: /w/', mimetype='text/plain')
        elif os.path.exists(data):
            return flask.send_from_directory('./', data)
        else:
            return main_error_404.main_error_404_2(conn)
    else:
        return main_error_404.main_error_404_2(conn)