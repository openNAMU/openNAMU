from .tool.func import *
from .main_func_error_404 import main_func_error_404

def main_view(name = ''):
    with get_db_connect() as conn:
        file_name = re.search(r'([^/]+)$', name)
        if not file_name:
            return main_func_error_404()
        else:
            file_name = file_name.group(1)
            dir_name = './views/' + re.sub(r'\.{2,}', '', name[:-len(file_name)])

            file_name = re.sub(r'\.cache_v(?:[0-9]+)$', '', file_name)

            mime_type = file_name.split('.')
            if len(mime_type) < 2:
                mime_type = 'text/plain'
            else:
                mime_type = mime_type[len(mime_type) - 1].lower()
                image_type = ['jpeg', 'jpg', 'gif', 'png', 'webp', 'ico', 'svg']
                if mime_type in image_type:
                    if not mime_type == 'svg':
                        mime_type = 'image/' + mime_type
                    else:
                        mime_type = 'image/svg+xml'
                elif mime_type == 'js':
                    mime_type = 'text/javascript'
                elif mime_type == 'txt':
                    mime_type = 'text/plain'
                else:
                    mime_type = 'text/' + mime_type

            return flask.send_from_directory(
                dir_name, file_name, 
                mimetype = mime_type
            )