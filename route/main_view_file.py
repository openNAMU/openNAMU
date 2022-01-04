from .tool.func import *
from .main_error_404 import main_error_404

def main_view_file(data = ''):
    if data == 'robots.txt' and not os.path.exists('robots.txt'):
        return flask.Response('User-agent: *\nDisallow: /\nAllow: /$\nAllow: /w/', mimetype = 'text/plain')
    elif os.path.exists(data):
        if re.search(r'\.txt$', data, flags = re.I):
            return flask.send_from_directory('./', data, mimetype = 'text/plain')
        else:
            return flask.send_from_directory('./', data, mimetype = 'text/xml')

    return main_error_404()