from .tool.func import *
from .main_error_404 import main_error_404

def main_view_image(name = ''):
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
        return main_error_404()