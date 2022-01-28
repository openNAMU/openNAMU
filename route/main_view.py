from .tool.func import *
from .main_error_404 import main_error_404

def main_view(name = ''):
    file_name = re.search(r'([^/]+)$', name)
    if not file_name:
        return main_error_404(conn)
    else:
        file_name = file_name.group(1)
        dir_name = './views/' + re.sub(r'\.{2,}', '', re.sub(r'([^/]+)$', '', name))

        mime_type = file_name.split('.')
        if len(mime_type) < 2:
            mime_type = 'text/plain'
        else:
            mime_type = mime_type[len(mime_type) - 1].lower()
            image_type = ['jpeg', 'jpg', 'gif', 'png', 'webp']
            if mime_type in image_type:
                mime_type = 'image/' + mime_type
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
