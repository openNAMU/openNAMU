from .tool.func import *

def api_func_language(db_set, data = 'Test'):
    other_set = {}
    if flask.request.method == 'POST':
        other_set["data"] = flask.request.form.get('data', '')
        other_set["data"] = other_set["data"].split(' ')
    else:
        other_set["data"] = [data]
    
    other_set = json.dumps(other_set)

    if platform.system() == 'Linux':
        if platform.machine() in ["AMD64", "x86_64"]:
            data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.amd64.bin"), sys._getframe().f_code.co_name, db_set, other_set], stdout = subprocess.PIPE).communicate()[0]
        else:
            data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.arm64.bin"), sys._getframe().f_code.co_name, db_set, other_set], stdout = subprocess.PIPE).communicate()[0]
    else:
        if platform.machine() in ["AMD64", "x86_64"]:
            data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.amd64.exe"), sys._getframe().f_code.co_name, db_set, other_set], stdout = subprocess.PIPE).communicate()[0]
        else:
            data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.arm64.exe"), sys._getframe().f_code.co_name, db_set, other_set], stdout = subprocess.PIPE).communicate()[0]

    data = data.decode('utf8')

    return flask.Response(response = data, status = 200, mimetype = 'application/json')