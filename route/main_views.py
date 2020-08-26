from .tool.func import *
from . import main_error_404

def main_views_2(conn, name):
    curs = conn.cursor()

    file_name = re.search(r'([^/]+)$', name)
    if not file_name:
        return main_error_404.main_error_404_2(conn)
    else:
        file_name = file_name.group(1)
        dir_name = './views/' + re.sub(r'\.{2,}', '', re.sub(r'([^/]+)$', '', name))

        mime_type = re.search(r'([^.]+)$', file_name)
        image_type = [
            '.jpeg', 
            '.jpg', 
            '.gif', 
            '.png', 
            '.webp'
        ]
        if mime_type:
            mime_type = mime_type.group(1).lower()
            if mime_type in image_type:
                mime_type = 'image/' + mime_type
            else:
                mime_type = 'text/' + mime_type
        else:
            mime_type = 'text/plain'

        return flask.send_from_directory(
            dir_name, file_name, 
            mimetype = mime_type
        )