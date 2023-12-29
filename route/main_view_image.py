from .tool.func import *
from .main_func_error_404 import main_func_error_404

def main_view_image(name = ''):
    with get_db_connect() as conn:
        name = re.sub(r'\.cache_v(?:[0-9]+)$', '', name)
        mime_type = re.search(r'([^.]+)$', name)
        if mime_type:
            mime_type = mime_type.group(1).lower()
            if mime_type == 'svg':
                mime_type = 'svg+xml'

            return flask.send_from_directory(
                './' + load_image_url(), name, 
                mimetype = 'image/' + mime_type
            )
        else:
            return main_func_error_404()