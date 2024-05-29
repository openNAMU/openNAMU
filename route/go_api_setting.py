from .tool.func import *

def api_setting(name = 'Test'):
    other_set = {}
    other_set["set_name"] = name
    other_set["ip"] = ip_check()

    func_name = sys._getframe().f_code.co_name
    if flask.request.method == 'PUT':
        func_name += '_edit'
        other_set['data'] = flask.request.form.get('data', 'Test')
    
    other_set = json.dumps(other_set)

    if platform.system() == 'Linux':
        if platform.machine() in ["AMD64", "x86_64"]:
            data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.amd64.bin"), func_name, other_set], stdout = subprocess.PIPE).communicate()[0]
        else:
            data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.arm64.bin"), func_name, other_set], stdout = subprocess.PIPE).communicate()[0]
    else:
        if platform.machine() in ["AMD64", "x86_64"]:
            data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.amd64.exe"), func_name, other_set], stdout = subprocess.PIPE).communicate()[0]
        else:
            data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.arm64.exe"), func_name, other_set], stdout = subprocess.PIPE).communicate()[0]

    data = data.decode('utf8')

    return flask.Response(response = data, status = 200, mimetype = 'application/json')