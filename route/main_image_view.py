from .tool.func import *
from . import main_error_404

def main_image_view_2(conn, name, app_var):
    curs = conn.cursor()

    mime_type = re.search(r'([^.]+)$', name)
    if mime_type:
        return flask.send_from_directory(
            './' + app_var['path_data_image'], name, 
            mimetype = 'image/' + mime_type.group(1).lower()
        )
    else:
        return main_error_404.main_error_404_2(conn)