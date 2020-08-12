from .tool.func import *
from . import main_error_404

def main_image_view_2(conn, name, app_var):
    curs = conn.cursor()

    if os.path.exists(os.path.join(app_var['path_data_image'], name)):
        return flask.send_from_directory(
            './' + app_var['path_data_image'], name, 
            mimetype = 'image/' + re.search(r'\.([^\.]+)$', name).group(1)
        )
    else:
        return main_error_404.main_error_404_2(conn)