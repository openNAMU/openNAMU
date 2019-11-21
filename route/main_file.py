from .tool.func import *

def main_file_2(conn, data):
    curs = conn.cursor()

    if data == 'robots.txt' and not os.path.exists('robots.txt'):
        return flask.Response('User-agent: *\nDisallow: /\nAllow: /$\nAllow: /w/', mimetype='text/plain')
    else:
        return flask.send_from_directory('./', data)