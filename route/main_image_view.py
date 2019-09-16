from .tool.func import *

def main_image_view_2(name, app_var):
    
    
    if os.path.exists(os.path.join(app_var['path_data_image'], name)):
        return flask.send_from_directory('./' + app_var['path_data_image'], name)
    else:
        return redirect()